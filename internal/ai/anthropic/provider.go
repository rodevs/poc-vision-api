package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"poc-vision-api/internal/ai"
)

// Provider implementa ai.Provider para Anthropic API directa
type Provider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// New crea una nueva instancia del proveedor Anthropic
func New(ctx context.Context, cfg ai.ProviderConfig) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	model := cfg.Model
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 120
	}

	return &Provider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: "https://api.anthropic.com/v1",
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

func (p *Provider) Name() string {
	return "Anthropic"
}

func (p *Provider) SupportsVision() bool {
	return true
}

func (p *Provider) SupportsToolUse() bool {
	return true
}

// Estructuras de request/response Anthropic
type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []message `json:"messages"`
	Tools     []tool    `json:"tools,omitempty"`
}

type message struct {
	Role    string        `json:"role"`
	Content []contentPart `json:"content"`
}

type contentPart struct {
	Type      string       `json:"type"`
	Text      string       `json:"text,omitempty"`
	Source    *imageSource `json:"source,omitempty"`
	ID        string       `json:"id,omitempty"`
	Name      string       `json:"name,omitempty"`
	Input     interface{}  `json:"input,omitempty"`
	ToolUseID string       `json:"tool_use_id,omitempty"`
	Content   string       `json:"content,omitempty"`
	IsError   bool         `json:"is_error,omitempty"`
}

type imageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

type messagesResponse struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Role         string        `json:"role"`
	Content      []contentPart `json:"content"`
	StopReason   string        `json:"stop_reason"`
	StopSequence *string       `json:"stop_sequence"`
	Usage        usage         `json:"usage"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// AnalyzeImageWithTools implementa el análisis de imagen con Tool Use
func (p *Provider) AnalyzeImageWithTools(ctx context.Context, request ai.AnalysisRequest) (*ai.AnalysisResponse, error) {
	systemPrompt := request.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = ai.GetDefaultSystemPrompt()
	}

	userPrompt := request.UserPrompt
	if userPrompt == "" {
		userPrompt = ai.GetDefaultUserPrompt()
	}

	messages := []message{
		{
			Role: "user",
			Content: []contentPart{
				{
					Type: "image",
					Source: &imageSource{
						Type:      "base64",
						MediaType: request.MediaType,
						Data:      request.ImageBase64,
					},
				},
				{
					Type: "text",
					Text: userPrompt,
				},
			},
		},
	}

	tools := p.buildTools(request.Tools)
	toolCalls := 0
	maxIterations := 5

	maxTokens := 4096
	if request.MaxTokens > 0 {
		maxTokens = request.MaxTokens
	}

	for i := 0; i < maxIterations; i++ {
		reqBody := messagesRequest{
			Model:     p.model,
			MaxTokens: maxTokens,
			System:    systemPrompt,
			Messages:  messages,
			Tools:     tools,
		}

		respBody, err := p.doRequest(ctx, "/messages", reqBody)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		var resp messagesResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Respuesta final
		if resp.StopReason == "end_turn" {
			responseText := p.extractText(resp.Content)
			result := ai.ParseAnalysisResponse(responseText, p.model, p.Name())
			result.ToolCallsCount = toolCalls
			return result, nil
		}

		// Tool Use
		if resp.StopReason == "tool_use" {
			messages = append(messages, message{
				Role:    "assistant",
				Content: resp.Content,
			})

			var toolResults []contentPart
			for _, part := range resp.Content {
				if part.Type == "tool_use" {
					toolCalls++
					log.Printf("[Anthropic] Executing tool: %s", part.Name)

					var toolInput map[string]interface{}
					if part.Input != nil {
						inputBytes, _ := json.Marshal(part.Input)
						json.Unmarshal(inputBytes, &toolInput)
					}

					toolResult, err := request.ToolExecutor.ExecuteTool(ctx, part.Name, toolInput)
					if err != nil {
						toolResults = append(toolResults, contentPart{
							Type:      "tool_result",
							ToolUseID: part.ID,
							Content:   err.Error(),
							IsError:   true,
						})
					} else {
						toolResults = append(toolResults, contentPart{
							Type:      "tool_result",
							ToolUseID: part.ID,
							Content:   toolResult,
						})
					}
				}
			}

			messages = append(messages, message{
				Role:    "user",
				Content: toolResults,
			})
			continue
		}

		// Otra razón
		responseText := p.extractText(resp.Content)
		result := ai.ParseAnalysisResponse(responseText, p.model, p.Name())
		result.ToolCallsCount = toolCalls
		return result, nil
	}

	return nil, fmt.Errorf("max iterations reached")
}

func (p *Provider) buildTools(tools []ai.ToolDefinition) []tool {
	if len(tools) == 0 {
		tools = []ai.ToolDefinition{ai.GetDefaultToolDefinition()}
	}

	var anthropicTools []tool
	for _, t := range tools {
		anthropicTools = append(anthropicTools, tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	return anthropicTools
}

func (p *Provider) extractText(content []contentPart) string {
	for _, part := range content {
		if part.Type == "text" {
			return part.Text
		}
	}
	return ""
}

func (p *Provider) doRequest(ctx context.Context, endpoint string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// init registra el proveedor en la fábrica global
func init() {
	ai.RegisterProvider(ai.ProviderAnthropic, func(ctx context.Context, config ai.ProviderConfig) (ai.Provider, error) {
		return New(ctx, config)
	})
}
