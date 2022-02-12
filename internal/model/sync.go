package model

type SyncEventChannel chan SyncEvent

type (
	SyncStatus    int
	SyncLineType  int
	SyncLoginType int
)

const (
	SyncStatusStarted  = SyncStatus(1)
	SyncStatusFinished = SyncStatus(2)
	SyncStatusAborted  = SyncStatus(3)

	SyncLineTypeInfo    = SyncLineType(0)
	SyncLineTypeSuccess = SyncLineType(1)
	SyncLineTypeError   = SyncLineType(2)

	SyncLoginTypeContinue = SyncLoginType(1)
	SyncLoginTypeText     = SyncLoginType(2)
)

type SyncEvent struct {
	Status SyncStatus
	Lines  []SyncLine
	Login  *SyncInput
	Error  error
}

type SyncInput struct {
	Content     string
	Placeholder string
	Type        SyncLoginType
	Input       chan SyncInputResult
}

type SyncLine struct {
	Type  SyncLineType
	Value string
}

type SyncInputResult struct {
	Continue bool
	Abort    bool
	Text     string
}
