# GitHub Gist

Available for: macOS, Linux

The GitHub Gist manager lets you provide snippets via multiple GitHub Gist accounts. The snippets are cached locally and
synchronized manually.

## Configuration

The configuration for the GitHub Gist manager may look similar to this:

```yaml title="config.yaml"
manager:
    githubGist:
      # If set to false, github gist is disabled completely.
      enabled: true
      # You can define multiple independent Github Gist sources.
      gists:
        - # If set to false, this github gist url is ignored.
          enabled: true
          # URL to the GitHub gist account.
          url: gist.github.com/lemoony
          # Supported values: None, OAuthDeviceFlow, Token. Default value: None (which means no authentication). In order to retrieve secret gists, you must be authenticated.
          authenticationMethod: OAuthDeviceFlow
          # If this list is not empty, only those gists that match the listed tags will be provided to you.
          includeTags: [snipkitExample]
          # Only gist files with endings which match one of the listed suffixes will be considered.
          suffixRegex: [.sh]
          # Defines where the snippet name is extracted from (see also titleHeaderEnabled). Allowed values: DESCRIPTION, FILENAME, COMBINE, COMBINE_PREFER_DESCRIPTION.
          nameMode: COMBINE_PREFER_DESCRIPTION
          # If set to true, any tags will be removed from the description.
          removeTagsFromDescription: true
          # If set to true, the snippet title can be overwritten by defining a title header within the gist.
          titleHeaderEnabled: true
          # If set to true, the title header comment will not be shown in the preview window.
          hideTitleInPreview: true
```

## Synchronization

All gists are cached locally. If there are updates, you have to manually trigger a synchronization
process via

```sh 
snipkit manager sync
```

There is also a shorthand alias:

```sh
snipkit sync
```

## Authentication

If `authenticationMethod` is set to `None`, only public gists are available. In order to retrieve secret gists,
use one of the following values:

- `PAT`
- `OAuthDeviceFlow`

Option `PAT` refers to a
[personalized access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token),
which has to be created manually. `OAuthDeviceFlow` refers to the OAuth Device Flow supported by GitHub. In both cases, 
the scope of the token is limited to `gist`.


After specifying the desired authentication mechanism, just trigger a synchronization via `snipkit sync`. Snipkit will
ask you for the PAT or redirect you. The access token will be stored securely (e.g. by means of Keychain on macOS).

!!! tip "PAT vs OAuth"
    If you want to use a PAT or the OAuth device flow is up to you. Both mechanisms have different advantages.
    E.g. with a personalized access token, you have full control over the expiration date. The OAuth mechanism,
    on the other hand, is more easy to use since you don't have to create a token yourself. With an OAuth token, however,
    you may have to perform the authentication process more often.


## Snippet Names

The name of a snippet can be retrieved by multiple ways. Set config option `nameMode` to one of 
the following values:

- `DESCRIPTION`: The description of the gist will be used as name. If a gist contains multiple files, all snippets will 
    have the same name.
- `FILENAME`: The filename of the gist will be used as name.
- `COMBINE`: The description and the filename of the gist will be concatenated (`<description> - <filename>`). 
- `COMBINE_PREFER_DESCRIPTION`:  The description of the gist will be used as name. If a gist contains multiple files, 
    the description and the filename of the gist will be concatenated.

!!! tip "Comment Syntax"
    Moreover, SnipKit lets you also provide a different snippet name via a special comment syntax. For a more detailed
    description, please see section *Snippet Names* in [file system directory][fslibrary].
    In order to enable this feature, set `titleHeaderEnabled` to `true`. If a snippet does not contain a title header
    comment, the specified `nameMode` will decide the snippet name.

## Tags

The description of a gist may contain multiple tags which can be used for filtering via the `includeTags` option.

E.g., a gist with the description `Example gist title #test #snipkit` is tagged with `test` and `snipkit`. 
If you have set `removeTagsFromDescription` to `true`, only `Example gist title` will be used as Snippet Name.

[fslibrary]: ./fslibrary.md

