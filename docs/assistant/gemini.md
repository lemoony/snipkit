# Gemini Assistant

## Configuration

```yaml title="config.yaml"
version: 1.2.0
config:
  assistant:
    gemini:
      # If set to false, Gemini will not be used as an AI assistant.
      enabled: true
      # Gemini API endpoint.
      endpoint: https://generativelanguage.googleapis.com
      # Gemini Model to be used (e.g., openai/gpt-4o)
      model: gemini-1.5-flash
      # The name of the environment variable holding the Gemini API key.
      apiKeyEnv: SNIPKIT_GEMINI_API_KEY
```

!!! info
    For this configuration, you will need to provide the API key for the Gemini API via the environment variable `SNIPKIT_GEMINI_API_KEY`.