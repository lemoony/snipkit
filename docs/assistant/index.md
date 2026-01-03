# SnipKit Assistant

The SnipKit Assistant helps you create SnipKit snippets using AI through an interactive chat-style interface. Generate scripts from natural language prompts, execute them with real-time output preview, and refine with follow-up prompts—all in a unified workflow.

::: tip
Alternatively, you may try using [SnipKit GPT](https://chatgpt.com/g/g-A2y9cPWE7-snipkit-assistant) to generate scripts compatible with SnipKit.
:::

![Assistant Demo](/images/assistant/assistant.gif)

## Supported AI Providers

- [OpenAI](openai)
- [Anthropic](anthropic)
- [Google Gemini](gemini)
- [Ollama](ollama) (local models)
- [OpenAI-Compatible](openai-compatible) (Together.ai, Groq, Azure OpenAI, etc.)

## Generate Scripts

```sh [Generate a Script]
snipkit assistant generate
```

```sh [Root-Level Command for Convenience]
snipkit ai
```

After entering your prompt, SnipKit displays the generated script directly in the chat interface. You can then choose from the action bar:

- **Execute** (`E`) - Run the script immediately (prompts for parameter values if needed)
- **Open Editor** (`O`) - Edit the script in your [configured editor](../configuration/overview.md#editor) before execution
- **Revise** (`R`) - Provide a follow-up prompt to refine the script
- **Cancel** (`C`) - Exit without executing

## Execution Preview

When you execute a script, SnipKit displays the output in real-time as the command runs. After execution completes, you'll see:

- The complete script output (stdout and stderr)
- Exit code
- Execution duration

This information helps you quickly understand whether the script worked as expected and diagnose any issues.

## Revise Prompts

If the script didn't work as expected or you want to refine it, select **Revise** (`R`) from the action bar after execution. The chat interface preserves your full conversation history, including:

- Previous prompts
- Generated scripts
- Execution output and results

![Assistant Wizard](/images/assistant//assistant-zip.gif)

When you provide a follow-up prompt, SnipKit automatically includes the context from previous interactions, so you can simply describe what needs to change.

::: tip
If the script failed due to errors visible in the output, try revising with an empty prompt—the AI will use the error output to fix the issue automatically.
:::

## Save Generated Snippets

SnipKit supports saving generated snippets to your [File System Library][fslibrary].

![Assistant Wizard](/images/assistant/assistant-save.gif)

After executing a script, select **Save & Exit** (`S`) from the action bar. A modal dialog lets you specify:

- **Filename** - The file name for the saved script
- **Snippet Name** - A descriptive title for the snippet

If you set `saveMode: FS_LIBRARY`, the save option will be available in the post-execution action bar.

```yaml [config.yaml]
version: 1.3.0
config:
  assistant:
    saveMode: FS_LIBRARY
```

::: tip
The [File System Library manager][fslibrary] must be enabled.
:::

[fslibrary]: ../managers/fslibrary

## Configuration

This command lets you enable the assistant by editing your SnipKit configuration file:

```sh [Enable or Switch to a Different AI Provider]
snipkit assistant choose
```

![Assistant Choose](/images/assistant/assistant-choose.gif)

You will need to provide an API key for the corresponding AI provider via an environment variable.

[configuration]: ../configuration/overview

```yaml [config.yaml]
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

::: tip
The first provider with `enabled: true` will be used. If all providers are set to `enabled: false`, the assistant will not function.
:::
