package idutil

import (
	"encoding/base64"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
)

type IDPrefix string

func FormatSnippetID(id string, prefix IDPrefix) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s#%s", prefix, id)))
}

func NanoID() string {
	result, _ := gonanoid.Nanoid()
	return result
}
