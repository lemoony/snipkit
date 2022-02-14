# GitHub Gist

Available for: macOS, Linux

The GitHub Gist manager lets you provide snippets via multiple GitHub accounts. Each gist may contain multiple files which 
are mapped to single snippets. The gists are cached locally and synchronized manually, so accessing them is very fast.

!!! tip "Example Gist"
    Upon adding the GitHub Gist manager via `snipkit manager add`, SnipKit will configure a working 
    [example gist](https://gist.github.com/lemoony/4905e7468b8f0a7991d6122d7d09e40d), so you can quickly see how it works.

## Configuration

The configuration for the GitHub Gist manager may look similar to this:

```yaml title="config.yaml"
manager:
    githubGist:
      # If set to false, github gist is disabled completely.
      enabled: true
      # You can define multiple independent GitHub Gist sources.
      gists:
        - # If set to false, this GitHub gist url is ignored.
          enabled: true
          # URL to the GitHub gist account.
          url: gist.github.com/lemoony
          # Supported values: None, OAuthDeviceFlow, PAT. Default value: None (which means no authentication). In order to retrieve secret gists, you must be authenticated.
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

Value `PAT` refers to a
[personalized access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token),
which has to be created manually. `OAuthDeviceFlow` refers to the OAuth Device Flow supported by GitHub. In both cases, 
the scope of the token is limited to `gist`.


After specifying the desired authentication mechanism, just trigger a synchronization via `snipkit sync`. Snipkit will
ask you for the PAT or perform the OAuth authorization. The access token will be stored securely (e.g., by means of Keychain 
on macOS).

!!! tip "PAT vs OAuth"
    If you want to use a PAT or OAuth. Both mechanisms have different advantages. With a personalized access token, you have 
    full control over the expiration date, but you have to create the token yourself. The OAuth mechanism is more convenient 
    since you don't have to create a token yourself. However, the token may expire sooner and you have to perform the 
    authentication process more often.

### Custom OAuth Client ID

One option which is not listed by default, since it won't be required very often, is the following:

```yaml title="config.yaml"
manager:
    githubGist:
      OAuthClientID: <client_id>
```

The `OAuthClientID` lets you specify a custom client ID when performing the OAuth authentication. This is only required if 
your gists are not hosted on github.com and you want to use `OAuthDeviceFlow`.

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
    SnipKit lets you also provide a different snippet name via a special comment syntax. For a more detailed
    description, please see section *Snippet Names* in [file system directory][fslibrary].
    In order to enable this feature, set `titleHeaderEnabled` to `true`. If a snippet does not contain a title header
    comment, the specified `nameMode` will decide the snippet name.

## Tags

The description of a gist may contain multiple tags which can be used for filtering via the `includeTags` option.

E.g., a gist with the description `Example gist title #test #snipkit` is tagged with `test` and `snipkit`. 
If you have set `removeTagsFromDescription` to `true`, only `Example gist title` will be used as Snippet Name.

[fslibrary]: ./fslibrary.md

