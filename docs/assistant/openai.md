# OpenAI Assistant

Use OpenAI's GPT models as your SnipKit assistant.

## Configuration

```yaml [config.yaml]
version: 1.3.0
config:
  assistant:
    providers:
      - type: openai
        # If set to false, OpenAI will not be used as an AI assistant.
        enabled: true
        # OpenAI Model to be used.
        model: gpt-4.1
        # The name of the environment variable holding the OpenAI API key.
        apiKeyEnv: SNIPKIT_OPENAI_API_KEY
        # Optional: Custom API endpoint.
        # endpoint: https://api.openai.com/v1
```

## API Key

You need to provide the API key for the OpenAI API via the environment variable specified in `apiKeyEnv`.

```sh [Set OpenAI API Key]
export SNIPKIT_OPENAI_API_KEY="your-api-key-here"
```

Get your API key from the [OpenAI Platform](https://platform.openai.com/api-keys). See [OpenAI Models](https://platform.openai.com/docs/models) for all available models.
