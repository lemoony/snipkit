package model

var SyncNotSupported = SyncResult{NotSupported: true}

type SyncResult struct {
	Added   int
	Updated int
	Deleted int

	NotSupported bool

	Error error
}
