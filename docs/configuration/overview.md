# Configuration

Snipkit requires a configuration file to be present. The configuration file resides in the snipkit home directory.

## Home directory

The path to the home directory is assumed to be `{$XDG_CONFIG_HOME}/snipkit`. The value of `XDG_CONFIG_HOME` is specified by
the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html). If not
set explicitly, a sensible default value for your operating system is assumed:

- `macOS`: ~/Library/Application Support/snipkit
- `Linux`: ~/.config/snipkit

You can specify another directory to be used by SnipKit be setting the environment variable `SNIPKIT_HOME`. E.g., you may
want to put the following into your `~/.zshrc` file:

```bash 
export SNIPKIT_HOME=~/.snipkit
```

## Initialization

In order to create a config file for SnipKit, execute:

```bash
snipkit config init
```

This command creates a config file in the SnipKit home directory. The initial config file looks similar to this:

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
```

No snippet manager has been added at this time. In order to add a one execute:

```bash
snipkit manager add
```

For more information on the different managers supported, please see [Managers][managers].

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
  theme: "default"
```

If the theme is not shipped with snipkit, it will try to look up a custom theme. If the theme is named `<xx>`, the theme file
must be located at `<SNIPKIT_HOME>/<xxx>.yaml`.

For a list of supported default themes, have a look at the [Themes][themes] page.

## Default Root Command

Most of the time, you want to call the same subcommand, e.g. `print` or `exec`. You can configure `snipkit` so that this
command gets executed by default:

```yaml 
version: 1.0.0
config:
  defaultRootCommand: "exec"
```

This way, calling `snipkit` will yield the same result as `snipkit exec`. If you want to call the `print` command instead,
you can still call `snipkit print`.


## Clean up

The config file as well as all custom themes can be deleted with:

```bash 
snipkit config clean
```

The cleanup method is a way to remove all SnipKit artifacts from your hard drive. It only deletes contents of the SnipKit
home directory. If this directory is empty at the end of the cleanup process, it will be deleted as well.

[managers]: ../managers/overview.md
[themes]: themes.md
