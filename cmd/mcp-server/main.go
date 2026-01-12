package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSdk "github.com/aws/aws-sdk-go-v2/service/lambda"
)

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
	AcceptedPaymentField   string   `json:"accepted_payment_field"`
	AdditionalInformation  string   `json:"additional_information,omitempty"`
}

type ProviderCatalogResponse struct {
	Providers   []Provider `json:"providers"`
	TotalCount  int        `json:"total_count"`
	LastUpdated string     `json:"last_updated"`
}

type MCPRequest struct {
	Action      string `json:"action"`
	ServiceType string `json:"service_type,omitempty"`
	Region      string `json:"region,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CatalogLambdaRequest struct {
	ServiceType string `json:"service_type,omitempty"`
	Region      string `json:"region,omitempty"`
	Status      string `json:"status,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

type MCPProviderServer struct {
	CatalogFunctionName string
	LambdaClient        *lambdaSdk.Client
}

func NewMCPProviderServer(ctx context.Context) (*MCPProviderServer, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &MCPProviderServer{
		CatalogFunctionName: os.Getenv("CATALOG_FUNCTION_NAME"),
		LambdaClient:        lambdaSdk.NewFromConfig(cfg),
	}, nil
}

func (s *MCPProviderServer) GetAllProviders(ctx context.Context, serviceType, region, status string) (*ProviderCatalogResponse, error) {
	if status == "" {
		status = "active"
	}

	// Build request for catalog Lambda
	catalogRequest := CatalogLambdaRequest{
		ServiceType: serviceType,
		Region:      region,
		Status:      status,
		Limit:       1000,
	}

	payload, err := json.Marshal(catalogRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal catalog request: %w", err)
	}

	log.Printf("Invoking catalog Lambda: %s", s.CatalogFunctionName)

	// Invoke the catalog Lambda directly
	result, err := s.LambdaClient.Invoke(ctx, &lambdaSdk.InvokeInput{
		FunctionName: &s.CatalogFunctionName,
		Payload:      payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke catalog Lambda: %w", err)
	}

	if result.FunctionError != nil {
		return nil, fmt.Errorf("catalog Lambda error: %s", *result.FunctionError)
	}

	var providerResponse ProviderCatalogResponse
	if err := json.Unmarshal(result.Payload, &providerResponse); err != nil {
		return nil, fmt.Errorf("failed to decode catalog response: %w", err)
	}

	log.Printf("Fetched %d providers from catalog Lambda", len(providerResponse.Providers))
	return &providerResponse, nil
}

func Handler(ctx context.Context, request MCPRequest) (string, error) {
	server, err := NewMCPProviderServer(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create MCP server: %w", err)
	}

	serviceType := request.ServiceType
	region := request.Region
	status := request.Status

	if status == "" {
		status = "active"
	}

	providers, err := server.GetAllProviders(ctx, serviceType, region, status)
	if err != nil {
		log.Printf("Error fetching providers: %v", err)
		return "", err
	}

	jsonProviders, err := json.Marshal(providers.Providers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal providers: %w", err)
	}

	return string(jsonProviders), nil
}

func main() {
	lambda.Start(Handler)
}
