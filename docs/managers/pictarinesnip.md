# Snip

Available for: macOS

[Homepage](https://snip.picta-hub.io/)

[Repository](https://github.com/Pictarine/macos-snippets)

## Configuration

The configuration for Snip may look similar to this:

```yaml title="config.yaml" 
manager:
    pictarineSnip:
      # Set to true if you want to use Snip.
      enabled: true
      # Path to the snippets file.
      libraryPath: /Users/<user>/Library/Containers/com.pictarine.Snip/Data/Library/Application Support/Snip/snippets
      # If this list is not empty, only those snippets that match the listed tags will be provided to you.
      includeTags:
        - snipkit
        - othertag
```

Upon adding Snip as a manager, SnipKit will try to detect the `librayPath` automatically. If the library file was not
found, `enabled` will be set to `false`.

With this example configuration, SnipKit gets all snippets from Snip which are tagged `snipkit` or `othertag`. All other
snippets will not be presented to you. If you don't want to filter for tags, set `includeTags: []`.
