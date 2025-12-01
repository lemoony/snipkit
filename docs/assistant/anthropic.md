# Anthropic Claude Assistant

Use Anthropic's Claude models as your SnipKit assistant.

## Configuration

```yaml title="config.yaml"
version: 1.3.0
config:
  assistant:
    providers:
      - type: anthropic
        # If set to false, Anthropic will not be used as an AI assistant.
        enabled: true
        # Claude Model to be used.
        model: claude-sonnet-4.5
        # The name of the environment variable holding the Anthropic API key.
        apiKeyEnv: SNIPKIT_ANTHROPIC_API_KEY
```

## API Key

You need to provide the API key for the Anthropic API via the environment variable specified in `apiKeyEnv`.

```sh title="Set Anthropic API Key"
export SNIPKIT_ANTHROPIC_API_KEY="your-api-key-here"
```

Get your API key from the [Anthropic Console](https://console.anthropic.com/settings/keys). See [Anthropic Models](https://docs.anthropic.com/en/docs/about-claude/models) for all available models.
