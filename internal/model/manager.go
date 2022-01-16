package model

type ManagerInfo struct {
	Lines []ManagerInfoLine
}

type ManagerInfoLine struct {
	Key     string
	Value   string
	IsError bool
}

type ManagerDescription struct {
	Name        string
	Description string
	Enabled     bool
}
