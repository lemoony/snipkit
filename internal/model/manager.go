package model

type ManagerKey string

type ManagerInfo struct {
	Lines []ManagerInfoLine
}

type ManagerInfoLine struct {
	Key     string
	Value   string
	IsError bool
}

type ManagerDescription struct {
	Key         ManagerKey
	Name        string
	Description string
	Enabled     bool
}
