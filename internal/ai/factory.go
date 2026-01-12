package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ProviderType representa los tipos de proveedores soportados
type ProviderType string

const (
	ProviderBedrock   ProviderType = "bedrock"
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderGoogle    ProviderType = "google"
	ProviderAzure     ProviderType = "azure"
)

// ProviderFactory crea instancias de proveedores
type ProviderFactory struct {
	registry map[ProviderType]ProviderConstructor
}

// ProviderConstructor es una función que crea un proveedor
type ProviderConstructor func(ctx context.Context, config ProviderConfig) (Provider, error)

// NewProviderFactory crea una nueva fábrica de proveedores
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		registry: make(map[ProviderType]ProviderConstructor),
	}
}

// Register registra un constructor de proveedor
func (f *ProviderFactory) Register(providerType ProviderType, constructor ProviderConstructor) {
	f.registry[providerType] = constructor
}

// Create crea una instancia de proveedor según el tipo
func (f *ProviderFactory) Create(ctx context.Context, providerType ProviderType, config ProviderConfig) (Provider, error) {
	constructor, exists := f.registry[providerType]
	if !exists {
		return nil, fmt.Errorf("provider type '%s' not registered", providerType)
	}
	return constructor(ctx, config)
}

// CreateFromEnv crea un proveedor basado en variables de entorno
func (f *ProviderFactory) CreateFromEnv(ctx context.Context) (Provider, error) {
	providerType := ProviderType(strings.ToLower(os.Getenv("AI_PROVIDER")))
	if providerType == "" {
		providerType = ProviderBedrock // default
	}

	config := ProviderConfig{
		APIKey:     os.Getenv("AI_API_KEY"),
		Model:      os.Getenv("AI_MODEL"),
		Region:     os.Getenv("AWS_REGION"),
		MaxRetries: 3,
		Timeout:    120,
	}

	// Configuración específica por proveedor
	switch providerType {
	case ProviderBedrock:
		if config.Model == "" {
			config.Model = os.Getenv("BEDROCK_MODEL_ID")
		}
	case ProviderOpenAI:
		if config.APIKey == "" {
			config.APIKey = os.Getenv("OPENAI_API_KEY")
		}
		if config.Model == "" {
			config.Model = "gpt-4o"
		}
	case ProviderAnthropic:
		if config.APIKey == "" {
			config.APIKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		if config.Model == "" {
			config.Model = "claude-sonnet-4-20250514"
		}
	case ProviderGoogle:
		if config.APIKey == "" {
			config.APIKey = os.Getenv("GOOGLE_API_KEY")
		}
		if config.Model == "" {
			config.Model = "gemini-2.0-flash"
		}
	}

	return f.Create(ctx, providerType, config)
}

// ListRegistered retorna los tipos de proveedores registrados
func (f *ProviderFactory) ListRegistered() []ProviderType {
	types := make([]ProviderType, 0, len(f.registry))
	for t := range f.registry {
		types = append(types, t)
	}
	return types
}

// DefaultFactory es la fábrica global de proveedores
var DefaultFactory = NewProviderFactory()

// RegisterProvider registra un proveedor en la fábrica global
func RegisterProvider(providerType ProviderType, constructor ProviderConstructor) {
	DefaultFactory.Register(providerType, constructor)
}

// CreateProvider crea un proveedor usando la fábrica global
func CreateProvider(ctx context.Context, providerType ProviderType, config ProviderConfig) (Provider, error) {
	return DefaultFactory.Create(ctx, providerType, config)
}

// CreateProviderFromEnv crea un proveedor desde variables de entorno
func CreateProviderFromEnv(ctx context.Context) (Provider, error) {
	return DefaultFactory.CreateFromEnv(ctx)
}
