package model

type SyncState int

const (
	SyncStateStarted  = SyncState(1)
	SyncStateFinished = SyncState(2)
)

type SyncEvent struct {
	State SyncState
	Error error
}

type SyncFeedback struct {
	Events chan SyncEvent
}
