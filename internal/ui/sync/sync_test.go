package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func Test_SyncScreen_inout_token(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Syncing all managers...")
		c.ExpectString("Manager X sync started")
		c.ExpectString("Please type in something")
		c.Send("test_token")
		c.SendKey(termtest.KeyEnter)

		c.ExpectString("All done.")
	}, func(stdio termutil.Stdio) {
		screen := New(WithIn(stdio.In), WithOut(stdio.Out))
		syncChannel := make(chan struct{})
		go func() {
			defer close(syncChannel)
			syncChannel <- struct{}{}
			screen.Start()
		}()
		<-syncChannel // wait for screen.Start()
		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusStarted})
		time.Sleep(time.Millisecond * 100)

		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Lines:  []appModel.SyncLine{{Type: appModel.SyncLineTypeInfo, Value: "Manager X sync started"}},
			},
		})
		time.Sleep(time.Millisecond * 100)

		inputChannel := make(chan appModel.SyncInputResult)
		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Input: &appModel.SyncInput{
					Content:     "Please type in something",
					Placeholder: "Type here...",
					Type:        appModel.SyncLoginTypeText,
					Input:       inputChannel,
				},
			},
		})

		assert.Equal(t, "test_token", (<-inputChannel).Text)

		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusFinished})
		<-syncChannel // wait for screen.Start() to return
	})
}

func Test_SyncScreen_input_pressKeyToContinue(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Syncing all managers...")
		c.ExpectString("Manager X sync started")
		c.ExpectString("Press [enter] to continue")
		c.SendKey(termtest.KeyEnter)

		c.ExpectString("All done.")
	}, func(stdio termutil.Stdio) {
		screen := New(WithIn(stdio.In), WithOut(stdio.Out))
		syncChannel := make(chan struct{})
		go func() {
			defer close(syncChannel)
			syncChannel <- struct{}{}
			screen.Start()
		}()
		<-syncChannel // wait for screen.Start()
		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusStarted})
		time.Sleep(time.Millisecond * 100)

		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Lines:  []appModel.SyncLine{{Type: appModel.SyncLineTypeInfo, Value: "Manager X sync started"}},
			},
		})
		time.Sleep(time.Millisecond * 100)

		inputChannel := make(chan appModel.SyncInputResult)
		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Input: &appModel.SyncInput{
					Content: "Press [enter] to continue",
					Type:    appModel.SyncLoginTypeContinue,
					Input:   inputChannel,
				},
			},
		})

		assert.True(t, (<-inputChannel).Continue)

		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusFinished})
		<-syncChannel // wait for screen.Start() to return
	})
}

func Test_SyncScreen_input_abort(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Syncing all managers...")
		c.ExpectString("Manager X sync started")
		c.ExpectString("Please type in something")
		c.SendKey(termtest.KeyStrC)
		c.ExpectString("Aborted")
	}, func(stdio termutil.Stdio) {
		screen := New(WithIn(stdio.In), WithOut(stdio.Out))
		syncChannel := make(chan struct{})
		go func() {
			defer close(syncChannel)
			syncChannel <- struct{}{}
			screen.Start()
		}()
		<-syncChannel // wait for screen.Start()
		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusStarted})
		time.Sleep(time.Millisecond * 100)

		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Lines:  []appModel.SyncLine{{Type: appModel.SyncLineTypeInfo, Value: "Manager X sync started"}},
			},
		})
		time.Sleep(time.Millisecond * 100)

		inputChannel := make(chan appModel.SyncInputResult)
		screen.Send(UpdateStateMsg{
			Status: appModel.SyncStatusStarted,
			ManagerState: &ManagerState{
				Status: appModel.SyncStatusStarted,
				Input: &appModel.SyncInput{
					Content: "Please type in something",
					Type:    appModel.SyncLoginTypeText,
					Input:   inputChannel,
				},
			},
		})

		input := <-inputChannel
		assert.True(t, input.Abort)
		assert.False(t, input.Continue)

		time.Sleep(time.Millisecond * 100)
		screen.Send(UpdateStateMsg{Status: appModel.SyncStatusAborted})
		<-syncChannel // wait for screen.Start() to return
	})
}
