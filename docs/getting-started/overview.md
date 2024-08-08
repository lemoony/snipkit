# Overview

## Usage

```bash
$ snipkit -h
Snipkit helps you to execute scripts saved in your favorite snippets manager without even leaving the terminal.

Usage:
  snipkit [flags]
  snipkit [command]

Available Commands:
  browse      Browse all snippets without executing them
  completion  Generate the autocompletion script for the specified shell
  config      Manage your snipkit configuration file
  copy        Copies the snippet to the clipboard
  exec        Execute a snippet directly from the terminal
  export      Exports snippets on stdout
  help        Help about any command
  info        Provides useful information about the snipkit configuration
  manager     Manage the snippet managers snipkit connects to
  print       Prints the snippet on stdout
  sync        Synchronizes all snippet managers


Flags:
  -c, --config string      config file (default "/Users/pse/Library/Application Support/snipkit/config.yaml")
  -h, --help               help for snipkit
  -l, --log-level string   log level used for debugging problems (supported values: trace,debug,info,warn,error,fatal,panic) (default "panic")
  -v, --version            version for snipkit

Use "snipkit [command] --help" for more information about a command.
```

## Init a config file

Upon first usage, you have to create a config file.

```sh title="Create a new config file"
snipkit config init
```

!!! tip "Edit the config file manually"
    SnipKit has a lot more configuration options. Please see [Configuration][configuration] if you encounter problems or want
    to modify the behavior.

As of now, no external snippet manager is configured.

```sh title="Add an external snippet manager"
snipkit manager add
```

You will be presented with a list of supported managers. Pick the one you want to use. After that, you should be ready to go.

!!! tip "Different manager need different configuration"
    Every manager has unique configuration options. Have a look at [Managers][managers] for more information.

[configuration]: ../configuration/overview.md
[managers]: ../managers/overview.md

## Snippet Commands

#### Execute snippets

```sh title="Execute a snippet"
snipkit exec
```

!!! tip "Confirm commands"
    If you want to confirm a command before execution (with all parameters being resolved) add the 
    flag `--confirm`:
    ```bash 
    snipkit exec --confirm
    ```
    Snpkit will print the command on stdout and ask you to explicitly confirm its execution.

!!! tip "Print snippet on stdout"
    If you want to print the command that is executed add the flag `-p` or `--print`.

#### Print snippets

You can print snippets to stdout without executing them.

```sh title="Print a snippet"
snipkit print
```

#### Browse snippets

You can browse all available snippets without executing or printing them.

```sh title="Browse all snippets"
snipkit browse
```

#### Copy snippet to clipboard

You can copy a snippet to the clipboard in two ways:

```sh title="Copy to clipboard"
snipkit copy # Copies the snippet directly to the clipboard without printing
snipkit print --copy # Prints the snippet on stdout and, additionally, copies it to the clipboard
```

#### Export snippets

```bash
$ snipkit export -h

Exports all snippets on stdout as JSON including parsed meta information like parameters.

Usage:
  snipkit export [flags]

Flags:
  -f, --fields strings   Fields to be exported (default [id,title,content,parameters])
  -o, --output string    Output format. One of: json,json-pretty,xml (default "json")
```

The number of exported fields can be limited with the `--fields` flag. For example `snipkit export --fields=title,content -o json-pretty` results in an export like

```json 
{
  "snippets": [

    {
      "title": "List files in directory",
      "content": "ls -l"
    },
    {
      "title": "Check system uptime",
      "content": "uptime"
    }
  ]
}
```
