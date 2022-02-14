package tagutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

func Test_HasValidTags(t *testing.T) {
	tests := []struct {
		name       string
		validTags  []string
		actualTags []string
		expected   bool
	}{
		{name: "no tags at all", validTags: []string{}, actualTags: []string{}, expected: true},
		{name: "no valid tags but actual tags", validTags: []string{}, actualTags: []string{"foo"}, expected: true},
		{name: "valid tags not found", validTags: []string{"foo"}, actualTags: []string{}, expected: false},
		{name: "multiple valid tags", validTags: []string{"foo", "zoo"}, actualTags: []string{"zoo"}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validTagSet := stringutil.NewStringSet(tt.validTags)
			assert.Equal(t, tt.expected, HasValidTag(validTagSet, tt.actualTags))
		})
	}
}
