package google

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

// Provider implementa ai.Provider para Google Gemini
type Provider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// New crea una nueva instancia del proveedor Google
func New(ctx context.Context, cfg ai.ProviderConfig) (*Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Google API key is required")
	}

	model := cfg.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 120
	}

	return &Provider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

func (p *Provider) Name() string {
	return "Google Gemini"
}

func (p *Provider) SupportsVision() bool {
	return true
}

func (p *Provider) SupportsToolUse() bool {
	return true
}

// Estructuras de request/response Gemini
type generateRequest struct {
	Contents          []content          `json:"contents"`
	SystemInstruction *systemInstruction `json:"systemInstruction,omitempty"`
	Tools             []toolDeclaration  `json:"tools,omitempty"`
	GenerationConfig  *generationConfig  `json:"generationConfig,omitempty"`
}

type systemInstruction struct {
	Parts []part `json:"parts"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type part struct {
	Text             string            `json:"text,omitempty"`
	InlineData       *inlineData       `json:"inlineData,omitempty"`
	FunctionCall     *functionCall     `json:"functionCall,omitempty"`
	FunctionResponse *functionResponse `json:"functionResponse,omitempty"`
}

type inlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type toolDeclaration struct {
	FunctionDeclarations []functionDeclaration `json:"functionDeclarations"`
}

type functionDeclaration struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type functionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type functionResponse struct {
	Name     string          `json:"name"`
	Response json.RawMessage `json:"response"`
}

type generationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type generateResponse struct {
	Candidates    []candidate    `json:"candidates"`
	UsageMetadata *usageMetadata `json:"usageMetadata,omitempty"`
}

type candidate struct {
	Content       content `json:"content"`
	FinishReason  string  `json:"finishReason"`
	SafetyRatings []any   `json:"safetyRatings,omitempty"`
}

type usageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
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

	contents := []content{
		{
			Role: "user",
			Parts: []part{
				{
					InlineData: &inlineData{
						MimeType: request.MediaType,
						Data:     request.ImageBase64,
					},
				},
				{
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
		reqBody := generateRequest{
			Contents: contents,
			SystemInstruction: &systemInstruction{
				Parts: []part{{Text: systemPrompt}},
			},
			Tools: tools,
			GenerationConfig: &generationConfig{
				Temperature:     request.Temperature,
				MaxOutputTokens: maxTokens,
			},
		}

		endpoint := fmt.Sprintf("/models/%s:generateContent?key=%s", p.model, p.apiKey)
		respBody, err := p.doRequest(ctx, endpoint, reqBody)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		var resp generateResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if len(resp.Candidates) == 0 {
			return nil, fmt.Errorf("no candidates in response")
		}

		candidate := resp.Candidates[0]

		// Buscar function calls
		var pendingCalls []functionCall
		for _, part := range candidate.Content.Parts {
			if part.FunctionCall != nil {
				pendingCalls = append(pendingCalls, *part.FunctionCall)
			}
		}

		// Si no hay function calls, es respuesta final
		if len(pendingCalls) == 0 || candidate.FinishReason == "STOP" {
			responseText := p.extractText(candidate.Content.Parts)
			result := ai.ParseAnalysisResponse(responseText, p.model, p.Name())
			result.ToolCallsCount = toolCalls
			return result, nil
		}

		// Agregar respuesta del modelo
		contents = append(contents, candidate.Content)

		// Ejecutar cada function call
		var responseParts []part
		for _, fc := range pendingCalls {
			toolCalls++
			log.Printf("[Google] Executing tool: %s", fc.Name)

			toolResult, err := request.ToolExecutor.ExecuteTool(ctx, fc.Name, fc.Args)
			var resultContent map[string]interface{}
			if err != nil {
				resultContent = map[string]interface{}{"error": err.Error()}
			} else {
				json.Unmarshal([]byte(toolResult), &resultContent)
				if resultContent == nil {
					resultContent = map[string]interface{}{"result": toolResult}
				}
			}

			resultJSON, _ := json.Marshal(resultContent)
			responseParts = append(responseParts, part{
				FunctionResponse: &functionResponse{
					Name:     fc.Name,
					Response: resultJSON,
				},
			})
		}

		contents = append(contents, content{
			Role:  "user",
			Parts: responseParts,
		})
	}

	return nil, fmt.Errorf("max iterations reached")
}

func (p *Provider) buildTools(tools []ai.ToolDefinition) []toolDeclaration {
	if len(tools) == 0 {
		tools = []ai.ToolDefinition{ai.GetDefaultToolDefinition()}
	}

	var declarations []functionDeclaration
	for _, t := range tools {
		declarations = append(declarations, functionDeclaration{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.InputSchema,
		})
	}

	return []toolDeclaration{{FunctionDeclarations: declarations}}
}

func (p *Provider) extractText(parts []part) string {
	for _, part := range parts {
		if part.Text != "" {
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
	ai.RegisterProvider(ai.ProviderGoogle, func(ctx context.Context, config ai.ProviderConfig) (ai.Provider, error) {
		return New(ctx, config)
	})
}
