package migrations

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const latestConfigPath = "../testdata/example-config.yaml"

func Test_Migrate(t *testing.T) {
	tests := []struct {
		from     string
		fromPath string
		to       string
		toPath   string
	}{
		{
			from:     "1.1.0",
			fromPath: configVersionPath("1.1.0"),
			to:       "1.1.1",
			toPath:   latestConfigPath,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Migrate_%s-%s", tt.from, tt.to), func(t *testing.T) {
			input, err := ioutil.ReadFile(tt.fromPath)
			assert.NoError(t, err)

			expected, err := ioutil.ReadFile(tt.toPath)
			assert.NoError(t, err)

			actual := Migrate(input)

			assert.YAMLEq(t, string(actual), string(expected))
		})
	}
}

func configVersionPath(version string) string {
	return fmt.Sprintf("../testdata/migrations/config-%s.yaml", strings.ReplaceAll(version, ".", "-"))
}
