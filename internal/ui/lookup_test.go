package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fuzzyMatcher(t *testing.T) {
	tests := []struct {
		str    string
		substr string
		match  bool
		ranges [][2]int
		score  int
	}{
		{
			"Find process listening to port",
			"find port",
			true,
			[][2]int{{0, 4}, {26, 30}},
			8,
		},
		{
			"Find process listening to port",
			"proc",
			true,
			[][2]int{{5, 9}},
			4,
		},
		{
			"Find process listening to port",
			"app",
			false,
			[][2]int{},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.str+" - "+tt.substr, func(t *testing.T) {
			ranges, score, match := fuzzyMatcher(tt.str, tt.substr)

			assert.Equal(t, tt.match, match)
			if tt.match {
				assert.Equal(t, tt.ranges, ranges)
			} else {
				assert.Len(t, tt.ranges, 0)
			}
			assert.Equal(t, tt.score, score)
		})
	}
}
