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
export SNIPKIT_HOME=/Users/<user>/.snipkit
```

!!! warning "Use absolute paths for `SNIPKIT_HOME`"
    Please make sure to use an absolute path for `SNIPKIT_HOME` and do not use the tilde (~) character to point to your 
    home directory. Otherwise, SnipKit will not be able to find your configuration file.

## Initialization

In order to create a config file for SnipKit, execute:

```bash
snipkit config init
```

This command creates a config file in the SnipKit home directory. The initial config file looks similar to this:

```yaml title="config.yaml"
version: 1.1.1
config:
  style:
    # The theme defines the terminal colors used by Snipkit.
    # Available themes:default,dracula.
    theme: default
  # Your preferred editor to open the config file when typing 'snipkit config edit'.
  editor: "" # Defaults to a reasonable value for your operation system when empty.
  # The command which should run if you don't provide any subcommand.
  defaultRootCommand: "" # If not set, the help text will be shown.
  # Enable fuzzy searching for snippet titles.
  fuzzySearch: true
```

No snippet manager has been added at this time. In order to add a one execute:

```bash
snipkit manager add
```

For more information on the different managers supported, please see [Managers][managers].

## Migration

The config file may change from time to time. If you use an outdated config file version, Snipkit will refuse to run. 
Fortunately, you can migrate your config file to the latest version:

```bash
snipkit config migrate
```

## Config options

### Editor

When typing `snipkit config edit` the configuration file will be opened in an editor of your choice.

The default editor is defined by the `$VISUAL` or `$EDITOR` environment variables. This behavior can be overwritten by
setting the `editor` field in the configuration file to a non-empty string, e.g.:

```yaml title="config.yaml"
version: 1.1.1
config:
  editor: "code"
```

If no value is provided at all, SnipKit will try to use `vim`.

### Default Root Command

Most of the time, you want to call the same subcommand, e.g. `print` or `exec`. You can configure `snipkit` so that this
command gets executed by default:

```yaml title="config.yaml"
version: 1.1.1
config:
  defaultRootCommand: "exec"
```

This way, calling `snipkit` will yield the same result as `snipkit exec`. If you want to call the `print` command instead,
you can still call `snipkit print`.

### Fuzzy search

Enable fuzzy searching for snippet titles. This leads to potentially more snippets matching the search criteria. Snipkit
will try to rank them according to similarity. Disable fuzzy search for performance reason or if you just don't like.

```yaml title="config.yaml"
version: 1.1.1
config:
  fuzzySearch: true
```

### Style

#### Theme 

SnipKit supports multiple themes out of the box and also allows you to define your own themes:

```yaml title="config.yaml"
version: 1.1.1
config:
  style:
    theme: "default"
```

If the theme is not shipped with snipkit, it will try to look up a custom theme. If the theme is named `<xx>`, the theme file
must be located at `<SNIPKIT_HOME>/<xxx>.yaml`.

For a list of supported default themes, have a look at the [Themes][themes] page.

#### Hide Keymap

By default, a help for the key mapping is displayed at the bottom of the screen. To save same screen space, this can be 
disabled:

```yaml title="config.yaml"
version: 1.1.1
config:
  style:
    hideKeyMap: true
```

### Script Options

#### Shell

The shell for script executions is defined by the `$SHELL` environment variable. This behavior can be overwritten by setting
the `shell` option to a non-empty string, e.g.:

```yaml title="config.yaml"
version: 1.1.1 
config:
  script:
    shell: "/bin/zsh"
```

If neither `$SHELL` nor the config option `shell` is defined, SnipKit will try to use `/bin/bash` as a fallback value.

#### Parameter mode

How values are injected into your snippet for the defined parameters is defined by the `parameterMode` option:

```yaml title="config.yaml"
version: 1.1.1
config:
  script:
    parameterMode: SET
```

The default value is `SET`, defining that values should be set as variables. This means that the following script

```sh title="Raw snippet before execution"
# ${VAR} Description: What to print
echo ${VAR}
```

will be updated in the following way, e.g. for `VAR = "Hello word"`:

```sh  title="Example for parameterMode SET"
# ${VAR} Description: What to print
VAR="Hello world"
echo ${VAR}
```

Alternatively, all occurrences of a parameter can be replaced with the actual value when 
specifying `REPLACE` for `parameterMode`:

```sh title="Example for parameterMode = REPLACE"
echo "Hello world"
```

#### Remove Comments

SnipKit will remove all parameter comments from a snippet when specifying `removeComments`:

```yaml title="config.yaml"
version: 1.1.1
config:
  script:
    removeComments: true
```

This means that the following script

```sh title="Raw snippet before execution"
# ${VAR} Description: What to print
echo ${VAR}
```

will be formatted in the following way:

```sh  title="Example for removeComments = true"
echo ${VAR}
```

!!! info
    Comments will always be removed if `parameterMode` is set to `REPLACE`.

#### Confirm Commands

If you always want to explicitly confirm the command before execution, specify the `execConfirm` option as follows:

```yaml title="config.yaml"
version: 1.1.1
config:
  script:
    execConfirm: true
```

!!! tip "Flag --confirm"
    The same functionality can be achieved by means of the `--confirm` flag:
    ```bash
    snipkit exec --confirm
    ```
    Use the flag instead of the config option if you only want to explicitly confirm the command in some cases.


#### Print Commands

SnipKit will print the command to be executed on stdout if specified by the `execPrint` commands:

```yaml title="config.yaml"
version: 1.1.1
config:
  script:
    execPrint: true
```

!!! tip "Flag -p or --print"
    The same functionality can be achieved by means of the `-p` or `--print` flag:
    ```bash 
    snipkit exec --print
    ```
    Use the flag instead of the config option if you only want to print the command every now and then.

## Clean up

The config file as well as all custom themes can be deleted with:

```bash 
snipkit config clean
```

The cleanup method is a way to remove all SnipKit artifacts from your hard drive. It only deletes contents of the SnipKit
home directory. If this directory is empty at the end of the cleanup process, it will be deleted as well.

[managers]: ../managers/overview.md

[themes]: themes.md
