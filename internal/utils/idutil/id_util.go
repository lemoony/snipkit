package idutil

import (
	"encoding/base64"
	"fmt"
)

type IDPrefix string

func FormatSnippetID(id string, prefix IDPrefix) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s#%s", prefix, id)))
}
