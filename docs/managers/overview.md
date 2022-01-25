# Overview

Managers are the actual provider of snippetss.

## Supported managers

- [SnippetsLab](https://www.renfei.org/snippets-lab/)
- [Snip](https://github.com/Pictarine/macos-snippets)

Moreover, snipkit allows you to provide snippets via a simple [file system directory][fslibrary].

## Adding a manager

Adding a manager means that snipkit will retrieve snippets from it each time it is started. 

This command lets you add a manager to your [configuration][configuration]:

```sh
snipkit manager add
```

It represents a list of all supported managers that have not been added or enabled in your configuration. It will try to
automatically detect the manager on your system and configure everything automatically. If SnipKit thinks has found the 
corresponding manager and everythings looks good so far, it will be enabled. Otherwise, all required config options will
be added to your config file, however, the manager will be disabled.


## Enabling & Disabling

Each manager can be enabled or disabled. By default, all managers are disabled:

```yaml title="config.yaml"
manager:
    <managerName>:
      # If set to false, the <managerName> is disabled 
      enabled: true
```

If a manager does not work, snipkit refuses to startup. In this case, disable the manager by setting `enabled: false`.


[configuration]: ../configuration/overview.md
[fslibrary]: ./fslibrary.md
