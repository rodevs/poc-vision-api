package models

import "time"

type Provider struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	OfficialName           string   `json:"official_name"`
	CommonNames            []string `json:"common_names"`
	Aliases                []string `json:"aliases"`
	ServiceType            string   `json:"service_type"`
	Region                 string   `json:"region"`
	Status                 string   `json:"status"`
	OnlineServiceAvailable bool     `json:"online_service_available"`
}

type ProviderCatalogResponse struct {
	Providers   []Provider `json:"providers"`
	TotalCount  int        `json:"total_count"`
	LastUpdated string     `json:"last_updated"`
}

type AnalyzeRequest struct {
	ImageBase64 string `json:"image_base64"`
	MediaType   string `json:"media_type"`
	RequestID   string `json:"request_id,omitempty"`
}

type AnalysisResult struct {
	RequestID         string           `json:"request_id"`
	HasMatch          bool             `json:"has_match"`
	MatchedProvider   *Provider        `json:"matched_provider,omitempty"`
	Justification     string           `json:"justification"`
	ConfidenceLevel   string           `json:"confidence_level"`
	ExtractedEntities []string         `json:"extracted_entities,omitempty"`
	ExtractedInfo     *ExtractedInfo   `json:"extracted_info,omitempty"`
	ProcessingTime    string           `json:"processing_time"`
	ModelUsed         string           `json:"model_used"`
	Latency           map[string]int64 `json:"latency,omitempty"`
	Timestamp         time.Time        `json:"timestamp"`
}

type ExtractedInfo struct {
	ProviderName  string `json:"provider_name,omitempty"`
	ServiceType   string `json:"service_type,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	TotalAmount   string `json:"total_amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	ReceiptDate   string `json:"receipt_date,omitempty"`
	DueDate       string `json:"due_date,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
}

type APIResponse struct {
	Success bool            `json:"success"`
	Data    *AnalysisResult `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type MCPRequest struct {
	Action      string `json:"action"`
	ServiceType string `json:"service_type,omitempty"`
	Region      string `json:"region,omitempty"`
	Status      string `json:"status,omitempty"`
}
