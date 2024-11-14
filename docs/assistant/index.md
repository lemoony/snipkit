# SnipKit Assistant

The SnipKit Assistant helps you create SnipKit snippets using AI.

!!! tip
    Alternatively, you may try using [SnipKit GPT](https://chatgpt.com/g/g-A2y9cPWE7-snipkit-assistant) to generate scripts compatible with SnipKit.

![Assistant Demo](../images/assistant/assistant.gif)

!!! warning
    SnipKit Assistant is currently in beta for OpenAI and Gemini. Several improvements are in progress, so stay tuned!

## Supported AI Providers

- OpenAI
- Gemini

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
version: 1.2.0
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
version: 1.2.0
config:
  assistant:
    # Defines if you want to save the snippets created by the assistant. Possible values: NEVER | FS_LIBRARY
    saveMode: NEVER
    openai:
      enabled: false
      # ....
    gemini:
      enabled: false
      # ....
```

!!! warning
    Only one AI provider can be set to `enabled: true` at a time. If all providers are set to `enabled: false`, the assistant will not function.

