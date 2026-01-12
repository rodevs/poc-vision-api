package bedrock

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"poc-vision-api/internal/ai"
)

// Provider implementa ai.Provider para AWS Bedrock
type Provider struct {
	client   *bedrockruntime.Client
	modelID  string
	region   string
	maxRetry int
}

// New crea una nueva instancia del proveedor Bedrock
func New(ctx context.Context, cfg ai.ProviderConfig) (*Provider, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRetryMode(aws.RetryModeStandard),
		config.WithRetryMaxAttempts(cfg.MaxRetries),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if cfg.Region != "" {
		awsCfg.Region = cfg.Region
	}

	client := bedrockruntime.NewFromConfig(awsCfg)

	modelID := cfg.Model
	if modelID == "" {
		modelID = "anthropic.claude-3-5-sonnet-20241022-v2:0"
	}

	return &Provider{
		client:   client,
		modelID:  modelID,
		region:   awsCfg.Region,
		maxRetry: cfg.MaxRetries,
	}, nil
}

func (p *Provider) Name() string {
	return "AWS Bedrock"
}

func (p *Provider) SupportsVision() bool {
	return true
}

func (p *Provider) SupportsToolUse() bool {
	return true
}

// AnalyzeImageWithTools implementa el análisis de imagen con Tool Use
func (p *Provider) AnalyzeImageWithTools(ctx context.Context, request ai.AnalysisRequest) (*ai.AnalysisResponse, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(request.ImageBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64: %w", err)
	}

	imageFormat := p.mapImageFormat(request.MediaType)

	systemPrompt := request.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = ai.GetDefaultSystemPrompt()
	}

	userPrompt := request.UserPrompt
	if userPrompt == "" {
		userPrompt = ai.GetDefaultUserPrompt()
	}

	messages := []types.Message{
		{
			Role: types.ConversationRoleUser,
			Content: []types.ContentBlock{
				&types.ContentBlockMemberImage{
					Value: types.ImageBlock{
						Format: imageFormat,
						Source: &types.ImageSourceMemberBytes{
							Value: imageBytes,
						},
					},
				},
				&types.ContentBlockMemberText{
					Value: userPrompt,
				},
			},
		},
	}

	toolConfig := p.buildToolConfig(request.Tools)
	toolCalls := 0
	maxIterations := 5

	maxTokens := int32(4096)
	if request.MaxTokens > 0 {
		maxTokens = int32(request.MaxTokens)
	}

	temperature := float32(0.0)
	if request.Temperature > 0 {
		temperature = request.Temperature
	}

	for i := 0; i < maxIterations; i++ {
		output, err := p.client.Converse(ctx, &bedrockruntime.ConverseInput{
			ModelId:    aws.String(p.modelID),
			Messages:   messages,
			System:     []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: systemPrompt}},
			ToolConfig: toolConfig,
			InferenceConfig: &types.InferenceConfiguration{
				MaxTokens:   aws.Int32(maxTokens),
				Temperature: aws.Float32(temperature),
				TopP:        aws.Float32(0.9),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("converse failed: %w", err)
		}

		// Respuesta final
		if output.StopReason == types.StopReasonEndTurn {
			responseText := p.extractTextFromOutput(output)
			result := ai.ParseAnalysisResponse(responseText, p.modelID, p.Name())
			result.ToolCallsCount = toolCalls
			return result, nil
		}

		// Tool Use solicitado
		if output.StopReason == types.StopReasonToolUse {
			assistantContent := output.Output.(*types.ConverseOutputMemberMessage).Value.Content
			messages = append(messages, types.Message{
				Role:    types.ConversationRoleAssistant,
				Content: assistantContent,
			})

			var toolResults []types.ContentBlock
			for _, block := range assistantContent {
				if toolUse, ok := block.(*types.ContentBlockMemberToolUse); ok {
					toolCalls++
					log.Printf("[Bedrock] Executing tool: %s", *toolUse.Value.Name)

					var toolInput map[string]interface{}
					if toolUse.Value.Input != nil {
						inputBytes, _ := json.Marshal(toolUse.Value.Input)
						json.Unmarshal(inputBytes, &toolInput)
					}

					// Ejecutar herramienta
					toolResult, err := request.ToolExecutor.ExecuteTool(ctx, *toolUse.Value.Name, toolInput)
					if err != nil {
						toolResults = append(toolResults, &types.ContentBlockMemberToolResult{
							Value: types.ToolResultBlock{
								ToolUseId: toolUse.Value.ToolUseId,
								Status:    types.ToolResultStatusError,
								Content: []types.ToolResultContentBlock{
									&types.ToolResultContentBlockMemberText{Value: err.Error()},
								},
							},
						})
					} else {
						toolResults = append(toolResults, &types.ContentBlockMemberToolResult{
							Value: types.ToolResultBlock{
								ToolUseId: toolUse.Value.ToolUseId,
								Status:    types.ToolResultStatusSuccess,
								Content: []types.ToolResultContentBlock{
									&types.ToolResultContentBlockMemberText{Value: toolResult},
								},
							},
						})
					}
				}
			}

			messages = append(messages, types.Message{
				Role:    types.ConversationRoleUser,
				Content: toolResults,
			})
			continue
		}

		// Otra razón de parada
		responseText := p.extractTextFromOutput(output)
		result := ai.ParseAnalysisResponse(responseText, p.modelID, p.Name())
		result.ToolCallsCount = toolCalls
		return result, nil
	}

	return nil, fmt.Errorf("max iterations reached")
}

func (p *Provider) buildToolConfig(tools []ai.ToolDefinition) *types.ToolConfiguration {
	if len(tools) == 0 {
		tools = []ai.ToolDefinition{ai.GetDefaultToolDefinition()}
	}

	var bedrockTools []types.Tool
	for _, tool := range tools {
		bedrockTools = append(bedrockTools, &types.ToolMemberToolSpec{
			Value: types.ToolSpecification{
				Name:        aws.String(tool.Name),
				Description: aws.String(tool.Description),
				InputSchema: &types.ToolInputSchemaMemberJson{
					Value: document.NewLazyDocument(tool.InputSchema),
				},
			},
		})
	}

	return &types.ToolConfiguration{
		Tools: bedrockTools,
	}
}

func (p *Provider) mapImageFormat(mediaType string) types.ImageFormat {
	format := ai.NormalizeMediaType(mediaType)
	switch format {
	case ai.ImageFormatPNG:
		return types.ImageFormatPng
	case ai.ImageFormatGIF:
		return types.ImageFormatGif
	case ai.ImageFormatWEBP:
		return types.ImageFormatWebp
	default:
		return types.ImageFormatJpeg
	}
}

func (p *Provider) extractTextFromOutput(output *bedrockruntime.ConverseOutput) string {
	if msg, ok := output.Output.(*types.ConverseOutputMemberMessage); ok {
		for _, block := range msg.Value.Content {
			if textBlock, ok := block.(*types.ContentBlockMemberText); ok {
				return textBlock.Value
			}
		}
	}
	return ""
}

// init registra el proveedor en la fábrica global
func init() {
	ai.RegisterProvider(ai.ProviderBedrock, func(ctx context.Context, config ai.ProviderConfig) (ai.Provider, error) {
		return New(ctx, config)
	})
}
