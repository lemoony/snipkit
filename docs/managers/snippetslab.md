# SnippetsLab

Available for: macOS

[Homepeage](https://www.renfei.org/snippets-lab/)

## Configuration

The configuration for SnippetsLab may look similar to this:

```yaml title="config.yaml"
manager:
    snippetsLab:
      # Set to true if you want to use SnippetsLab.
      enabled: true
      # Path to your *.snippetslablibrary file.
      # SnipKit will try to detect this file automatically when generating the config.
      libraryPath: /path/to/main.snippetslablibrary
      # If this list is not empty, only those snippets that match the listed tags will be provided to you.
      includeTags:
        - snipkit
        - othertag
```

With this configuration, snipkit gets all snippets from SnippetsLab which are tagged with `snipkit` or `othertag`. All other
snippets will not be presented to you. If you don't want to filter for tags, set `includeTags: []`.

## Library path

Snipkit will try to automatically detect the path to the currently configured `*.snippetslablibrary` file.

If you have enabled iCloud sync, this path will be similar to:

```
/Users/<user>/Library/Containers/com.renfei.SnippetsLab/Data/Library/Application Support/com.renfei.SnippetsLab/main.snippetslablibrary
```
SnippetsLab lets you configure a custom library path. In this case, SnipKit will try to detect the preferences
file of SnippetsLab:

```
/Users/<user>/Library/Containers/com.renfei.SnippetsLab/Data/Library/Preferences/com.renfei.SnippetsLab.plist
```

The preferences file holds the path to the current `*.snippetslablibrary` file (if iCloud sync is turned off).
