package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/google/uuid"

	"poc-vision-api/internal/ai"
	// Importar proveedores para que se registren automáticamente
	_ "poc-vision-api/internal/ai/anthropic"
	_ "poc-vision-api/internal/ai/bedrock"
	_ "poc-vision-api/internal/ai/google"
	_ "poc-vision-api/internal/ai/openai"
)

// Handler maneja las solicitudes de la Lambda
type Handler struct {
	aiProvider      ai.Provider
	dynamoClient    *dynamodb.Client
	lambdaClient    *awslambda.Client
	tableName       string
	mcpServerLambda string
}

// AnalyzeRequest representa la solicitud de análisis
type AnalyzeRequest struct {
	ImageBase64 string `json:"image_base64"`
	MediaType   string `json:"media_type"`
	RequestID   string `json:"request_id,omitempty"`
}

// AnalysisResult representa el resultado del análisis
type AnalysisResult struct {
	RequestID       string            `json:"request_id"`
	HasMatch        bool              `json:"has_match"`
	MatchedProvider *ai.ProviderMatch `json:"matched_provider,omitempty"`
	Justification   string            `json:"justification"`
	ConfidenceLevel string            `json:"confidence_level"`
	ExtractedInfo   *ai.ExtractedInfo `json:"extracted_info,omitempty"`
	ProcessingTime  string            `json:"processing_time"`
	ModelUsed       string            `json:"model_used"`
	ProviderName    string            `json:"provider_name"`
	ToolCalls       int               `json:"tool_calls"`
	Latency         map[string]int64  `json:"latency,omitempty"`
}

// APIResponse representa la respuesta de la API
type APIResponse struct {
	Success bool            `json:"success"`
	Data    *AnalysisResult `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// Cache para el catálogo de proveedores (persiste entre invocaciones en warm starts)
var (
	cachedCatalog string
	cacheExpiry   time.Time
	cacheTTL      = 5 * time.Minute
)

// ToolExecutorImpl implementa ai.ToolExecutor
type ToolExecutorImpl struct {
	lambdaClient    *awslambda.Client
	mcpServerLambda string
}

// ExecuteTool ejecuta una herramienta por nombre
func (t *ToolExecutorImpl) ExecuteTool(ctx context.Context, toolName string, input map[string]interface{}) (string, error) {
	switch toolName {
	case "get_provider_catalog":
		return t.getProviderCatalog(ctx)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// GetAvailableTools retorna las herramientas disponibles
func (t *ToolExecutorImpl) GetAvailableTools() []ai.ToolDefinition {
	return []ai.ToolDefinition{ai.GetDefaultToolDefinition()}
}

func (t *ToolExecutorImpl) getProviderCatalog(ctx context.Context) (string, error) {
	// Usar cache si está válido
	if cachedCatalog != "" && time.Now().Before(cacheExpiry) {
		log.Printf("[Cache] Using cached catalog (expires in %v)", time.Until(cacheExpiry).Round(time.Second))
		return cachedCatalog, nil
	}

	log.Printf("[Cache] Cache miss or expired, fetching from MCP server")

	payload := map[string]interface{}{
		"action": "get_all_providers",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	mcpCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	output, err := t.lambdaClient.Invoke(mcpCtx, &awslambda.InvokeInput{
		FunctionName:   aws.String(t.mcpServerLambda),
		Payload:        payloadBytes,
		InvocationType: lambdatypes.InvocationTypeRequestResponse,
	})
	if err != nil {
		return "", fmt.Errorf("mcp invoke failed: %w", err)
	}

	result := string(output.Payload)

	// Guardar en cache
	cachedCatalog = result
	cacheExpiry = time.Now().Add(cacheTTL)
	log.Printf("[Cache] Catalog cached for %v", cacheTTL)

	return result, nil
}

// processImage procesa una imagen y retorna el resultado
func (h *Handler) processImage(ctx context.Context, req AnalyzeRequest) (*AnalysisResult, error) {
	startTime := time.Now()
	latency := make(map[string]int64)

	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	// Validar imagen
	if _, err := base64.StdEncoding.DecodeString(req.ImageBase64); err != nil {
		return nil, fmt.Errorf("invalid base64 image data")
	}

	validMediaTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !validMediaTypes[req.MediaType] {
		return nil, fmt.Errorf("unsupported media type: %s", req.MediaType)
	}

	// Crear ejecutor de herramientas
	toolExecutor := &ToolExecutorImpl{
		lambdaClient:    h.lambdaClient,
		mcpServerLambda: h.mcpServerLambda,
	}

	// Preparar solicitud de análisis
	analysisRequest := ai.AnalysisRequest{
		ImageBase64:  req.ImageBase64,
		MediaType:    req.MediaType,
		Tools:        toolExecutor.GetAvailableTools(),
		ToolExecutor: toolExecutor,
		MaxTokens:    4096,
		Temperature:  0.0,
	}

	// Ejecutar análisis con el proveedor de IA
	analysisStart := time.Now()
	response, err := h.aiProvider.AnalyzeImageWithTools(ctx, analysisRequest)
	latency["analysis_with_tools"] = time.Since(analysisStart).Milliseconds()

	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Construir resultado
	result := &AnalysisResult{
		RequestID:       req.RequestID,
		HasMatch:        response.HasMatch,
		MatchedProvider: response.MatchedProvider,
		Justification:   response.Justification,
		ConfidenceLevel: response.ConfidenceLevel,
		ExtractedInfo:   response.ExtractedInfo,
		ProcessingTime:  time.Since(startTime).String(),
		ModelUsed:       response.ModelUsed,
		ProviderName:    response.ProviderName,
		ToolCalls:       response.ToolCallsCount,
		Latency:         latency,
	}

	// Guardar en DynamoDB (asíncrono)
	if h.tableName != "" {
		go func() {
			saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := h.saveResults(saveCtx, result); err != nil {
				log.Printf("Failed to save results: %v", err)
			}
		}()
	}

	return result, nil
}

func (h *Handler) saveResults(ctx context.Context, result *AnalysisResult) error {
	resultJSON, _ := json.Marshal(result)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	item := map[string]ddbtypes.AttributeValue{
		"RequestID":        &ddbtypes.AttributeValueMemberS{Value: result.RequestID},
		"AnalysisResult":   &ddbtypes.AttributeValueMemberS{Value: string(resultJSON)},
		"ProcessedAt":      &ddbtypes.AttributeValueMemberS{Value: timestamp},
		"HasProviderMatch": &ddbtypes.AttributeValueMemberBOOL{Value: result.HasMatch},
		"ConfidenceLevel":  &ddbtypes.AttributeValueMemberS{Value: result.ConfidenceLevel},
		"ModelUsed":        &ddbtypes.AttributeValueMemberS{Value: result.ModelUsed},
		"ProviderName":     &ddbtypes.AttributeValueMemberS{Value: result.ProviderName},
		"ToolCalls":        &ddbtypes.AttributeValueMemberN{Value: fmt.Sprintf("%d", result.ToolCalls)},
	}

	if result.HasMatch && result.MatchedProvider != nil {
		providerJSON, _ := json.Marshal(result.MatchedProvider)
		item["MatchedProviderJSON"] = &ddbtypes.AttributeValueMemberS{Value: string(providerJSON)}
	}

	_, err := h.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(h.tableName),
		Item:      item,
	})

	return err
}

func (h *Handler) parseRequest(request events.APIGatewayProxyRequest) (AnalyzeRequest, error) {
	contentType := request.Headers["Content-Type"]
	if contentType == "" {
		contentType = request.Headers["content-type"]
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType = "application/json"
	}

	if strings.HasPrefix(mediaType, "multipart/form-data") {
		return h.parseMultipartRequest(request, params["boundary"])
	}

	return h.parseJSONRequest(request)
}

func (h *Handler) parseMultipartRequest(request events.APIGatewayProxyRequest, boundary string) (AnalyzeRequest, error) {
	var req AnalyzeRequest

	if boundary == "" {
		return req, fmt.Errorf("multipart boundary not found")
	}

	var body []byte
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return req, fmt.Errorf("failed to decode base64 body: %w", err)
		}
		body = decoded
	} else {
		body = []byte(request.Body)
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return req, fmt.Errorf("failed to read multipart: %w", err)
		}

		fieldName := part.FormName()
		switch fieldName {
		case "image", "file":
			data, err := io.ReadAll(part)
			if err != nil {
				return req, fmt.Errorf("failed to read image data: %w", err)
			}

			req.ImageBase64 = base64.StdEncoding.EncodeToString(data)

			partContentType := part.Header.Get("Content-Type")
			if partContentType != "" {
				req.MediaType = partContentType
			} else {
				req.MediaType = h.detectMediaType(part.FileName(), data)
			}

		case "request_id":
			data, _ := io.ReadAll(part)
			req.RequestID = string(data)

		case "media_type":
			data, _ := io.ReadAll(part)
			req.MediaType = string(data)
		}
	}

	return req, nil
}

func (h *Handler) parseJSONRequest(request events.APIGatewayProxyRequest) (AnalyzeRequest, error) {
	var req AnalyzeRequest

	body := request.Body
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return req, fmt.Errorf("failed to decode base64 body: %w", err)
		}
		body = string(decoded)
	}

	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return req, fmt.Errorf("invalid JSON request body: %w", err)
	}

	return req, nil
}

func (h *Handler) detectMediaType(filename string, data []byte) string {
	if len(data) >= 8 {
		switch {
		case data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
			return "image/jpeg"
		case data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47:
			return "image/png"
		case data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46:
			return "image/gif"
		case data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46:
			return "image/webp"
		}
	}

	detected := http.DetectContentType(data)
	if strings.HasPrefix(detected, "image/") {
		return detected
	}

	filename = strings.ToLower(filename)
	switch {
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filename, ".webp"):
		return "image/webp"
	}

	return "image/jpeg"
}

// HandleRequest maneja la solicitud HTTP
func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, x-api-key",
			},
		}, nil
	}

	if request.HTTPMethod != "POST" {
		response := APIResponse{Success: false, Error: "Method not allowed"}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	req, err := h.parseRequest(request)
	if err != nil {
		response := APIResponse{Success: false, Error: err.Error()}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	if req.ImageBase64 == "" {
		response := APIResponse{Success: false, Error: "image file is required"}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	if req.MediaType == "" {
		req.MediaType = "image/jpeg"
	}

	result, err := h.processImage(ctx, req)
	if err != nil {
		log.Printf("Processing error: %v", err)
		response := APIResponse{Success: false, Error: err.Error()}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	response := APIResponse{Success: true, Data: result}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRetryMode(aws.RetryModeStandard),
		config.WithRetryMaxAttempts(1),
	)
	if err != nil {
		log.Fatalf("Config failed: %v", err)
	}

	tableName := os.Getenv("DYNAMODB_TABLE")
	mcpServerLambda := os.Getenv("MCP_SERVER_FUNCTION_NAME")

	if mcpServerLambda == "" {
		log.Fatal("MCP_SERVER_FUNCTION_NAME is required")
	}

	// Crear proveedor de IA usando la fábrica
	aiProvider, err := ai.CreateProviderFromEnv(ctx)
	if err != nil {
		log.Fatalf("Failed to create AI provider: %v", err)
	}

	log.Printf("Using AI provider: %s (Vision: %v, ToolUse: %v)",
		aiProvider.Name(),
		aiProvider.SupportsVision(),
		aiProvider.SupportsToolUse())

	handler := &Handler{
		aiProvider:      aiProvider,
		dynamoClient:    dynamodb.NewFromConfig(cfg),
		lambdaClient:    awslambda.NewFromConfig(cfg),
		tableName:       tableName,
		mcpServerLambda: mcpServerLambda,
	}

	lambda.Start(handler.HandleRequest)
}
