package fslibrary

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

func getSnippetName(system *system.System, filePath string) string {
	contents := string(system.ReadFile(filePath))
	if title, ok := titleheader.ParseTitleFromHeader(contents); ok {
		return title
	}
	return filepath.Base(filePath)
}

func pruneTitleHeader(r io.Reader) string {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return titleheader.PruneTitleHeader(string(all))
}
