# OpenAI Assistant

## Configuration

```yaml title="config.yaml"
version: 1.2.0
config:
  assistant:
    openai:
      # If set to false, OpenAI will not be used as an AI assistant.
      enabled: true
      # OpenAI API endpoint.
      endpoint: https://api.openai.co
      # OpenAI Model to be used (e.g., openai/gpt-4o)
      model: openai/gpt-4o
      # The name of the environment variable holding the OpenAI API key.
      apiKeyEnv: SNIPKIT_OPENAI_API_KEY
```

!!! info
    For this configuration, you will need to provide the API key for the OpenAI API via the environment variable `SNIPKIT_OPENAI_API_KEY`.