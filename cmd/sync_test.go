package cmd

import (
	"testing"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Sync(t *testing.T) {
	app := mocks.App{}
	app.On("SyncManager").Return(nil, nil)

	runExecuteTest(t, []string{"sync"}, withApp(&app))

	app.AssertNumberOfCalls(t, "SyncManager", 1)
}
