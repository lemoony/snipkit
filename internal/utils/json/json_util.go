package jsonutil

import (
	"bytes"
	"encoding/json"
)

func CompactJSON(input []byte) string {
	var compactJSON bytes.Buffer
	err := json.Compact(&compactJSON, input)
	if err != nil {
		panic(err)
	}
	return compactJSON.String()
}
