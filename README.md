# SnipKit - Snippet CLI manager

SnipKit aims to paste code snippets from your favorite snippet manager into your terminal without even leaving it.

<p align="center">
  <a href="https://go.dev/"><img alt="Language" src="https://img.shields.io/badge/language-Go-blue.svg"></a>
  <a href="https://github.com/lemoony/snipkit/actions/workflows/build.yml"><img alt="Build" src="https://github.com/lemoony/snipkit/actions/workflows/build.yml/badge.svg"></a>
  <a href="https://opensource.org/licenses/Apache-2.0"><img alt="License" src="https://img.shields.io/badge/License-Apache_2.0-blue.svg"></a>
  <a href="https://goreportcard.com/report/github.com/lemoony/snipkit"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lemoony/snipkit"></a>
  <a href="https://codecov.io/gh/lemoony/snipkit"><img alt="codecov" src="https://codecov.io/gh/lemoony/snipkit/branch/main/graph/badge.svg?token=UOG4O1yscP"></a>
</p>

<p align="center">
  <img alt="Demo" src="./docs/images/demo.gif" />
</p>

<p align="center">
<a href="ttps://lemoony.github.io/snipkit/">Documentation</a> 
</p>

## Features

`snipkit` supports the following features:

- Load snippets form an external snippet manager (filtered by tags)
  - SnippetsLab
  - File system directory
- Parameter substitution
- Enum parameters
- Search for snippets by typing
- Root command can be adjusted (e.g. set to `print` or `exec`)
- Themes
  - Built-in themes (`default`, `dracula`, `solarized-light`, `example`)
  - Define custom themes

Inspired by [Pet](https://github.com/knqyf263/pet).

## Quick Start

#### Overview of all commands

```bash
snipkit -h
```
#### Configuration

```bash 
# Create a new config
snipkit config init
```

If you have SnippetsLab installed, the config should already point to the corresponding
library file. 

You can open & edit the config file easily:

```bash 
snipkit config edit
```

Have a look at the various configuration options. They should be self-explanatory
most of the time.

## Power Setup

### Alias

Always typing the full name `snipkit` in order to open the manager might be too 
cumbersome for you. Just define an alias (e.g. in your `.zshrc` file):

```bash 
# SnipKit alais
sn () {
  snipkit $1
}
```

Then you can just type `sn` instead of `snipkit` to open the app.

### Default Root Command

Most of the time, you want to call the same subcommand, e.g. `print` or `exec`. You
can configure `snipkit` so that this command gets executed by default by editing the config:

*Example:*

```yaml
# snipkit config edit 
defaultRootCommand: "exec"
```

With this setup, calling `sn` will yield the same result as `snipkit exec`. If you want to call
the `print` command instead, just type `sn print`.

## Installation

### Homebrew

```bash 
brew install lemoony/tap/snipkit
```

### apt 

```bash 
echo 'deb [trusted=yes] https://apt.fury.io/lemoony/ /' | sudo tee /etc/apt/sources.list.d/snipkit.list
sudo apt update
sudo apt install snipkit
```

### yum

```bash 
echo '[snipkit]
name=Snipkit Private Repo
baseurl=https://yum.fury.io/lemoony/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/snipkit.repo
sudo yum install snipkit
```
### deb, rpm and apk packages 

Download the .deb, .rpm or .apk packages from [releases page](https://github.com/lemoony/snipkit/releases) and install 
them with the appropriate tools.


### Go

```bash
go install github.com/lemoony/snipkit@latest
```

### Build

```bash 
git clone https://github.com/lemoony/snipkit.git
cd snipkit 
make build
```

After the build succeeds, go to `./dist` to find the binary for your operating system.


## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).
