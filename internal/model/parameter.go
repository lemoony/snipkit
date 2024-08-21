package model

type ParameterType int

const (
	ParameterTypeValue    = ParameterType(0)
	ParameterTypePath     = ParameterType(1)
	ParameterTypePassword = ParameterType(2)
)

type Parameter struct {
	Key          string
	Name         string
	Type         ParameterType
	Description  string
	DefaultValue string
	Values       []string
}

type ParameterValue struct {
	Key   string
	Value string
}
