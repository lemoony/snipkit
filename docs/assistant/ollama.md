# Ollama Assistant

Use local LLMs via Ollama as your SnipKit assistant. This allows you to run AI models entirely on your own machine without sending data to external APIs.

## Prerequisites

1. Install [Ollama](https://ollama.ai/)
2. Pull a model: `ollama pull llama3`
3. Ensure Ollama is running: `ollama serve`

## Configuration

```yaml title="config.yaml"
version: 1.3.0
config:
  assistant:
    providers:
      - type: ollama
        # If set to false, Ollama will not be used as an AI assistant.
        enabled: true
        # The model to use (must be pulled first).
        model: llama3
        # Ollama server URL (defaults to localhost).
        serverUrl: http://localhost:11434
```

!!! note
    No API key is required since Ollama runs locally.

Pull a model before using it:

```sh
ollama pull llama3
```

See [Ollama Library](https://ollama.ai/library) for all available models.

## Custom Server

If Ollama is running on a different machine or port, update the `serverUrl`:

```yaml
serverUrl: http://192.168.1.100:11434
```
