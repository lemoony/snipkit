# OpenAI-Compatible Assistant

Use any service that implements the OpenAI API format as your SnipKit assistant. This includes cloud providers like OpenRouter, Together.ai, and Groq, as well as local servers like llama.cpp and vLLM.

## Configuration

```yaml title="config.yaml"
version: 1.3.0
config:
  assistant:
    providers:
      - type: openai-compatible
        # If set to false, this provider will not be used.
        enabled: true
        # The model name (specific to the provider).
        model: anthropic/claude-sonnet-4
        # The API endpoint URL.
        endpoint: https://openrouter.ai/api/v1
        # The name of the environment variable holding the API key.
        apiKeyEnv: OPENROUTER_API_KEY
```

Get your API key from [OpenRouter](https://openrouter.ai/keys).

## Example Configurations

### OpenRouter

```yaml
- type: openai-compatible
  enabled: true
  model: anthropic/claude-sonnet-4
  endpoint: https://openrouter.ai/api/v1
  apiKeyEnv: OPENROUTER_API_KEY
```

### Azure OpenAI

```yaml
- type: openai-compatible
  enabled: true
  model: gpt-4.1
  endpoint: https://your-resource.openai.azure.com/openai/deployments/your-deployment
  apiKeyEnv: AZURE_OPENAI_API_KEY
```

### Local llama.cpp Server

```yaml
- type: openai-compatible
  enabled: true
  model: local-model
  endpoint: http://localhost:8080/v1
  # apiKeyEnv not needed for local servers
```

## Use Cases

This provider type is useful for:

- **Cloud AI platforms**: OpenRouter, Together.ai, Groq, Perplexity
- **Enterprise deployments**: Azure OpenAI, AWS Bedrock (with OpenAI-compatible proxies)
- **Local inference**: llama.cpp with `--api` flag, vLLM, text-generation-inference
- **Custom proxies**: Any OpenAI-compatible API proxy or gateway
