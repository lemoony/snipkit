package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func Test_Form_NoParameters(t *testing.T) {
	var parameters []model.Parameter

	term := NewTUI()
	values, ok := term.ShowParameterForm(parameters, OkButtonPrint)
	assert.Len(t, values, 0)
	assert.True(t, ok)
}
