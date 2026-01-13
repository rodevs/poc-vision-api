package ai

import "context"

// Provider define la interfaz para cualquier proveedor de IA generativa
type Provider interface {
	// AnalyzeImageWithTools analiza una imagen usando herramientas definidas
	AnalyzeImageWithTools(ctx context.Context, request AnalysisRequest) (*AnalysisResponse, error)

	// Name retorna el nombre del proveedor
	Name() string

	// SupportsVision indica si el proveedor soporta análisis de imágenes
	SupportsVision() bool

	// SupportsToolUse indica si el proveedor soporta Tool Use/Function Calling
	SupportsToolUse() bool
}

// ToolExecutor define la interfaz para ejecutar herramientas
type ToolExecutor interface {
	// ExecuteTool ejecuta una herramienta por nombre y retorna el resultado
	ExecuteTool(ctx context.Context, toolName string, input map[string]interface{}) (string, error)

	// GetAvailableTools retorna las definiciones de herramientas disponibles
	GetAvailableTools() []ToolDefinition
}

// AnalysisRequest representa una solicitud de análisis
type AnalysisRequest struct {
	ImageBase64  string
	MediaType    string
	SystemPrompt string
	UserPrompt   string
	Tools        []ToolDefinition
	ToolExecutor ToolExecutor
	MaxTokens    int
	Temperature  float32
}

// AnalysisResponse representa la respuesta del análisis
type AnalysisResponse struct {
	HasMatch        bool             `json:"has_match"`
	MatchedProvider *ProviderMatch   `json:"matched_provider,omitempty"`
	Justification   string           `json:"justification"`
	ConfidenceLevel string           `json:"confidence_level"`
	ExtractedInfo   *ExtractedInfo   `json:"extracted_info,omitempty"`
	RawResponse     string           `json:"raw_response,omitempty"`
	ToolCallsCount  int              `json:"tool_calls_count"`
	ModelUsed       string           `json:"model_used"`
	ProviderName    string           `json:"provider_name"`
	Latency         map[string]int64 `json:"latency,omitempty"`
}

// ProviderMatch representa un proveedor encontrado
type ProviderMatch struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	OfficialName           string   `json:"official_name"`
	CommonNames            []string `json:"common_names"`
	Aliases                []string `json:"aliases"`
	ServiceType            string   `json:"service_type"`
	Region                 string   `json:"region"`
	Status                 string   `json:"status"`
	OnlineServiceAvailable bool     `json:"online_service_available"`
	AcceptedPaymentField   string   `json:"accepted_payment_field"`
	AdditionalInformation  string   `json:"additional_information,omitempty"`
}

// ExtractedInfo representa información extraída de la imagen
type ExtractedInfo struct {
	ProviderName      string `json:"provider_name,omitempty"`
	ServiceType       string `json:"service_type,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	PaymentFieldName  string `json:"payment_field_name,omitempty"`
	PaymentFieldValue string `json:"payment_field_value,omitempty"`
	TotalAmount       string `json:"total_amount,omitempty"`
	Currency          string `json:"currency,omitempty"`
	ReceiptDate       string `json:"receipt_date,omitempty"`
	DueDate           string `json:"due_date,omitempty"`
	CustomerName      string `json:"customer_name,omitempty"`
}

// ToolDefinition define una herramienta disponible
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// ToolCall representa una llamada a herramienta solicitada por el modelo
type ToolCall struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// ToolResult representa el resultado de ejecutar una herramienta
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error"`
}

// ImageFormat representa el formato de imagen soportado
type ImageFormat string

const (
	ImageFormatJPEG ImageFormat = "image/jpeg"
	ImageFormatPNG  ImageFormat = "image/png"
	ImageFormatGIF  ImageFormat = "image/gif"
	ImageFormatWEBP ImageFormat = "image/webp"
	DocumentFormatPDF ImageFormat = "application/pdf"
)

// ProviderConfig contiene la configuración común para proveedores
type ProviderConfig struct {
	APIKey      string
	Model       string
	Region      string
	MaxRetries  int
	Timeout     int // segundos
	ExtraConfig map[string]interface{}
}
