package form

import (
	"io"

	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/ui/style"
)

type Option interface {
	apply(c *model)
}

type optionFunc func(o *model)

func (f optionFunc) apply(o *model) {
	f(o)
}

func WithIn(input io.Reader) Option {
	return optionFunc(func(c *model) {
		c.input = &input
	})
}

func WithOut(out io.Writer) Option {
	return optionFunc(func(c *model) {
		c.output = &out
	})
}

func WithStyler(styler style.Style) Option {
	return optionFunc(func(c *model) {
		c.styler = styler
	})
}

func WithFS(fs afero.Fs) Option {
	return optionFunc(func(c *model) {
		c.fs = fs
	})
}
