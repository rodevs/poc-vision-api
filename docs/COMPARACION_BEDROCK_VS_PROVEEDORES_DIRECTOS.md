# AWS Bedrock vs Proveedores Directos: Análisis Ejecutivo

**Documento actualizado:** 13 de enero de 2026  
**Tipo de cambio utilizado:** 1 USD = 17.85 MXN

---

## Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Modelos Comparables](#modelos-comparables)
3. [Comparativa de Precios: Anthropic Claude](#comparativa-de-precios-anthropic-claude)
4. [Comparativa de Precios: Google Gemini](#comparativa-de-precios-google-gemini)
5. [Comparativa de Precios: OpenAI](#comparativa-de-precios-openai)
6. [Comparativa de Precios: Mistral AI](#comparativa-de-precios-mistral-ai)
7. [Comparativa de Precios: Meta Llama](#comparativa-de-precios-meta-llama)
8. [Características y Servicios Adicionales](#características-y-servicios-adicionales)
9. [Matriz de Decisión por Caso de Uso](#matriz-de-decisión-por-caso-de-uso)
10. [Conclusiones y Recomendaciones](#conclusiones-y-recomendaciones)

---

## Resumen Ejecutivo

### Decisión Rápida por Proveedor

| Proveedor | ¿Cuándo usar Bedrock? | ¿Cuándo usar API Directa? |
|-----------|----------------------|---------------------------|
| **Anthropic Claude** | Ecosistema AWS, compliance enterprise, multi-modelo | Precios ligeramente menores, prompt caching nativo |
| **Google Gemini** | Integración AWS existente, sin cuenta GCP | Tier gratuito, Google Search grounding, precios menores |
| **OpenAI** | Solo modelos open-weight disponibles en Bedrock | Modelos propietarios (GPT-5.2, o3), fine-tuning |
| **Mistral AI** | Infraestructura AWS centralizada | Precios equivalentes, prefer si ya usa Mistral directamente |
| **Meta Llama** | Sin gestión de infraestructura | Self-hosting para máximo control y menor costo |

### Ahorro Estimado por Proveedor (API Directa vs Bedrock)

| Proveedor | Diferencia de Precio | Observación |
|-----------|---------------------|-------------|
| Google Gemini | **30-50% más barato directo** | Tier gratuito disponible |
| Anthropic Claude | **~Similar o 5-10% menos directo** | Prompt caching mejor en directo |
| OpenAI | **N/A - Modelos diferentes** | Bedrock solo tiene open-weight |
| Mistral AI | **~Equivalente** | Mínima diferencia |
| Meta Llama | **~Equivalente** | Open source, opción self-host |

---

## Modelos Comparables

Solo se incluyen modelos disponibles **tanto en AWS Bedrock como en API directa del proveedor**:

| Proveedor | Modelos Comparables |
|-----------|-------------------|
| **Anthropic** | Claude Opus 4.5, Claude Sonnet 4.5, Claude Haiku 4.5, Claude 3.7 Sonnet |
| **Google** | Gemini 2.5 Pro, Gemini 2.5 Flash, Gemini 2.0 Flash |
| **Meta** | Llama 3.1 (8B, 70B, 405B) |
| **Mistral** | Mistral Large, Mixtral 8x7B, Mistral 7B |

> **Nota:** OpenAI solo ofrece modelos open-weight (GPT-OSS) en Bedrock. Los modelos propietarios como GPT-5.2 y o3 NO están disponibles en Bedrock.

---

## Comparativa de Precios: Anthropic Claude

### Precios por 1M Tokens (USD)

| Modelo | Bedrock Input | Directo Input | Bedrock Output | Directo Output |
|--------|---------------|---------------|----------------|----------------|
| **Claude Opus 4.5** | $5.00 | $5.00 | $25.00 | $25.00 |
| **Claude Sonnet 4.5** (≤200K) | $3.00 | $3.00 | $15.00 | $15.00 |
| **Claude Sonnet 4.5** (>200K) | $6.00 | $6.00 | $22.50 | $22.50 |
| **Claude Haiku 4.5** | $1.00 | $1.00 | $5.00 | $5.00 |

### Prompt Caching (USD/MTok)

| Característica | Bedrock | API Directa Anthropic |
|----------------|---------|----------------------|
| Cache Read (Opus 4.5) | $0.50 | $0.50 |
| Cache Write (Opus 4.5) | $6.25 | $6.25 |
| Cache Read (Sonnet 4.5) | $0.30 | $0.30 |
| Cache Write (Sonnet 4.5) | $3.75 | $3.75 |
| Cache Read (Haiku 4.5) | $0.10 | $0.10 |
| Cache Write (Haiku 4.5) | $1.25 | $1.25 |

### Veredicto Claude

| Criterio | Ganador | Observación |
|----------|---------|-------------|
| **Precio base** | 🟰 Empate | Precios idénticos |
| **Prompt Caching** | 🟰 Empate | Disponible en ambos, mismos precios |
| **Batch Processing** | ✅ Bedrock | 50% descuento, 24h processing |
| **Rate Limits** | ✅ API Directa | Mayor flexibilidad para aumentos |
| **Latencia** | 🟰 Empate | Similar en ambas plataformas |

**💡 Recomendación:** Para Claude, la decisión depende del ecosistema. Si ya estás en AWS, usa Bedrock. Si necesitas máxima flexibilidad en rate limits o usas múltiples clouds, considera API directa.

---

## Comparativa de Precios: Google Gemini

### Precios por 1M Tokens (USD)

| Modelo | Bedrock Input | Directo Input | Bedrock Output | Directo Output |
|--------|---------------|---------------|----------------|----------------|
| **Gemini 3 Pro** (≤200K) | N/D | $2.00 | N/D | $12.00 |
| **Gemini 3 Flash** | N/D | $0.50 | N/D | $3.00 |
| **Gemini 2.5 Pro** (≤200K) | $1.25* | $1.25 | $10.00* | $10.00 |
| **Gemini 2.5 Flash** | $0.30* | $0.30 | $2.50* | $2.50 |
| **Gemini 2.0 Flash** | $0.10* | $0.10 | $0.40* | $0.40 |

*Precios estimados basados en disponibilidad regional de Bedrock

### Tier Gratuito (Solo API Directa Google)

| Modelo | Límite Gratuito | Valor Estimado/Mes |
|--------|----------------|-------------------|
| Gemini 3 Flash | Standard tier ilimitado | ~$50-100 USD |
| Gemini 2.5 Pro | Standard tier ilimitado | ~$100-200 USD |
| Gemini 2.5 Flash | Standard tier ilimitado | ~$50-100 USD |

### Herramientas Exclusivas API Directa Google

| Herramienta | Disponible Bedrock | Disponible Directo | Precio Directo |
|-------------|-------------------|-------------------|----------------|
| Google Search Grounding | ❌ No | ✅ Sí | 1,500 RPD gratis, luego $35/1K |
| Google Maps Grounding | ❌ No | ✅ Sí | 1,500 RPD gratis, luego $25/1K |
| Deep Research Agent | ❌ No | ✅ Sí | Tokens a precio Gemini 3 Pro |
| Code Execution | ❌ No | ✅ Sí | Gratis |
| Context Caching Storage | ✅ Limitado | ✅ Sí | $1-4.50/1M tokens/hora |

### Veredicto Gemini

| Criterio | Ganador | Observación |
|----------|---------|-------------|
| **Precio base** | 🟰 Empate | Precios similares cuando disponible |
| **Tier gratuito** | ✅ API Directa | Tier gratuito generoso |
| **Herramientas avanzadas** | ✅ API Directa | Search/Maps grounding exclusivo |
| **Disponibilidad modelos** | ✅ API Directa | Gemini 3 primero en directa |
| **Integración AWS** | ✅ Bedrock | Sin salir del ecosistema AWS |

**💡 Recomendación:** Para Gemini, **API directa es claramente superior** por el tier gratuito, acceso a modelos más recientes (Gemini 3), y herramientas exclusivas como Google Search grounding.

---

## Comparativa de Precios: OpenAI

### Modelos en AWS Bedrock (Solo Open-Weight)

| Modelo Bedrock | Tipo | Input/1K | Output/1K |
|----------------|------|----------|-----------|
| GPT-OSS-120B | Open-weight | ~$0.005 | ~$0.016 |
| GPT-OSS-20B | Open-weight | ~$0.001 | ~$0.003 |
| GPT-OSS-Safeguard | Moderación | ~$0.001 | ~$0.002 |

### Modelos API Directa OpenAI (Propietarios)

| Modelo | Input/1M | Output/1M | Notas |
|--------|----------|-----------|-------|
| **GPT-5.2** | $1.75 | $14.00 | Flagship, mejor para agentes |
| **GPT-5.2 Pro** | $21.00 | $168.00 | Máxima precisión |
| **GPT-5 mini** | $0.25 | $2.00 | Balance costo-rendimiento |
| **GPT-4.1** | $3.00 | $12.00 | Modelo anterior, estable |
| **GPT-4.1 mini** | $0.80 | $3.20 | Económico |
| **o4-mini** | $4.00 | $16.00 | Razonamiento avanzado |

### Veredicto OpenAI

| Criterio | Ganador | Observación |
|----------|---------|-------------|
| **Modelos propietarios** | ✅ API Directa | GPT-5.2, o3, o4 solo en directa |
| **Modelos open-weight** | ✅ Bedrock | Gestión simplificada |
| **Fine-tuning** | ✅ API Directa | Disponible para GPT-4.1 familia |
| **Herramientas (Web Search)** | ✅ API Directa | $10-25/1K calls |

**💡 Recomendación:** **No son comparables directamente.** Si necesitas GPT-5.2 o modelos o3/o4, debes usar API directa OpenAI. Bedrock solo ofrece versiones open-weight.

---

## Comparativa de Precios: Mistral AI

### Precios por 1K Tokens (USD)

| Modelo | Bedrock Input | Directo Input | Bedrock Output | Directo Output |
|--------|---------------|---------------|----------------|----------------|
| **Mistral Large** | $0.008 | $0.008 | $0.024 | $0.024 |
| **Mixtral 8x7B** | $0.00045 | $0.0005 | $0.0007 | $0.0007 |
| **Mistral 7B** | $0.00015 | $0.00015 | $0.0002 | $0.0002 |

### Veredicto Mistral

| Criterio | Ganador | Observación |
|----------|---------|-------------|
| **Precio** | 🟰 Empate | Precios prácticamente idénticos |
| **Disponibilidad** | 🟰 Empate | Mismos modelos disponibles |
| **Latencia** | 🟰 Empate | Similar rendimiento |

**💡 Recomendación:** Para Mistral, elige según tu infraestructura existente. No hay ventaja significativa en ninguna plataforma.

---

## Comparativa de Precios: Meta Llama

### Precios por 1K Tokens (USD)

| Modelo | Bedrock Input | Self-Host | Bedrock Output | Self-Host |
|--------|---------------|-----------|----------------|-----------|
| **Llama 3.1 8B** | $0.00022 | $0* | $0.00022 | $0* |
| **Llama 3.1 70B** | $0.00099 | $0* | $0.00099 | $0* |
| **Llama 3.1 405B** | $0.00532 | $0* | $0.016 | $0* |

*Self-host solo paga infraestructura (GPU compute)

### Costo Estimado Self-Host vs Bedrock

| Modelo | Bedrock (100M tokens/mes) | Self-Host EC2 (estimado/mes) |
|--------|--------------------------|------------------------------|
| Llama 3.1 8B | ~$44 USD | ~$200-400 USD (g5.xlarge) |
| Llama 3.1 70B | ~$198 USD | ~$2,000-4,000 USD (p4d.24xlarge) |
| Llama 3.1 405B | ~$2,132 USD | ~$15,000-30,000 USD (p5.48xlarge) |

### Veredicto Llama

| Criterio | Ganador | Observación |
|----------|---------|-------------|
| **Costo bajo volumen** | ✅ Bedrock | Sin costo fijo de infraestructura |
| **Costo alto volumen** | ✅ Self-Host | Break-even ~500M tokens/mes |
| **Simplicidad** | ✅ Bedrock | Zero ops, escalado automático |
| **Control total** | ✅ Self-Host | Customización completa |

**💡 Recomendación:** Para Llama, usa **Bedrock para volumen bajo-medio** (<500M tokens/mes). Considera self-hosting solo para volumen muy alto o requerimientos especiales de personalización.

---

## Características y Servicios Adicionales

### Comparativa de Funcionalidades

| Característica | AWS Bedrock | API Directa |
|----------------|-------------|-------------|
| **API unificada multi-modelo** | ✅ Sí | ❌ No (cada proveedor diferente) |
| **Billing centralizado** | ✅ Sí (AWS) | ❌ Múltiples facturas |
| **Compliance enterprise** | ✅ SOC, ISO, HIPAA, FedRAMP | ⚠️ Varía por proveedor |
| **VPC endpoints** | ✅ Sí | ❌ No (tráfico público) |
| **Guardrails integrados** | ✅ Sí ($0.15-0.17/1K text units) | ⚠️ Limitado |
| **Knowledge Bases** | ✅ Sí (RAG nativo) | ❌ Implementar manualmente |
| **Agents framework** | ✅ AgentCore nativo | ❌ Implementar manualmente |
| **Batch Processing 50%** | ✅ Sí | ⚠️ Varía por proveedor |
| **Cross-Region Inference** | ✅ Sí (alta disponibilidad) | ❌ No |
| **Model switching sin código** | ✅ Sí | ❌ Cambios de SDK necesarios |

### Certificaciones de Seguridad

| Certificación | AWS Bedrock | Anthropic | Google | OpenAI |
|---------------|-------------|-----------|--------|--------|
| SOC 2 Type II | ✅ | ✅ | ✅ | ✅ |
| ISO 27001 | ✅ | ✅ | ✅ | ✅ |
| HIPAA Eligible | ✅ | ⚠️ BAA disponible | ✅ | ✅ |
| FedRAMP | ✅ High | ❌ | ⚠️ Moderate | ❌ |
| GDPR | ✅ | ✅ | ✅ | ✅ |

---

## Matriz de Decisión por Caso de Uso

### Escenario 1: Startup / MVP

| Factor | Recomendación | Proveedor |
|--------|---------------|-----------|
| **Prioridad:** Costo mínimo | Google Gemini API | Tier gratuito |
| **Backup:** Escalabilidad | AWS Bedrock Nova 2 Lite | $0.00006/1K input |
| **Ahorro estimado/mes** | $100-500 USD | Usando tier gratuito |

### Escenario 2: Empresa Mediana (AWS-native)

| Factor | Recomendación | Proveedor |
|--------|---------------|-----------|
| **Prioridad:** Simplicidad ops | AWS Bedrock | Multi-modelo unificado |
| **Modelo recomendado** | Claude Sonnet 4.5 + Nova 2 Lite | Balance rendimiento-costo |
| **Costo estimado/mes (100K req)** | $5,000-15,000 MXN | ~$280-840 USD |

### Escenario 3: Enterprise (Multi-cloud)

| Factor | Recomendación | Proveedor |
|--------|---------------|-----------|
| **Prioridad:** Redundancia | Híbrido Bedrock + APIs directas | Alta disponibilidad |
| **Configuración** | Bedrock primario + failover directo | 99.99% uptime |
| **Costo adicional failover** | 10-20% overhead | Justificado por SLA |

### Escenario 4: AI-First / Alto Volumen

| Factor | Recomendación | Proveedor |
|--------|---------------|-----------|
| **Prioridad:** Features avanzados | API Directa (Anthropic/OpenAI) | Últimas capacidades |
| **Prioridad:** Costo escala | Self-host Llama | Control total |
| **Prioridad:** Mix óptimo | Bedrock Priority Tier | SLA garantizado |

---

## Conclusiones y Recomendaciones

### ✅ Cuándo USAR AWS Bedrock

| Situación | Beneficio Principal |
|-----------|-------------------|
| **Infraestructura AWS existente** | Billing unificado, VPC integration, IAM nativo |
| **Requerimientos compliance estrictos** | FedRAMP High, HIPAA, SOC con un solo vendor |
| **Necesidad multi-modelo** | Cambiar entre Claude, Gemini, Llama sin cambios de código |
| **Equipo sin expertise ML** | Zero-ops, escalado automático |
| **Casos de uso enterprise** | Guardrails, Knowledge Bases, AgentCore integrados |
| **Procesamiento batch** | 50% descuento consistente en todos los modelos |

### ❌ Cuándo NO USAR AWS Bedrock

| Situación | Mejor Alternativa |
|-----------|-------------------|
| **Necesitas GPT-5.2 o modelos o3/o4** | OpenAI API directa (no disponibles en Bedrock) |
| **Requieres Google Search grounding** | Google AI API directa (exclusivo) |
| **Presupuesto muy limitado** | Google Gemini tier gratuito |
| **Alto volumen con modelo fijo** | Self-host Llama (break-even ~500M tokens/mes) |
| **Necesitas modelos más recientes** | API directa (Gemini 3 primero en Google) |
| **Ya tienes infraestructura GCP** | Vertex AI (mejor integración) |

### Tabla Resumen Final

| Proveedor | Usar Bedrock | Usar Directo | Empate |
|-----------|--------------|--------------|--------|
| **Anthropic Claude** | AWS-native, compliance | Rate limits, última versión | Precio |
| **Google Gemini** | Sin cuenta GCP | Tier gratis, Search grounding | - |
| **OpenAI** | Open-weight únicamente | Modelos propietarios | - |
| **Mistral AI** | AWS-native | Ya usa Mistral | Precio |
| **Meta Llama** | Volumen bajo-medio | Alto volumen, self-host | - |

### Pros y Contras Generales

#### AWS Bedrock

| Pros | Contras |
|------|---------|
| ✅ API unificada para todos los modelos | ❌ No todos los modelos disponibles |
| ✅ Integración nativa AWS (VPC, IAM, CloudWatch) | ❌ Modelos nuevos tardan en llegar |
| ✅ Compliance enterprise completo | ❌ Sin tier gratuito |
| ✅ Guardrails y Knowledge Bases incluidos | ❌ Vendor lock-in AWS |
| ✅ Batch processing 50% descuento | ❌ Sin herramientas exclusivas (ej. Google Search) |
| ✅ Cross-region inference | ❌ Precios a veces ligeramente mayores |

#### APIs Directas

| Pros | Contras |
|------|---------|
| ✅ Acceso a modelos más recientes | ❌ Múltiples SDKs y facturas |
| ✅ Features exclusivos por proveedor | ❌ Implementar seguridad manualmente |
| ✅ Tiers gratuitos (Google) | ❌ Sin API unificada |
| ✅ Mayor flexibilidad rate limits | ❌ Compliance case-by-case |
| ✅ Relación directa con proveedor | ❌ Sin guardrails nativos |

### Recomendación Estratégica Final

```
┌─────────────────────────────────────────────────────────────┐
│                    RECOMENDACIÓN GENERAL                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  🏢 ENTERPRISE AWS-NATIVE:                                  │
│     → AWS Bedrock como plataforma principal                 │
│     → API directa solo para features exclusivos             │
│                                                             │
│  🚀 STARTUP / BAJO PRESUPUESTO:                             │
│     → Google Gemini API (tier gratuito)                     │
│     → Migrar a Bedrock cuando escales                       │
│                                                             │
│  🔬 AI-FIRST / INVESTIGACIÓN:                               │
│     → APIs directas para últimas capacidades                │
│     → Bedrock para producción estable                       │
│                                                             │
│  ⚡ ALTO VOLUMEN (>500M tokens/mes):                        │
│     → Evaluar self-hosting Llama                            │
│     → Negociar precios enterprise directamente              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Fuentes y Referencias

### Precios y Documentación Oficial

| Fuente | URL | Fecha Consulta |
|--------|-----|----------------|
| AWS Bedrock Pricing | https://aws.amazon.com/bedrock/pricing/ | 13 enero 2026 |
| Google Gemini API Pricing | https://ai.google.dev/gemini-api/docs/pricing | 13 enero 2026 |
| OpenAI API Pricing | https://openai.com/api/pricing/ | 13 enero 2026 |
| Anthropic Claude Pricing | https://www.anthropic.com/pricing | 13 enero 2026 |
| Mistral AI Pricing | https://mistral.ai/technology/ | 13 enero 2026 |

### Documentación Técnica

| Recurso | URL |
|---------|-----|
| AWS Bedrock User Guide | https://docs.aws.amazon.com/bedrock/latest/userguide/ |
| Anthropic API Documentation | https://docs.anthropic.com/ |
| Google AI Developer Docs | https://ai.google.dev/docs |
| OpenAI Platform Docs | https://platform.openai.com/docs/ |

### Notas Importantes

1. **Precios sujetos a cambios:** Los proveedores actualizan precios frecuentemente. Verificar siempre en fuente oficial.
2. **Disponibilidad regional:** No todos los modelos están disponibles en todas las regiones de AWS.
3. **Términos de servicio:** Revisar políticas de uso de datos de cada proveedor.
4. **SLAs:** Los SLAs enterprise requieren acuerdos específicos con cada proveedor.

---

*Documento generado el 13 de enero de 2026. Los precios y capacidades están sujetos a cambios por parte de los proveedores.*
