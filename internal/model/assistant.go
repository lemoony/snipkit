package model

type AssistantKey string

type AssistantDescription struct {
	Key         AssistantKey
	Name        string
	Description string
	Enabled     bool
}
