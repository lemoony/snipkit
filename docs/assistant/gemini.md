# Google Gemini Assistant

Use Google's Gemini models as your SnipKit assistant.

## Configuration

```yaml [config.yaml]
version: 1.3.0
config:
  assistant:
    providers:
      - type: gemini
        # If set to false, Gemini will not be used as an AI assistant.
        enabled: true
        # Gemini Model to be used.
        model: gemini-1.5-flash
        # The name of the environment variable holding the Gemini API key.
        apiKeyEnv: SNIPKIT_GEMINI_API_KEY
        # Optional: Custom API endpoint.
        # endpoint: https://generativelanguage.googleapis.com
```

## API Key

You need to provide the API key for the Gemini API via the environment variable specified in `apiKeyEnv`.

```sh [Set Gemini API Key]
export SNIPKIT_GEMINI_API_KEY="your-api-key-here"
```

Get your API key from [Google AI Studio](https://aistudio.google.com/app/apikey). See [Gemini Models](https://ai.google.dev/gemini-api/docs/models/gemini) for all available models.
