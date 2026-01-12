package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

// DirectInvokeRequest is used when Lambda is invoked directly (not via API Gateway)
type DirectInvokeRequest struct {
	ServiceType string `json:"service_type,omitempty"`
	Region      string `json:"region,omitempty"`
	Status      string `json:"status,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

func getProviders() ProviderCatalogResponse {
	return ProviderCatalogResponse{
		Providers: []Provider{
			{
				ID:                     "cfe-001",
				Name:                   "CFE",
				OfficialName:           "Comision Federal de Electricidad",
				CommonNames:            []string{"CFE", "Comision Federal de Electricidad", "Luz"},
				Aliases:                []string{"Luz CFE", "CFE Luz", "Electricidad CFE"},
				ServiceType:            "electricidad",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "referencia",
				AdditionalInformation:  "inicia con 01 y tiene una longitud de 30 digitos. La referencia puede estar separada por espacios por ejemplo: 01 230121518891 251120 000000500 5. En ese caso tomala como un solo valor y respondela sin espacios.",
			},
			{
				ID:                     "telmex-001",
				Name:                   "Telmex",
				OfficialName:           "Telefonos de Mexico",
				CommonNames:            []string{"Telmex", "Telefonos de Mexico"},
				Aliases:                []string{"Infinitum", "Telmex Internet", "Telmex Telefono"},
				ServiceType:            "telecomunicaciones",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "numero de cuenta",
			},
			{
				ID:                     "izzi-001",
				Name:                   "izzi",
				OfficialName:           "izzi Telecom",
				CommonNames:            []string{"izzi", "izzi Telecom"},
				Aliases:                []string{"izzi Internet", "Cablevision"},
				ServiceType:            "telecomunicaciones",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "servicio",
			},
			{
				ID:                     "totalplay-001",
				Name:                   "Totalplay",
				OfficialName:           "Totalplay Telecomunicaciones",
				CommonNames:            []string{"Totalplay", "Total Play"},
				Aliases:                []string{"Totalplay Internet", "Totalplay TV"},
				ServiceType:            "telecomunicaciones",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "numero de cuenta",
			},
			{
				ID:                     "megacable-001",
				Name:                   "Megacable",
				OfficialName:           "Megacable Comunicaciones",
				CommonNames:            []string{"Megacable"},
				Aliases:                []string{"Megacable Internet", "Megacable TV"},
				ServiceType:            "telecomunicaciones",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "suscriptor",
			},
			{
				ID:                     "naturgy-001",
				Name:                   "Naturgy",
				OfficialName:           "Naturgy Mexico",
				CommonNames:            []string{"Naturgy", "Gas Natural"},
				Aliases:                []string{"Gas Natural Fenosa", "Naturgy Gas"},
				ServiceType:            "gas",
				Region:                 "nacional",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "referencia",
			},
			{
				ID:                     "sacmex-001",
				Name:                   "SACMEX",
				OfficialName:           "Sistema de Aguas de la Ciudad de Mexico",
				CommonNames:            []string{"SACMEX", "Agua CDMX"},
				Aliases:                []string{"Agua Ciudad de Mexico", "Aguas CDMX"},
				ServiceType:            "agua",
				Region:                 "cdmx",
				Status:                 "active",
				OnlineServiceAvailable: true,
				AcceptedPaymentField:   "numero de cuenta",
			},
		},
		TotalCount:  7,
		LastUpdated: time.Now().UTC().Format(time.RFC3339),
	}
}

// Handler that supports both API Gateway and direct Lambda invocation
func handler(ctx context.Context, rawEvent json.RawMessage) (interface{}, error) {
	providers := getProviders()

	// Try to parse as API Gateway request
	var apiRequest events.APIGatewayProxyRequest
	if err := json.Unmarshal(rawEvent, &apiRequest); err == nil && apiRequest.HTTPMethod != "" {
		// This is an API Gateway request
		body, _ := json.Marshal(providers)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: string(body),
		}, nil
	}

	// Direct Lambda invocation - return the providers directly
	return providers, nil
}

func main() {
	lambda.Start(handler)
}
