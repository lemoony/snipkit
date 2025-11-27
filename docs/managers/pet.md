# Pet

Available for: macOS, Linux

[Repository](https://github.com/knqyf263/pet)

## Configuration

The configuration for Pet may look similar to this:

```yaml title="config.yaml" 
manager:
    pet:
      # Set to true if you want to use pet.
      enabled: true
      # List of pet snippet files.
      libraryPaths:
        - /Users/testuser/.config/pet/snippet.toml
      # If this list is not empty, only those snippets that match the listed tags will be provided to you.
      includeTags:
        - snipkit
        - othertag
```

Upon adding Pet as a manager, SnipKit will try to detect a default `libraryPath` automatically. If the library file was not
found, `enabled` will be set to `false`.

With this example configuration, SnipKit gets all snippets from Pet which are tagged `snipkit` or `othertag`. All other
snippets will not be presented to you. If you don't want to filter for tags, set `includeTags: []`.

## Parameter

Pet comes with its own parameter syntax in the form of `<param>`, `<param=default_value>` or `<param=|_value1_||_value2_|>`. 
SnipKit supports this syntax and you should have no problems using your Pet snippets the same way in SnipKit.

!!! tip
    While being easy to use, Pet's parameter syntax is less expressive than the one of SnipKit.
    Migrate to the  [file system directory][fslibrary] manager if you want to take advantage of the additional SnipKit 
    features like multiple value options or parameter descriptions.


[fslibrary]: ./fslibrary.md
