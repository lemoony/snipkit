package masscode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_AutoDiscoveryConfig(t *testing.T) {
	tests := []struct {
		name        string
		userHomeDir string
		expected    Config
	}{
		{
			name:        "found v1",
			userHomeDir: testDataUserHomeV1,
			expected: Config{
				Enabled:      true,
				MassCodeHome: fmt.Sprintf("%s/%s", testDataUserHomeV1, defaultMassCodeHomePath),
				Version:      version1,
				IncludeTags:  []string{},
			},
		},
		{
			name:        "found v2",
			userHomeDir: testDataUserHomeV2,
			expected: Config{
				Enabled:      true,
				MassCodeHome: fmt.Sprintf("%s/%s", testDataUserHomeV2, defaultMassCodeHomePath),
				Version:      version2,
				IncludeTags:  []string{},
			},
		},
		{
			name:        "not found",
			userHomeDir: "testdata/userhome-not-found",
			expected: Config{
				Enabled:      false,
				MassCodeHome: "testdata/userhome-not-found/massCode",
				Version:      version1,
				IncludeTags:  []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutil.NewTestSystem(system.WithUserHome(tt.userHomeDir))
			cfg := AutoDiscoveryConfig(s)
			assert.Equal(t, tt.expected, *cfg)
		})
	}
}
