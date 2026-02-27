# Power setup

::: tip Customize Snipkit
Also have a look at [fzf][fzf] to get an understanding of how to customize Snipkit even more to fit your needs.
:::

### Alias

Always typing the full name `snipkit` in order to open the manager might be too
cumbersome for you. Instead, define an alias (e.g. in your `.zshrc` file):

```bash 
# SnipKit alias
sn () {
  snipkit exec "$@"
}
```

Then you can just type `sn` instead of `snipkit` to execute a script via SnipKit.

### Inline command for ZSH

The `print -z` command in Zsh is used to push a command onto the Zsh input buffer, which effectively allows you to 
simulate typing a command into the terminal. 

The specified command appears as if you had typed it at the prompt, but it's not executed immediately; instead, it 
waits for you to press Enter. This can be used as an alternative to SnipKit confirmation mechanism (via the 
`--confirm` flag). For ease of convenience, define another alias:

```bash 
# SnipKit alias
sn () {
  print -z $(snipkit print)
}
```

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

[fzf]: ./fzf
