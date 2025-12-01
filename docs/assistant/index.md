# SnipKit Assistant

The SnipKit Assistant helps you create SnipKit snippets using AI.

!!! tip
    Alternatively, you may try using [SnipKit GPT](https://chatgpt.com/g/g-A2y9cPWE7-snipkit-assistant) to generate scripts compatible with SnipKit.

![Assistant Demo](../images/assistant/assistant.gif)

## Supported AI Providers

- [OpenAI](openai.md) (GPT-4, GPT-4o)
- [Anthropic](anthropic.md) (Claude)
- [Google Gemini](gemini.md)
- [Ollama](ollama.md) (local models)
- [OpenAI-Compatible](openai-compatible.md) (Together.ai, Groq, Azure OpenAI, etc.)

## Generate Scripts

```sh title="Generate a Script"
snipkit assistant generate
```

```sh title="Root-Level Command for Convenience"
snipkit ai
```

SnipKit opens the generated script in the [configured editor](../configuration/overview.md#editor), allowing you to review and modify it if necessary. The script will be executed once you close the editor.

## Retry or Tweak Prompts

If the script didn't work as expected or if you want to add more information, you can do so after the execution.

![Assistant Wizard](../images/assistant//assistant-zip.gif)

SnipKit remembers the previous prompt and script output. Everything is automatically included for the next prompt, so you don't need to provide anything unless you want to add new details.

!!! tip
    If the script didn't work due to errors visible in the output, you may not need to provide a new prompt. Just leave it empty and press enter.

## Save Generated Snippets

SnipKit supports saving generated snippets to your [File System Library][fslibrary].

![Assistant Wizard](../images/assistant/assistant-save.gif)

If you set `saveMode: FS_LIBRARY`, the assistant will ask whether you want to save the generated script after its execution.

```yaml title="config.yaml"
version: 1.3.0
config:
  assistant:
    saveMode: FS_LIBRARY
```

!!! note
    The [File System Library manager][fslibrary] must be enabled.

[fslibrary]: ../managers/fslibrary.md

## Configuration

This command lets you enable the assistant by editing your SnipKit configuration file:

```sh title="Enable or Switch to a Different AI Provider"
snipkit assistant choose
```

![Assistant Choose](../images/assistant/assistant-choose.gif)

You will need to provide an API key for the corresponding AI provider via an environment variable.

[configuration]: ../configuration/overview.md

```yaml title="config.yaml"
version: 1.3.0
config:
  assistant:
    # Defines if you want to save the snippets created by the assistant. Possible values: NEVER | FS_LIBRARY
    saveMode: NEVER
    providers:
      - type: openai
        enabled: true
        model: gpt-4.1
        apiKeyEnv: SNIPKIT_OPENAI_API_KEY
      - type: anthropic
        enabled: false
        model: claude-sonnet-4.5
        apiKeyEnv: SNIPKIT_ANTHROPIC_API_KEY
      - type: gemini
        enabled: false
        model: gemini-1.5-flash
        apiKeyEnv: SNIPKIT_GEMINI_API_KEY
      - type: ollama
        enabled: false
        model: llama3
        serverUrl: http://localhost:11434
```

!!! note
    The first provider with `enabled: true` will be used. If all providers are set to `enabled: false`, the assistant will not function.
