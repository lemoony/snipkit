# Configuration

Snipkit requires a configuration file to be present. The configuration file resides in the snipkit home directory.

## Home directory

The path to the home directory is assumed to be `{$XDG_CONFIG_HOME}/snipkit`. The value of `XDG_CONFIG_HOME` is specified by
the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html). If not
set explicitly, a sensible default value for your operating system is assumed:

- `macOS`: ~/Library/Application Support/snipkit
- `Linux`: ~/.config/snipkit

You can specifiy another directory to be used by snipkit be setting the environment variable `SNIPKIT_HOME`. E.g., you may
want to put the following into your `~/.zshr` file:

```bash 
export SNIPKIT_HOME=~/.snipkit
```

## Config file

Upon initialization, the configuration file looks similar to this:

```yaml 
version: 1.0.0
config:
  style:
    # The theme defines the terminal colors used by Snipkit.
    # Available themes:default,dracula.
    theme: default
  # Your preferred editor to open the config file when typing 'snipkit config edit'.
  editor: "" # Defaults to a reasonable value for your operation system when empty.
  # The command which should run if you don't provide any subcommand.
  defaultRootCommand: "" # If not set, the help text will be shown.
  provider:
    snippetsLab:
      # Set to true if you want to use SnippetsLab.
      enabled: fales
      # Path to your *.snippetslablibrary file.
      # SnipKit will try to detect this file automatically when generating the config.
      libraryPath: "/path/to/main.snippetslablibrary/"
      # If this list is not empty, only those snippets that match the listed tags will be provided to you.
      includeTags: []
    fsLibrary:
      # If set to false, the files specified via libraryPath will not be provided to you.
      enabled: false
      # Paths directories that hold snippets files. Each file must hold one snippet only.
      libraryPath:
        - /path/to/file/system/library
      # Only files with endings which match one of the listed suffixes will be considered.
      suffixRegex:
        - .sh
      # If set to true, the files will not be parsed in advance. This means, only the filename can be used as the snippet name.
      lazyOpen: false
      # If set to true, the title comment will not be shown in the preview window.
      hideTitleInPreview: false
```

## Editor

When typing `snipkit config edit` the configuration file will be opened in an editor of your choice.

The default editor is defined by the `$VISUAL` or `$EDITOR` environment variables. This behavior can be overwritten by
setting the `editor` field in the configuration file to a non-empty string, e.g.:

```yaml
version: 1.0.0
config:
  editor: "code"
```

If no value is provided at all, SnipKit will try to use `vim`.

## Theme

SnipKit supports multiple themes out of the box and also allows you to define your own themes:

```yaml
version: 1.0.0
config:
  theme: "dracula"
```

If the theme is not shipped with snipkit, it will try to lookup a custom theme. If the theme is named `<xx>`, the theme file
must be located at `<SNIPKIT_HOME>/<xxx>.yaml`.

For a list of supported default themes, have a look at the Themes page.

## Default Root Command

Most of the time, you want to call the same subcommand, e.g. `print` or `exec`. You can configure `snipkit` so that this
command gets executed by default by editing the config:

```yaml 
version: 1.0.0
config:
  defaultRootCommand: "exec"
```

This way, calling `snipkit` will yield the same result as `snipkit exec`. If you want to call the `print` command instead,
you can still call `snipkit print`.

## Provider

> TODO: Explain the different providers
