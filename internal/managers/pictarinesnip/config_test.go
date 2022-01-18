package pictarinesnip

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_findDefaultSnippetsLibrary(t *testing.T) {
	tests := []struct {
		name              string
		userContainersDir string
		expected          string
	}{
		{name: "found", userContainersDir: "testdata/userhome/Library/Containers", expected: testDataDefaultLibraryPath},
		{name: "not found", userContainersDir: "testdata/other", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := testutil.NewTestSystem(system.WithUserContainersDir(tt.userContainersDir))
			path := findDefaultSnippetsLibrary(s)
			assert.Equal(t, tt.expected, path)
		})
	}
}
