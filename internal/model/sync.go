package model

type SyncEventChannel chan SyncEvent

type SyncState int

type SyncLineType int

const (
	SyncStateStarted  = SyncState(1)
	SyncStateFinished = SyncState(2)

	SyncLineTypeInfo    = SyncLineType(0)
	SyncLineTypeSuccess = SyncLineType(1)
	SyncLineTypeError   = SyncLineType(2)
)

type SyncEvent struct {
	State SyncState
	Lines []SyncLine
	Login *SyncLogin
	Error error
}

type SyncLogin struct {
	Title    string
	Content  string
	Continue chan struct{}
}

type SyncLine struct {
	Type  SyncLineType
	Value string
}
