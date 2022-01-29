package model

type ManagerKey string

type ManagerDescription struct {
	Key         ManagerKey
	Name        string
	Description string
	Enabled     bool
}
