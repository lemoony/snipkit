package model

type ProviderInfo struct {
	Lines []ProviderLine
}

type ProviderLine struct {
	Key     string
	Value   string
	IsError bool
}
