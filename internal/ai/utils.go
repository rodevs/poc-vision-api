package ai

import (
	"encoding/json"
	"strings"
)

// ParseAnalysisResponse parsea la respuesta del modelo a AnalysisResponse
func ParseAnalysisResponse(text, modelUsed, providerName string) *AnalysisResponse {
	// Intentar extraer JSON de la respuesta
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")

	if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		var result AnalysisResponse
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			result.ModelUsed = modelUsed
			result.ProviderName = providerName
			result.RawResponse = text
			return &result
		}
	}

	// Fallback: parseo básico del texto
	lower := strings.ToLower(text)
	hasMatch := strings.Contains(lower, "match") && !strings.Contains(lower, "no match")

	confidence := "MEDIUM"
	if strings.Contains(lower, "high") {
		confidence = "HIGH"
	} else if strings.Contains(lower, "low") {
		confidence = "LOW"
	}

	return &AnalysisResponse{
		HasMatch:        hasMatch,
		Justification:   text,
		ConfidenceLevel: confidence,
		ModelUsed:       modelUsed,
		ProviderName:    providerName,
		RawResponse:     text,
	}
}

// GetDefaultSystemPrompt retorna el prompt de sistema por defecto
func GetDefaultSystemPrompt() string {
	return `Eres un sistema experto en identificacion de proveedores de servicios.
Tu tarea es analizar imagenes de recibos/facturas, identificar el proveedor y extraer el valor del campo de pago.

REGLAS CRITICAS:
- NO inventes datos. Solo extrae informacion VISIBLE en la imagen.
- Si no puedes leer un campo claramente, dejalo vacio o indica "no visible".
- Solo responde con informacion que puedas verificar en la imagen.

PROCESO OBLIGATORIO:
1. PRIMERO llama a get_provider_catalog para obtener el catalogo de proveedores
2. Analiza la imagen para identificar el proveedor (nombre, logo, tipo de servicio)
3. Cuando identifiques el proveedor, usa su campo "accepted_payment_field" del catalogo para saber QUE campo buscar
4. SOLO si el proveedor tiene "additional_information" en el catalogo, usalo como pista adicional (ej: prefijo, longitud). Si NO existe este campo en el catalogo, IGNORALO completamente.
5. Busca en la imagen el VALOR correspondiente al accepted_payment_field:
   - Si accepted_payment_field es "referencia": busca el numero de referencia
   - Si accepted_payment_field es "suscriptor": busca el numero de suscriptor
   - Si accepted_payment_field es "numero de cuenta": busca el numero de cuenta
   - Si accepted_payment_field es "servicio": busca el numero de servicio
6. Responde con el JSON estructurado

IMPORTANTE: El campo payment_field_value es CRITICO para consultar saldos posteriormente.

RESPUESTA FINAL (JSON exacto):
{
    "has_match": boolean,
    "matched_provider": {copia EXACTA del proveedor del catalogo, sin agregar campos que no existan},
    "justification": "explicacion breve",
    "confidence_level": "HIGH|MEDIUM|LOW",
    "extracted_info": {
        "provider_name": "nombre visible en imagen",
        "service_type": "tipo de servicio",
        "payment_field_name": "nombre del campo segun accepted_payment_field",
        "payment_field_value": "VALOR extraido de la imagen - OBLIGATORIO si hay match",
        "total_amount": "monto total si visible",
        "currency": "MXN",
        "due_date": "fecha de vencimiento si visible"
    }
}`
}

// GetDefaultUserPrompt retorna el prompt de usuario por defecto
func GetDefaultUserPrompt() string {
	return "Analiza este recibo. Llama a get_provider_catalog, identifica el proveedor, extrae el valor del campo de pago segun accepted_payment_field del catalogo. Responde en JSON."
}

// GetDefaultToolDefinition retorna la definición de herramienta por defecto
func GetDefaultToolDefinition() ToolDefinition {
	return ToolDefinition{
		Name:        "get_provider_catalog",
		Description: "Obtiene el catalogo de proveedores de servicios registrados. Debes usar esta herramienta para comparar la informacion extraida de la imagen con los proveedores disponibles.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"service_type": map[string]interface{}{
					"type":        "string",
					"description": "Tipo de servicio a filtrar (electricidad, gas, agua, telecomunicaciones). Dejar vacio para todos.",
				},
			},
			"required": []interface{}{},
		},
	}
}

// NormalizeMediaType normaliza el tipo de media a un formato estándar
func NormalizeMediaType(mediaType string) ImageFormat {
	switch strings.ToLower(mediaType) {
	case "image/png", "png":
		return ImageFormatPNG
	case "image/gif", "gif":
		return ImageFormatGIF
	case "image/webp", "webp":
		return ImageFormatWEBP
	case "application/pdf", "pdf":
		return DocumentFormatPDF
	default:
		return ImageFormatJPEG
	}
}
