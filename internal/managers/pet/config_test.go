package pet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_AutoDiscoveryConfig(t *testing.T) {
	tests := []struct {
		name        string
		userHomeDir string
		enabled     bool
	}{
		{name: "found", userHomeDir: testDataUserHome, enabled: true},
		{name: "not found", userHomeDir: "testdata/not-found-dir", enabled: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutil.NewTestSystem(system.WithUserHome(tt.userHomeDir))
			cfg := AutoDiscoveryConfig(s)
			assert.Equal(t, tt.enabled, cfg.Enabled)
		})
	}
}
