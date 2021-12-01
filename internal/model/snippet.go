package model

import (
	"fmt"
	"strings"
)

type Snippet struct {
	UUID     string
	Title    string
	Language string
	TagUUIDs []string
	Content  string
}

func (s Snippet) String() string {
	return fmt.Sprintf(
		"UUD: %s, Title: %s, Tags: [%s], Language: %s Content: %s",
		s.UUID,
		s.Title,
		strings.Join(s.TagUUIDs, ","),
		s.Language,
		s.Content,
	)
}
