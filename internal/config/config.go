package config

import (
	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
)

type versionWrapper struct {
	Version string
	Config  Config
}

type Config struct {
	Style struct {
		Theme string
	}
	Providers struct {
		SnippetsLab snippetslab.Config
		FileSystem  filesystem.Config
	}
}
