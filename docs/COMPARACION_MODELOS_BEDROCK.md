# Comparacion de Modelos de IA en AWS Bedrock

**Documento actualizado:** 13 de enero de 2026  
**Tipo de cambio utilizado:** 1 USD = 17.85 MXN

---

## Tabla de Contenidos

1. [Introduccion](#introduccion)
2. [Resumen Ejecutivo](#resumen-ejecutivo)
3. [Modelos Amazon Nova](#modelos-amazon-nova)
4. [Modelos Anthropic Claude](#modelos-anthropic-claude)
5. [Modelos OpenAI](#modelos-openai)
6. [Modelos Google Gemini](#modelos-google-gemini)
7. [Modelos Meta Llama](#modelos-meta-llama)
8. [Modelos Mistral AI](#modelos-mistral-ai)
9. [Comparativa de Precios](#comparativa-de-precios)
10. [Limites de Throughput y Rate Limits](#limites-de-throughput-y-rate-limits)
11. [Soporte de Prompt Caching](#soporte-de-prompt-caching)
12. [Casos de Uso Empresarial](#casos-de-uso-empresarial)
13. [Recomendaciones por Escenario](#recomendaciones-por-escenario)

---

## Introduccion

Este documento presenta una comparativa detallada de los principales modelos de inteligencia artificial disponibles a traves de AWS Bedrock, incluyendo sus capacidades, costos, limites de uso y casos de aplicacion en entornos empresariales.

AWS Bedrock ofrece acceso a modelos de multiples proveedores a traves de una API unificada, permitiendo a las organizaciones seleccionar el modelo mas adecuado para cada caso de uso sin necesidad de gestionar infraestructura.

---

## Resumen Ejecutivo

| Categoria | Modelo Recomendado | Justificacion |
|-----------|-------------------|---------------|
| Mejor relacion costo-rendimiento | Amazon Nova 2 Lite | Razonamiento rapido, bajo costo, ideal para tareas cotidianas |
| Mayor capacidad de razonamiento | Claude Opus 4.5 | Inteligencia frontera, excelente para agentes y codigo |
| Velocidad maxima | Claude Haiku 4.5 / Nova 2 Lite | Latencia minima para aplicaciones en tiempo real |
| Vision y multimodal | Claude Sonnet 4.5 / Nova 2 Omni | Comprension avanzada de imagenes, video y documentos |
| Generacion de codigo | Claude Sonnet 4.5 / Claude Opus 4 | Rendimiento superior en tareas de programacion |
| Procesamiento masivo | Cualquier modelo con Batch API | 50% de descuento en procesamiento por lotes |

---

## Modelos Amazon Nova

Amazon Nova es la familia de modelos propios de AWS, optimizados para ofrecer el mejor balance entre rendimiento y costo.

### Nova 2 Lite

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de razonamiento rapido y economico |
| **Contexto maximo** | 1,000,000 tokens |
| **Modalidades** | Texto, imagenes, video |
| **Casos de uso** | Chatbots, procesamiento de documentos, automatizacion empresarial |
| **Herramientas integradas** | Code interpreter, web grounding, soporte MCP remoto |

**Capacidades principales:**
- Control de esfuerzo de razonamiento para desarrolladores
- Herramientas integradas como interprete de codigo y anclaje web
- Soporte para herramientas MCP remotas
- Ventana de contexto de 1M de tokens
- Ideal para cargas de trabajo de IA cotidianas

### Nova 2 Pro (Preview)

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo multiproposito de alto rendimiento |
| **Contexto maximo** | 1,000,000 tokens |
| **Modalidades** | Texto, imagenes, video, audio |
| **Casos de uso** | Tareas empresariales complejas, analisis multimodal |

### Nova 2 Omni (Preview)

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo multimodal unificado |
| **Modalidades** | Texto, imagenes, video, voz |
| **Casos de uso** | Experiencias conversacionales multimodales, analisis creativo |

### Nova 2 Sonic

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de voz speech-to-speech |
| **Latencia** | Baja latencia para conversaciones naturales |
| **Casos de uso** | Asistentes de voz, recepcionistas virtuales, IVR inteligente |

### Precios Nova (por 1,000 tokens) - Region US East

| Modelo | Entrada (USD) | Entrada (MXN) | Salida (USD) | Salida (MXN) |
|--------|---------------|---------------|--------------|--------------|
| Nova 2 Lite | $0.00006 | $0.0011 | $0.00024 | $0.0043 |
| Nova 2 Pro | $0.0008 | $0.0143 | $0.0032 | $0.0571 |
| Nova Canvas (imagen) | - | - | $0.04/imagen | $0.71/imagen |
| Nova Reel (video) | - | - | $0.08/seg | $1.43/seg |

---

## Modelos Anthropic Claude

Anthropic Claude es una familia de modelos de lenguaje grande (LLM) conocida por su enfoque en seguridad y razonamiento avanzado.

### Claude Opus 4.5

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de inteligencia frontera, el mas potente de Claude |
| **Contexto maximo** | 200,000 tokens |
| **Modalidades** | Texto, imagenes, documentos |
| **Capacidades especiales** | Razonamiento hibrido (respuestas instantaneas + pensamiento extendido) |
| **Casos de uso** | Codigo de produccion, agentes sofisticados, tareas de oficina complejas |

### Claude Sonnet 4.5

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de alto rendimiento balanceado |
| **Contexto maximo** | 200,000+ tokens |
| **Modalidades** | Texto, imagenes, video, documentos |
| **Capacidades especiales** | Mejor modelo para coding y agentes complejos |
| **Casos de uso** | Razonamiento sostenido, flujos de trabajo multi-paso, computer use |

### Claude Haiku 4.5

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo rapido y economico |
| **Velocidad** | El mas rapido de la familia Claude |
| **Casos de uso** | Alto volumen, baja latencia, tareas bien definidas |

### Claude 3.7 Sonnet

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Primer modelo hibrido de razonamiento |
| **Capacidades** | Extended thinking para problemas complejos |
| **Casos de uso** | Coding agentivo, resolucion creativa de problemas |

### Precios Claude (por 1,000 tokens)

| Modelo | Entrada (USD) | Entrada (MXN) | Salida (USD) | Salida (MXN) |
|--------|---------------|---------------|--------------|--------------|
| Claude Opus 4.5 | $0.005 | $0.089 | $0.025 | $0.446 |
| Claude Sonnet 4.5 (<=200K) | $0.003 | $0.054 | $0.015 | $0.268 |
| Claude Sonnet 4.5 (>200K) | $0.006 | $0.107 | $0.0225 | $0.402 |
| Claude Haiku 4.5 | $0.001 | $0.018 | $0.005 | $0.089 |
| Claude 3.5 Sonnet | $0.003 | $0.054 | $0.015 | $0.268 |

---

## Modelos OpenAI

OpenAI ofrece modelos a traves de AWS Bedrock con sus modelos de codigo abierto (open weight).

### GPT-OSS-120B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de razonamiento de codigo abierto de alto rendimiento |
| **Parametros** | 120 mil millones |
| **Casos de uso** | Coding avanzado, matematicas complejas, tareas agentivas, investigacion profunda |

### GPT-OSS-20B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo compacto de razonamiento |
| **Parametros** | 20 mil millones |
| **Casos de uso** | Matematicas, coding, tareas analiticas con despliegue eficiente |

### GPT-OSS-Safeguard-120B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de moderacion de contenido de gran escala |
| **Casos de uso** | Politicas personalizadas, moderacion con razonamiento detallado |

### GPT-OSS-Safeguard-20B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de seguridad de contenido |
| **Casos de uso** | Clasificacion de contenido danino, workflows de trust & safety |

### Precios OpenAI API Directa (Referencia)

| Modelo | Entrada (USD/1M) | Entrada (MXN/1M) | Salida (USD/1M) | Salida (MXN/1M) |
|--------|------------------|------------------|-----------------|-----------------|
| GPT-5.2 | $1.75 | $31.24 | $14.00 | $249.90 |
| GPT-5.2 Pro | $21.00 | $374.85 | $168.00 | $2,998.80 |
| GPT-5 Mini | $0.25 | $4.46 | $2.00 | $35.70 |
| GPT-4.1 | $3.00 | $53.55 | $12.00 | $214.20 |
| GPT-4.1 Mini | $0.80 | $14.28 | $3.20 | $57.12 |

---

## Modelos Google Gemini

Google ofrece su familia Gemini con capacidades multimodales avanzadas.

### Gemini 3 Pro

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo multimodal mas potente de Google |
| **Contexto** | 200,000+ tokens |
| **Modalidades** | Texto, imagenes, video, audio |
| **Casos de uso** | Comprension multimodal, agentes, vibe-coding |

### Gemini 3 Flash

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo inteligente optimizado para velocidad |
| **Capacidades** | Busqueda superior y anclaje |
| **Casos de uso** | Aplicaciones de alta velocidad con grounding en busqueda |

### Gemini 2.5 Pro

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo multiproposito de ultima generacion |
| **Contexto** | 1,000,000 tokens |
| **Casos de uso** | Coding, razonamiento complejo, tareas STEM |

### Gemini 2.5 Flash

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo hibrido de razonamiento |
| **Contexto** | 1,000,000 tokens |
| **Casos de uso** | Procesamiento a escala, baja latencia, uso agentivo |

### Precios Gemini API (por 1,000,000 tokens)

| Modelo | Entrada (USD) | Entrada (MXN) | Salida (USD) | Salida (MXN) |
|--------|---------------|---------------|--------------|--------------|
| Gemini 3 Pro (<=200K) | $2.00 | $35.70 | $12.00 | $214.20 |
| Gemini 3 Pro (>200K) | $4.00 | $71.40 | $18.00 | $321.30 |
| Gemini 3 Flash | $0.50 | $8.93 | $3.00 | $53.55 |
| Gemini 2.5 Pro (<=200K) | $1.25 | $22.31 | $10.00 | $178.50 |
| Gemini 2.5 Flash | $0.30 | $5.36 | $2.50 | $44.63 |
| Gemini 2.5 Flash-Lite | $0.10 | $1.79 | $0.40 | $7.14 |
| Gemini 2.0 Flash | $0.10 | $1.79 | $0.40 | $7.14 |

---

## Modelos Meta Llama

Meta ofrece modelos Llama de codigo abierto a traves de Bedrock.

### Llama 3.1 (8B, 70B, 405B)

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelos de proposito general |
| **Variantes** | 8B, 70B, 405B parametros |
| **Contexto** | Hasta 128,000 tokens |
| **Casos de uso** | Generacion de texto, clasificacion, extraccion de informacion |

### Precios Meta Llama (por 1,000 tokens)

| Modelo | Entrada (USD) | Entrada (MXN) | Salida (USD) | Salida (MXN) |
|--------|---------------|---------------|--------------|--------------|
| Llama 3.1 8B | $0.00022 | $0.0039 | $0.00022 | $0.0039 |
| Llama 3.1 70B | $0.00099 | $0.0177 | $0.00099 | $0.0177 |
| Llama 3.1 405B | $0.00532 | $0.0950 | $0.016 | $0.2856 |
| Llama 2 Chat 13B | $0.00075 | $0.0134 | $0.001 | $0.0179 |

---

## Modelos Mistral AI

Mistral AI ofrece modelos eficientes con excelente rendimiento.

### Mistral Large

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo de alto rendimiento |
| **Casos de uso** | Razonamiento complejo, generacion de texto avanzada |

### Mixtral 8x7B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo Mixture of Experts |
| **Casos de uso** | Balance entre rendimiento y eficiencia |

### Mistral 7B

| Caracteristica | Especificacion |
|----------------|----------------|
| **Tipo** | Modelo compacto y eficiente |
| **Casos de uso** | Tareas de bajo costo y alta velocidad |

### Precios Mistral AI (por 1,000 tokens)

| Modelo | Entrada (USD) | Entrada (MXN) | Salida (USD) | Salida (MXN) |
|--------|---------------|---------------|--------------|--------------|
| Mistral Large | $0.008 | $0.143 | $0.024 | $0.428 |
| Mixtral 8x7B | $0.00045 | $0.008 | $0.0007 | $0.012 |
| Mistral 7B | $0.00015 | $0.003 | $0.0002 | $0.004 |

---

## Comparativa de Precios

### Tabla Comparativa General (por 1,000 tokens en MXN)

| Modelo | Entrada | Salida | Costo Total* |
|--------|---------|--------|--------------|
| **Economicos** | | | |
| Nova 2 Lite | $0.0011 | $0.0043 | $0.0054 |
| Mistral 7B | $0.003 | $0.004 | $0.007 |
| Llama 3.1 8B | $0.004 | $0.004 | $0.008 |
| Mixtral 8x7B | $0.008 | $0.012 | $0.020 |
| **Gama Media** | | | |
| Claude Haiku 4.5 | $0.018 | $0.089 | $0.107 |
| Nova 2 Pro | $0.014 | $0.057 | $0.071 |
| Llama 3.1 70B | $0.018 | $0.018 | $0.036 |
| **Premium** | | | |
| Claude Sonnet 4.5 | $0.054 | $0.268 | $0.322 |
| Mistral Large | $0.143 | $0.428 | $0.571 |
| Claude Opus 4.5 | $0.089 | $0.446 | $0.535 |
| Llama 3.1 405B | $0.095 | $0.286 | $0.381 |

*Costo total asume 1,000 tokens de entrada y 1,000 de salida

### Costo Estimado por Caso de Uso (MXN)

| Caso de Uso | Tokens Aprox. | Nova 2 Lite | Claude Haiku | Claude Sonnet |
|-------------|---------------|-------------|--------------|---------------|
| Pregunta simple | 500 in / 200 out | $0.0014 | $0.027 | $0.081 |
| Analisis de documento | 5,000 in / 1,000 out | $0.010 | $0.179 | $0.536 |
| Generacion de reporte | 2,000 in / 3,000 out | $0.015 | $0.304 | $0.912 |
| Analisis de imagen | 1,500 in / 500 out | $0.004 | $0.072 | $0.215 |

---

## Limites de Throughput y Rate Limits

Los limites de throughput en AWS Bedrock varian segun el modelo, la region y el tier de servicio.

### Tiers de Servicio en Bedrock

| Tier | Caracteristicas | Precio |
|------|-----------------|--------|
| **Standard** | Rendimiento consistente a tarifas regulares | Precio base |
| **Priority** | Asignacion preferencial de computo para aplicaciones criticas | Premium |
| **Flex** | Descuento a cambio de procesamiento diferido | -50% aprox. |
| **Batch** | Procesamiento asincrono en 24 horas | -50% |

### Limites por Modelo (Valores Tipicos por Region)

| Modelo | Peticiones/min (Default) | Tokens/min (Default) | Solicitable |
|--------|--------------------------|----------------------|-------------|
| Nova 2 Lite | 100-500 | 100,000-500,000 | Si |
| Nova 2 Pro | 50-200 | 50,000-200,000 | Si |
| Claude Sonnet | 50-100 | 100,000-200,000 | Si |
| Claude Haiku | 100-500 | 200,000-500,000 | Si |
| Claude Opus | 20-50 | 50,000-100,000 | Si |
| Llama 3.1 70B | 50-100 | 100,000-200,000 | Si |

**Nota:** Los limites exactos dependen de la region, la cuenta y los acuerdos con AWS. Se pueden solicitar aumentos a traves de Service Quotas.

### Cross-Region Inference

Bedrock ofrece inferencia cross-region para mayor disponibilidad:

- **Global Cross-Region:** Enruta automaticamente a la region con menor carga
- **Geo Cross-Region:** Mantiene los datos dentro de una region geografica (US, EU, APAC)

---

## Soporte de Prompt Caching

El prompt caching permite reutilizar tokens de entrada previamente procesados, reduciendo costos y latencia.

### Modelos con Soporte de Prompt Caching

| Proveedor | Modelo | Soporte | Descuento Read | Costo Write |
|-----------|--------|---------|----------------|-------------|
| Anthropic | Claude Opus 4.5 | Si | 90% | +25% |
| Anthropic | Claude Sonnet 4.5 | Si | 90% | +25% |
| Anthropic | Claude Haiku 4.5 | Si | 90% | +25% |
| Amazon | Nova 2 Lite | Si* | Variable | Variable |
| Amazon | Nova 2 Pro | Si* | Variable | Variable |
| Google | Gemini 2.5 Pro | Si | 90% | Incluido |
| Google | Gemini 2.5 Flash | Si | 90% | Incluido |
| OpenAI | GPT-5.2 | Si | 90% | Incluido |

*Disponibilidad sujeta a actualizaciones del servicio

### Precios de Prompt Caching (Claude via API directa)

| Modelo | Lectura Cache (USD/MTok) | Escritura Cache (USD/MTok) |
|--------|--------------------------|----------------------------|
| Claude Opus 4.5 | $0.50 | $6.25 |
| Claude Sonnet 4.5 (<=200K) | $0.30 | $3.75 |
| Claude Haiku 4.5 | $0.10 | $1.25 |

### Casos de Uso para Prompt Caching

1. **System prompts largos:** Reutilizar instrucciones base en multiples conversaciones
2. **Documentos de contexto:** Mantener documentos de referencia en cache
3. **Few-shot examples:** Cachear ejemplos para clasificacion o extraccion
4. **Tool definitions:** Cachear definiciones de herramientas para agentes

---

## Casos de Uso Empresarial

### 1. Servicio al Cliente

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Chatbot de alto volumen | Claude Haiku 4.5 | Bajo costo, alta velocidad |
| Resolucion de tickets complejos | Claude Sonnet 4.5 | Razonamiento avanzado, tool use |
| Asistente de voz | Nova 2 Sonic | Speech-to-speech nativo |
| Triaje automatico | Nova 2 Lite | Costo-efectivo para clasificacion |

### 2. Procesamiento de Documentos

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Extraccion de datos de facturas | Nova 2 Lite | Vision multimodal, bajo costo |
| Analisis de contratos legales | Claude Opus 4.5 | Razonamiento profundo, precision |
| OCR inteligente | Nova 2 Pro | Comprension contextual de documentos |
| Resumen de reportes | Claude Sonnet 4.5 | Balance calidad-costo |

### 3. Desarrollo de Software

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Generacion de codigo | Claude Sonnet 4.5 | Mejor rendimiento en coding |
| Code review automatizado | Claude Opus 4 | Analisis detallado, sugerencias |
| Documentacion tecnica | Nova 2 Pro | Generacion eficiente |
| Debugging asistido | Claude Opus 4.5 | Razonamiento extendido |

### 4. Analisis de Datos

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Interpretacion de graficas | Nova 2 Omni | Comprension visual avanzada |
| Generacion de insights | Claude Sonnet 4.5 | Razonamiento analitico |
| SQL generation | Nova 2 Lite | Rapido y economico |
| ETL inteligente | Claude Haiku 4.5 | Alto volumen, bajo costo |

### 5. Contenido y Marketing

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Redaccion de contenido | Claude Sonnet 4.5 | Calidad de escritura |
| Generacion de imagenes | Nova Canvas | Generacion nativa en Bedrock |
| Traduccion empresarial | Nova 2 Pro | Multilenguaje, contexto |
| Social media automation | Claude Haiku 4.5 | Volumen alto, costo bajo |

### 6. Seguridad y Compliance

| Requerimiento | Modelo Recomendado | Justificacion |
|---------------|-------------------|---------------|
| Moderacion de contenido | GPT-OSS-Safeguard | Especializado en seguridad |
| Deteccion de fraude | Claude Opus 4.5 | Analisis profundo de patrones |
| Analisis de logs | Nova 2 Lite | Procesamiento masivo |
| PII detection | Claude Haiku 4.5 | Rapido, preciso |

---

## Recomendaciones por Escenario

### Startup / MVP

**Objetivo:** Minimizar costos, validar concepto

| Componente | Recomendacion |
|------------|---------------|
| Modelo principal | Nova 2 Lite |
| Procesamiento batch | Batch API (-50%) |
| Fallback | Claude Haiku 4.5 |
| Costo mensual estimado (10K requests) | ~$50-100 MXN |

### Empresa Mediana

**Objetivo:** Balance costo-calidad, escalabilidad

| Componente | Recomendacion |
|------------|---------------|
| Tareas simples | Nova 2 Lite / Claude Haiku |
| Tareas complejas | Claude Sonnet 4.5 |
| Procesamiento masivo | Batch API + Prompt Caching |
| Costo mensual estimado (100K requests) | ~$5,000-15,000 MXN |

### Enterprise

**Objetivo:** Maxima calidad, SLAs, compliance

| Componente | Recomendacion |
|------------|---------------|
| Tareas criticas | Claude Opus 4.5 (Priority Tier) |
| Alto volumen | Claude Sonnet + Prompt Caching |
| Agentes autonomos | Nova 2 Pro + Claude Opus |
| Cross-region | Global Cross-Region Inference |
| Costo mensual estimado (1M requests) | $50,000-200,000+ MXN |

---

## Consideraciones Adicionales

### Seguridad y Compliance

- Todos los modelos en Bedrock cumplen con SOC 1/2/3, ISO 27001, HIPAA eligible
- Los datos no se utilizan para entrenar modelos de terceros
- Opcion de VPC endpoints para trafico privado
- Cifrado en transito y en reposo

### Optimizacion de Costos

1. **Usar el modelo adecuado:** No usar Claude Opus para tareas simples
2. **Implementar Prompt Caching:** Hasta 90% de ahorro en tokens de entrada
3. **Batch Processing:** 50% de descuento para procesamiento no urgente
4. **Monitorear uso:** CloudWatch metrics para identificar ineficiencias
5. **Cross-Region Inference:** Mejor disponibilidad sin costo adicional

### Roadmap y Tendencias

- **Nova 2:** Mejoras continuas en razonamiento y herramientas
- **Claude 4.x:** Enfoque en agentes autonomos y coding
- **Gemini 3:** Integracion con busqueda y herramientas de Google
- **Precios:** Tendencia a la baja con mayor competencia

---

## Referencias

- [AWS Bedrock Pricing](https://aws.amazon.com/bedrock/pricing/)
- [Amazon Nova Models](https://aws.amazon.com/nova/models/)
- [Anthropic Claude Documentation](https://docs.anthropic.com/)
- [Google Gemini API Pricing](https://ai.google.dev/gemini-api/docs/pricing)
- [OpenAI API Pricing](https://openai.com/api/pricing/)

---

*Documento generado el 13 de enero de 2026. Los precios y capacidades estan sujetos a cambios por parte de los proveedores.*
