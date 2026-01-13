# POC Vision API

API serverless para identificacion automatica de proveedores de servicios mediante analisis de imagenes de recibos/facturas usando modelos multimodales de IA.

## Tabla de Contenidos

1. [Descripcion General](#descripcion-general)
2. [Arquitectura](#arquitectura)
3. [Estructura del Proyecto](#estructura-del-proyecto)
4. [Requerimientos](#requerimientos)
5. [Configuracion](#configuracion)
6. [Proveedores de IA Soportados](#proveedores-de-ia-soportados)
7. [Instalacion y Despliegue](#instalacion-y-despliegue)
8. [Uso de la API](#uso-de-la-api)
9. [Ejemplos de Consumo](#ejemplos-de-consumo)
10. [Desarrollo Local](#desarrollo-local)
11. [Referencia de Configuracion](#referencia-de-configuracion)

---

## Descripcion General

Esta API permite analizar imagenes de recibos y facturas de servicios (electricidad, agua, gas, telecomunicaciones) para:

- Identificar automaticamente el proveedor del servicio
- Extraer informacion relevante como numero de referencia, cuenta, monto, fecha de vencimiento
- Comparar la informacion extraida con un catalogo de proveedores registrados
- Retornar el nivel de confianza del match encontrado

La solucion utiliza modelos multimodales con capacidades de vision y tool use para realizar el analisis.

---

## Arquitectura

El sistema esta compuesto por tres funciones Lambda que trabajan en conjunto:

```
                                    +------------------+
                                    |   API Gateway    |
                                    |  (x-api-key)     |
                                    +--------+---------+
                                             |
                    +------------------------+------------------------+
                    |                                                 |
                    v                                                 v
        +-----------+-----------+                       +-------------+-----------+
        |  Vision Analysis      |                       |    Catalog API          |
        |  Lambda               |                       |    Lambda               |
        |  (main-lambda)        |                       |    (catalog-api)        |
        +-----------+-----------+                       +-------------------------+
                    |                                                 ^
                    | Tool Call                                       |
                    v                                                 |
        +-----------+-----------+                                     |
        |    MCP Server         +-------------------------------------+
        |    Lambda             |  Invoke
        |    (mcp-server)       |
        +-----------+-----------+
                    |
                    v
        +-----------+-----------+
        |    DynamoDB           |
        |    (Results)          |
        +-----------------------+
```

### Componentes

| Componente | Descripcion |
|------------|-------------|
| **Vision Analysis Lambda** | Funcion principal que recibe imagenes, invoca el modelo de IA y coordina el analisis |
| **MCP Server Lambda** | Servidor de herramientas (Tool Use) que provee el catalogo de proveedores al modelo |
| **Catalog API Lambda** | API REST que expone el catalogo de proveedores de servicios |
| **DynamoDB** | Almacena los resultados de cada analisis para auditoria y trazabilidad |
| **API Gateway** | Punto de entrada HTTPS con autenticacion via API Key |

---

## Estructura del Proyecto

```
poc-vision-api/
├── cmd/                          # Puntos de entrada de las aplicaciones
│   ├── main-lambda/              # Lambda principal de analisis de vision
│   │   └── main.go
│   ├── mcp-server/               # Lambda servidor MCP (Tool Use)
│   │   └── main.go
│   └── catalog-api/              # Lambda API de catalogo de proveedores
│       └── main.go
├── internal/                     # Codigo interno del proyecto
│   ├── ai/                       # Abstraccion de proveedores de IA
│   │   ├── factory.go            # Fabrica de proveedores (patron Factory)
│   │   ├── types.go              # Interfaces y tipos compartidos
│   │   ├── utils.go              # Utilidades y prompts por defecto
│   │   ├── anthropic/            # Implementacion proveedor Anthropic
│   │   │   └── provider.go
│   │   ├── bedrock/              # Implementacion proveedor AWS Bedrock
│   │   │   └── provider.go
│   │   ├── google/               # Implementacion proveedor Google Gemini
│   │   │   └── provider.go
│   │   └── openai/               # Implementacion proveedor OpenAI
│   │       └── provider.go
│   └── models/                   # Modelos de datos compartidos
│       └── types.go
├── deployments/                  # Configuracion de infraestructura
│   └── template.yaml             # Template SAM/CloudFormation
├── scripts/                      # Scripts de build y deploy
│   ├── build.ps1                 # Build para Windows (PowerShell)
│   ├── build.sh                  # Build para Linux/macOS
│   ├── deploy.ps1                # Deploy para Windows (PowerShell)
│   ├── deploy.sh                 # Deploy para Linux/macOS
│   └── build_deploy.ps1          # Build + Deploy combinado (Windows)
├── build/                        # Directorio de artefactos compilados
├── config.env.example            # Ejemplo de archivo de configuracion
├── go.mod                        # Dependencias Go
└── README.md                     # Esta documentacion
```

### Descripcion de Directorios

| Directorio | Proposito |
|------------|-----------|
| `cmd/` | Contiene los puntos de entrada (main.go) de cada Lambda |
| `internal/ai/` | Implementacion del patron Factory para proveedores de IA con soporte multi-proveedor |
| `internal/models/` | Estructuras de datos compartidas entre componentes |
| `deployments/` | Templates de infraestructura como codigo (IaC) |
| `scripts/` | Automatizacion de build y despliegue |
| `build/` | Artefactos compilados (.zip) listos para desplegar |

---

## Requerimientos

### Software

| Herramienta | Version Minima | Proposito |
|-------------|----------------|-----------|
| Go | 1.21+ | Compilacion del codigo |
| AWS CLI | 2.x | Interaccion con servicios AWS |
| AWS SAM CLI | 1.x | Despliegue de infraestructura serverless |
| zip | - | Empaquetado de Lambdas (Linux/macOS) |

### Credenciales AWS

- Cuenta AWS con permisos para:
  - Lambda (crear, actualizar, invocar)
  - API Gateway (crear, configurar)
  - DynamoDB (crear tablas, PutItem, GetItem)
  - CloudFormation (crear stacks)
  - S3 (para artefactos de SAM)
  - IAM (crear roles)
  - Bedrock (InvokeModel, Converse) - si usa Bedrock

### Acceso a Modelos (segun proveedor)

| Proveedor | Requisito |
|-----------|-----------|
| AWS Bedrock | Solicitar acceso a modelos en la consola de Bedrock |
| OpenAI | API Key valida con acceso a modelos GPT-4 Vision |
| Anthropic | API Key valida con acceso a Claude 3+ |
| Google | API Key de Gemini API |

---

## Configuracion

### Archivo de Configuracion

Crea un archivo `config.env` en la raiz del proyecto basandote en `config.env.example`:

```bash
cp config.env.example config.env
```

### Variables de Entorno

| Variable | Descripcion | Valores Permitidos | Valor por Defecto |
|----------|-------------|-------------------|-------------------|
| `ENVIRONMENT` | Entorno de despliegue | `dev`, `staging`, `prod` | `dev` |
| `AWS_REGION` | Region de AWS | Cualquier region valida | `us-east-1` |
| `AI_PROVIDER` | Proveedor de IA a utilizar | `bedrock`, `openai`, `anthropic`, `google` | `bedrock` |
| `BEDROCK_MODEL_ID` | ID del modelo en Bedrock | Ver seccion de modelos | `amazon.nova-lite-v1:0` |
| `OPENAI_API_KEY` | API Key de OpenAI | - | (vacio) |
| `ANTHROPIC_API_KEY` | API Key de Anthropic | - | (vacio) |
| `GOOGLE_API_KEY` | API Key de Google | - | (vacio) |
| `CATALOG_API_TOKEN` | Token de autenticacion para la API | Cualquier string | `poc-vision-api-token-2026` |

### Ejemplo de Configuracion

**Usando AWS Bedrock (recomendado para produccion):**

```env
ENVIRONMENT=dev
AWS_REGION=us-east-1
AI_PROVIDER=bedrock
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
CATALOG_API_TOKEN=mi-token-seguro-123
```

**Usando OpenAI:**

```env
ENVIRONMENT=dev
AWS_REGION=us-east-1
AI_PROVIDER=openai
OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxx
CATALOG_API_TOKEN=mi-token-seguro-123
```

**Usando Anthropic:**

```env
ENVIRONMENT=dev
AWS_REGION=us-east-1
AI_PROVIDER=anthropic
ANTHROPIC_API_KEY=sk-ant-xxxxxxxxxxxxxxxxx
CATALOG_API_TOKEN=mi-token-seguro-123
```

**Usando Google Gemini:**

```env
ENVIRONMENT=dev
AWS_REGION=us-east-1
AI_PROVIDER=google
GOOGLE_API_KEY=AIzaxxxxxxxxxxxxxxxxxxxxxxx
CATALOG_API_TOKEN=mi-token-seguro-123
```

---

## Proveedores de IA Soportados

El proyecto implementa una arquitectura de proveedores intercambiables que permite usar diferentes modelos multimodales.

### Modelos Disponibles

| Proveedor | Modelo | ID | Caracteristicas |
|-----------|--------|----|-----------------| 
| **AWS Bedrock** | Amazon Nova Lite | `amazon.nova-lite-v1:0` | Rapido, economico, recomendado para dev |
| **AWS Bedrock** | Amazon Nova Pro | `amazon.nova-pro-v1:0` | Balance velocidad/calidad |
| **AWS Bedrock** | Claude 3.5 Sonnet | `anthropic.claude-3-5-sonnet-20241022-v2:0` | Alta precision |
| **OpenAI** | GPT-4o | `gpt-4o` | Modelo insignia de OpenAI |
| **OpenAI** | GPT-4 Turbo | `gpt-4-turbo` | Version Turbo optimizada |
| **Anthropic** | Claude Sonnet 4 | `claude-sonnet-4-20250514` | Ultima version de Claude |
| **Anthropic** | Claude 3.5 Sonnet | `claude-3-5-sonnet-20241022` | Version estable |
| **Google** | Gemini 2.0 Flash | `gemini-2.0-flash` | Rapido y eficiente |
| **Google** | Gemini 1.5 Pro | `gemini-1.5-pro` | Mayor capacidad |

### Capacidades por Proveedor

Todos los proveedores implementados soportan:

- Vision (analisis de imagenes)
- Tool Use / Function Calling (llamadas a herramientas)

### Arquitectura de Proveedores

El proyecto usa el patron Factory para registrar proveedores automaticamente:

```go
// Los proveedores se registran via init() al importarlos
import (
    _ "poc-vision-api/internal/ai/bedrock"
    _ "poc-vision-api/internal/ai/openai"
    _ "poc-vision-api/internal/ai/anthropic"
    _ "poc-vision-api/internal/ai/google"
)

// Crear proveedor desde variables de entorno
provider, err := ai.CreateProviderFromEnv(ctx)
```

---

## Instalacion y Despliegue

### Compilacion

**Windows (PowerShell):**

```powershell
.\scripts\build.ps1
```

**Linux/macOS:**

```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

Los artefactos compilados se generan en el directorio `build/`:

- `main-lambda.zip` - Lambda de analisis de vision
- `mcp-server.zip` - Lambda del servidor MCP
- `catalog-api.zip` - Lambda del catalogo de proveedores

### Despliegue

**Windows (PowerShell):**

```powershell
# Solo deploy (requiere build previo)
.\scripts\deploy.ps1

# Build + Deploy en un solo paso
.\scripts\build_deploy.ps1
```

**Linux/macOS:**

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

### Verificar Despliegue

Al finalizar el despliegue, se mostraran los endpoints y credenciales:

```
Deployment completed successfully
==================================
Analyze Endpoint: https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze
Catalog Endpoint: https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/catalog/providers
API Key: poc-vision-api-token-2026
Bedrock Model: amazon.nova-lite-v1:0
```

---

## Uso de la API

### Endpoints Disponibles

| Metodo | Endpoint | Descripcion | Autenticacion |
|--------|----------|-------------|---------------|
| POST | `/ai/analyze` | Analiza una imagen de recibo | API Key requerida |
| GET | `/catalog/providers` | Lista el catalogo de proveedores | API Key requerida |

### Autenticacion

Todas las peticiones requieren el header `x-api-key` con el token configurado:

```
x-api-key: poc-vision-api-token-2026
```

### POST /ai/analyze

Analiza una imagen de recibo y retorna la identificacion del proveedor.

**Request (JSON):**

```json
{
    "image_base64": "<imagen codificada en base64>",
    "media_type": "image/jpeg",
    "request_id": "opcional-uuid-personalizado"
}
```

**Request (multipart/form-data):**

- `image` o `file`: Archivo de imagen o PDF
- `request_id` (opcional): ID de solicitud personalizado

**Formatos soportados:**

- `image/jpeg`
- `image/png`
- `image/gif`
- `image/webp`
- `application/pdf`

**Response exitosa:**

```json
{
    "success": true,
    "data": {
        "request_id": "550e8400-e29b-41d4-a716-446655440000",
        "has_match": true,
        "matched_provider": {
            "id": "cfe-001",
            "name": "CFE",
            "official_name": "Comision Federal de Electricidad",
            "common_names": ["CFE", "Comision Federal de Electricidad", "Luz"],
            "aliases": ["Luz CFE", "CFE Luz", "Electricidad CFE"],
            "service_type": "electricidad",
            "region": "nacional",
            "status": "active",
            "online_service_available": true,
            "accepted_payment_field": "referencia",
            "additional_information": "inicia con 01 y tiene una longitud de 30 digitos..."
        },
        "justification": "Se identifico el logotipo de CFE y el formato de recibo de luz...",
        "confidence_level": "HIGH",
        "extracted_info": {
            "provider_name": "CFE",
            "service_type": "electricidad",
            "payment_field_name": "referencia",
            "payment_field_value": "012301215188912511200000005005",
            "total_amount": "500.00",
            "currency": "MXN",
            "due_date": "2026-02-15"
        },
        "processing_time": "2.5s",
        "model_used": "amazon.nova-lite-v1:0",
        "provider_name": "AWS Bedrock",
        "tool_calls": 1
    }
}
```

**Response de error:**

```json
{
    "success": false,
    "error": "invalid base64 image data"
}
```

### GET /catalog/providers

Lista todos los proveedores de servicios registrados.

**Response:**

```json
{
    "providers": [
        {
            "id": "cfe-001",
            "name": "CFE",
            "official_name": "Comision Federal de Electricidad",
            "service_type": "electricidad",
            "region": "nacional",
            "status": "active",
            "online_service_available": true,
            "accepted_payment_field": "referencia"
        },
        {
            "id": "telmex-001",
            "name": "Telmex",
            "official_name": "Telefonos de Mexico",
            "service_type": "telecomunicaciones",
            "region": "nacional",
            "status": "active",
            "online_service_available": true,
            "accepted_payment_field": "numero de cuenta"
        }
    ],
    "total_count": 7,
    "last_updated": "2026-01-12T10:30:00Z"
}
```

---

## Ejemplos de Consumo

### cURL - Enviar imagen como JSON (base64)

```bash
# Codificar imagen a base64
IMAGE_BASE64=$(base64 -w 0 recibo.jpg)

# Enviar peticion
curl -X POST "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze" \
    -H "x-api-key: poc-vision-api-token-2026" \
    -H "Content-Type: application/json" \
    -d "{\"image_base64\":\"$IMAGE_BASE64\",\"media_type\":\"image/jpeg\"}"
```

### cURL - Enviar imagen como multipart/form-data

```bash
curl -X POST "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze" \
    -H "x-api-key: poc-vision-api-token-2026" \
    -F "image=@recibo.jpg"
```

### cURL - Enviar PDF como multipart/form-data

```bash
curl -X POST "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze" \
    -H "x-api-key: poc-vision-api-token-2026" \
    -F "file=@factura.pdf"
```

### PowerShell

```powershell
$endpoint = "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze"
$apiKey = "poc-vision-api-token-2026"

# Leer y codificar imagen
$imageBytes = [System.IO.File]::ReadAllBytes("recibo.jpg")
$imageBase64 = [System.Convert]::ToBase64String($imageBytes)

# Construir body
$body = @{
    image_base64 = $imageBase64
    media_type = "image/jpeg"
} | ConvertTo-Json

# Enviar peticion
$response = Invoke-RestMethod -Uri $endpoint -Method POST `
    -Headers @{ "x-api-key" = $apiKey; "Content-Type" = "application/json" } `
    -Body $body

$response | ConvertTo-Json -Depth 10
```

### Python

```python
import base64
import requests

endpoint = "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze"
api_key = "poc-vision-api-token-2026"

# Leer y codificar imagen
with open("recibo.jpg", "rb") as f:
    image_base64 = base64.b64encode(f.read()).decode("utf-8")

# Enviar peticion
response = requests.post(
    endpoint,
    headers={
        "x-api-key": api_key,
        "Content-Type": "application/json"
    },
    json={
        "image_base64": image_base64,
        "media_type": "image/jpeg"
    }
)

print(response.json())
```

### JavaScript/Node.js

```javascript
const fs = require('fs');

const endpoint = 'https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze';
const apiKey = 'poc-vision-api-token-2026';

// Leer y codificar imagen
const imageBase64 = fs.readFileSync('recibo.jpg').toString('base64');

// Enviar peticion
fetch(endpoint, {
    method: 'POST',
    headers: {
        'x-api-key': apiKey,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        image_base64: imageBase64,
        media_type: 'image/jpeg'
    })
})
.then(response => response.json())
.then(data => console.log(JSON.stringify(data, null, 2)));
```

### Go

```go
package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    endpoint := "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/api/ai/analyze"
    apiKey := "poc-vision-api-token-2026"

    // Leer y codificar imagen
    imageData, _ := os.ReadFile("recibo.jpg")
    imageBase64 := base64.StdEncoding.EncodeToString(imageData)

    // Construir request
    body, _ := json.Marshal(map[string]string{
        "image_base64": imageBase64,
        "media_type":   "image/jpeg",
    })

    req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    result, _ := io.ReadAll(resp.Body)
    fmt.Println(string(result))
}
```

---

## Desarrollo Local

### Ejecutar Tests

```bash
go test ./...
```

### Agregar Nuevo Proveedor de IA

1. Crear directorio en `internal/ai/<nuevo-proveedor>/`

2. Implementar la interfaz `ai.Provider`:

```go
package nuevoproveedor

import (
    "context"
    "poc-vision-api/internal/ai"
)

type Provider struct {
    // campos del proveedor
}

func New(ctx context.Context, cfg ai.ProviderConfig) (*Provider, error) {
    // inicializacion
}

func (p *Provider) Name() string {
    return "Nuevo Proveedor"
}

func (p *Provider) SupportsVision() bool {
    return true
}

func (p *Provider) SupportsToolUse() bool {
    return true
}

func (p *Provider) AnalyzeImageWithTools(ctx context.Context, request ai.AnalysisRequest) (*ai.AnalysisResponse, error) {
    // implementacion del analisis
}

// Registrar el proveedor automaticamente al importar
func init() {
    ai.RegisterProvider(ai.ProviderType("nuevoproveedor"), func(ctx context.Context, cfg ai.ProviderConfig) (ai.Provider, error) {
        return New(ctx, cfg)
    })
}
```

3. Importar el proveedor en `cmd/main-lambda/main.go`:

```go
import (
    _ "poc-vision-api/internal/ai/nuevoproveedor"
)
```

### Agregar Nuevo Proveedor al Catalogo

Editar el archivo `cmd/catalog-api/main.go` y agregar el nuevo proveedor en la funcion `getProviders()`:

```go
{
    ID:                     "nuevo-001",
    Name:                   "Nombre",
    OfficialName:           "Nombre Oficial Completo",
    CommonNames:            []string{"Nombre", "Alias1"},
    Aliases:                []string{"Alias2", "Alias3"},
    ServiceType:            "tipo_servicio",
    Region:                 "nacional",
    Status:                 "active",
    OnlineServiceAvailable: true,
    AcceptedPaymentField:   "numero de cuenta",
},
```

---

## Referencia de Configuracion

### Parametros del Template SAM

| Parametro | Tipo | Default | Descripcion |
|-----------|------|---------|-------------|
| `Environment` | String | `dev` | Entorno de despliegue |
| `CatalogApiToken` | String | `poc-vision-api-token-2026` | Token de autenticacion API |
| `BedrockModelId` | String | `amazon.nova-lite-v1:0` | Modelo multimodal a utilizar |
| `AIProvider` | String | `bedrock` | Proveedor de IA |
| `AIApiKey` | String | (vacio) | API Key para proveedores externos |

### Recursos AWS Creados

| Recurso | Nombre | Descripcion |
|---------|--------|-------------|
| DynamoDB Table | `poc-vision-api-results-{env}` | Almacena resultados de analisis |
| API Gateway | `poc-vision-api-{env}` | Gateway REST con API Key |
| Lambda | `poc-vision-api-{env}-analysis` | Funcion principal de analisis |
| Lambda | `poc-vision-api-{env}-mcp-server` | Servidor de herramientas MCP |
| Lambda | `poc-vision-api-{env}-catalog-api` | API de catalogo de proveedores |
| API Key | `poc-vision-api-key-{env}` | Clave de acceso a la API |
| Usage Plan | `poc-vision-api-usage-plan-{env}` | Plan de uso con limites |

### Limites de la API

| Limite | Valor |
|--------|-------|
| Requests por dia | 1,000 |
| Burst limit | 50 requests |
| Rate limit | 25 requests/segundo |
| Timeout Lambda | 180 segundos |
| Memoria Lambda | 2048 MB |

### Formatos de Imagen

| Formato | MIME Type | Soportado |
|---------|-----------|-----------|
| JPEG | `image/jpeg` | Si |
| PNG | `image/png` | Si |
| GIF | `image/gif` | Si |
| WebP | `image/webp` | Si |

---

## Licencia

Proyecto de prueba de concepto (POC). Uso interno.
