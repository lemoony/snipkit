# Overview

Managers are the actual provider of snippets.

## Supported managers

- [SnippetsLab](https://www.renfei.org/snippets-lab/)
- [Snip](https://github.com/Pictarine/macos-snippets)
- GitHub Gist ([Example gist](https://gist.github.com/lemoony/4905e7468b8f0a7991d6122d7d09e40d))
- [Pet](https://github.com/knqyf263/pet)
- [MassCode](https://masscode.io/)

Moreover, SnipKit allows you to provide snippets via a simple [file system directory][fslibrary].

## Adding a manager

Adding a manager means that SnipKit will retrieve snippets from it each time it is started. 

This command lets you add a manager to your [configuration][configuration]:

```sh
snipkit manager add
```

It represents a list of all supported managers that have not been added or enabled in your configuration. It will try to
detect the path to the manager and configure everything automatically. If SnipKit thinks it has found the 
corresponding manager and everything looks good so far, it will be enabled. Otherwise, all required config options will
be added to your config file, however, the manager will be disabled.


## Enabling & Disabling

Each manager can be enabled or disabled. By default, all managers are disabled:

```yaml title="config.yaml"
manager:
    <managerName>:
      # If set to false, the <managerName> is disabled 
      enabled: true
```

If a manager does not work, SnipKit refuses to startup. In this case, disable the manager by setting `enabled: false` or 
fix the configuration.


[configuration]: ../configuration/overview.md
[fslibrary]: ./fslibrary.md
