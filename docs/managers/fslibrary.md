# File System Library

Available for: macOS, Linux

The file system library manager lets you provide snippets via multiple local directories. Each directory contains
files which correspond to snippets.

## Configuration

The configuration for the file system library may look similar to this:

```yaml title="config.yaml"
manager:
  fsLibrary:
    # If set to false, the files specified via libraryPath will not be provided to you.
    enabled: true
    # Paths directories that hold snippets files. Each file must hold one snippet only.
    libraryPath:
      - /path/to/file/system/library
      - /another/path
    # Only files with endings which match one of the listed suffixes will be considered.
    suffixRegex:
      - .sh
    # If set to true, the files will not be parsed in advance. This means, only the filename can be used as the snippet name.
    lazyOpen: false
    # If set to true, the title comment will not be shown in the preview window.
    hideTitleInPreview: true
```

## Snippet Names

By default, the file name will be used as the snippet name. E.g., snippet `/another/path/count-character.sh` will be 
presented to you as `count-character.sh` in the lookup window.

However, SnipKit lets you also provide a different snippet name via a special comment syntax:

```sh
#
# <custom snippet name>
#
```

The _start_ of this comment must be within the first 3 lines of the file. If the first of three consecutive lines starting
with `#` is after line 3, the snippet name will be ignored.

E.g., a valid snippet named `Do something` instead of `example.sh` may look like this:

```sh linenums="1" title="example.sh"
#!/bin/bash

#
# Do something
#

echo "here we go..."
```

!!! attention "Open snippets lazily"
    This only works if `lazyOpen` is set to false since the snippet files must be parsed in advance before presenting
    the lookup window. If set to `true`, only the filename can be used as snippet name.

!!! tip "Hide the title comment"
    If you don't want to show the title header in the snippet preview window, set `hideTitleInPreview: true`.
    SnipKit will remove the title header.
