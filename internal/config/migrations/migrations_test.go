package migrations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config/testdata"
)

func Test_Migrate(t *testing.T) {
	tests := []struct {
		from testdata.ConfigVersion
		to   testdata.ConfigVersion
	}{
		{
			from: testdata.ConfigV100,
			to:   testdata.ConfigV130,
		},
		{
			from: testdata.ConfigV120Providers,
			to:   testdata.ConfigV130Providers,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Migrate_%s-%s", tt.from, tt.to), func(t *testing.T) {
			input := testdata.ConfigBytes(t, tt.from)
			expected := testdata.ConfigBytes(t, tt.to)

			actual := Migrate(input)
			assert.YAMLEq(t, string(actual), string(expected))
		})
	}
}

func Test_Migrate_invalidYamlPanic(t *testing.T) {
	assert.Panics(t, func() {
		Migrate([]byte("{"))
	})
}

func Test_Migrate_invalidConfigVersion(t *testing.T) {
	assert.Panics(t, func() {
		Migrate([]byte("version: 3.0.0"))
	})
}
