package openai

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

// Provider implementa ai.Provider para OpenAI
type Provider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// New crea una nueva instancia del proveedor OpenAI
func New(ctx context.Context, cfg ai.ProviderConfig) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	model := cfg.Model
	if model == "" {
		model = "gpt-4o"
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 120
	}

	return &Provider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

func (p *Provider) Name() string {
	return "OpenAI"
}

func (p *Provider) SupportsVision() bool {
	return true
}

func (p *Provider) SupportsToolUse() bool {
	return true
}

// Estructuras de request/response OpenAI
type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Tools       []tool    `json:"tools,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

type message struct {
	Role       string      `json:"role"`
	Content    interface{} `json:"content"` // string o []contentPart
	ToolCalls  []toolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type tool struct {
	Type     string   `json:"type"`
	Function function `json:"function"`
}

type function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type toolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function functionCall `json:"function"`
}

type functionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type chatResponse struct {
	ID      string   `json:"id"`
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
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

	// Construir mensaje con imagen
	imageDataURL := fmt.Sprintf("data:%s;base64,%s", request.MediaType, request.ImageBase64)

	messages := []message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role: "user",
			Content: []contentPart{
				{
					Type: "image_url",
					ImageURL: &imageURL{
						URL:    imageDataURL,
						Detail: "high",
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
		reqBody := chatRequest{
			Model:       p.model,
			Messages:    messages,
			Tools:       tools,
			MaxTokens:   maxTokens,
			Temperature: request.Temperature,
		}

		respBody, err := p.doRequest(ctx, "/chat/completions", reqBody)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		var resp chatResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no choices in response")
		}

		choice := resp.Choices[0]

		// Si no hay tool calls, es respuesta final
		if len(choice.Message.ToolCalls) == 0 || choice.FinishReason == "stop" {
			responseText := ""
			if content, ok := choice.Message.Content.(string); ok {
				responseText = content
			}
			result := ai.ParseAnalysisResponse(responseText, p.model, p.Name())
			result.ToolCallsCount = toolCalls
			return result, nil
		}

		// Hay tool calls
		messages = append(messages, choice.Message)

		for _, tc := range choice.Message.ToolCalls {
			toolCalls++
			log.Printf("[OpenAI] Executing tool: %s", tc.Function.Name)

			var toolInput map[string]interface{}
			json.Unmarshal([]byte(tc.Function.Arguments), &toolInput)

			toolResult, err := request.ToolExecutor.ExecuteTool(ctx, tc.Function.Name, toolInput)
			resultContent := toolResult
			if err != nil {
				resultContent = fmt.Sprintf("Error: %s", err.Error())
			}

			messages = append(messages, message{
				Role:       "tool",
				Content:    resultContent,
				ToolCallID: tc.ID,
			})
		}
	}

	return nil, fmt.Errorf("max iterations reached")
}

func (p *Provider) buildTools(tools []ai.ToolDefinition) []tool {
	if len(tools) == 0 {
		tools = []ai.ToolDefinition{ai.GetDefaultToolDefinition()}
	}

	var openaiTools []tool
	for _, t := range tools {
		openaiTools = append(openaiTools, tool{
			Type: "function",
			Function: function{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		})
	}

	return openaiTools
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

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
	ai.RegisterProvider(ai.ProviderOpenAI, func(ctx context.Context, config ai.ProviderConfig) (ai.Provider, error) {
		return New(ctx, config)
	})
}
