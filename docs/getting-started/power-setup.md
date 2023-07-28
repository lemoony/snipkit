# Power setup

### Alias

Always typing the full name `snipkit` in order to open the manager might be too
cumbersome for you. Instead, define an alias (e.g. in your `.zshrc` file):

```bash 
# SnipKit alais
sn () {
  snipkit $1
}
```

Then you can just type `sn` instead of `snipkit` to open SnipKit.

### Default Root Command

Most of the time, you want to call the same subcommand, e.g. `print` or `exec`. You
can configure `snipkit` so that this command gets executed by default by editing the config:

*Example:*

```yaml
# snipkit config edit 
defaultRootCommand: "exec"
```

With this setup, calling `sn` will yield the same result as `snipkit exec`. If you want to call
the `print` command instead, type `sn print`.

### Interactive ZSH Widget

Similar to the history search in ZSH, it is possible to bind `snipkit print` to a keybinding that will copy the generated snippet to the clipboard.

Define a function:

```shell
snipkit-snippets-widget () {
        echoti rmkx
        exec </dev/tty
        local snipkit_output=$(mktemp ${TMPDIR:-/tmp}/snipkit.output.XXXXXXXX)
        ./snipkit print -o "${snipkit_output}"
        echoti smkx
        cat $snipkit_output | pbcopy
        rm -f $snipkit_output
}
```

Bind widget to `CTRL+x x`:

```shell
bindkey "^Xx" snipkit-snippets-widget
```
