package form

import (
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/muesli/termenv"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// config holds the configuration for the form.
type config struct {
	colorProfile termenv.Profile
	fs           afero.Fs
	input        *io.Reader
	output       *io.Writer
	help         help.Model
	styler       style.Style
}

// Option is a functional option for configuring the form.
type Option interface {
	apply(c *config)
}

type optionFunc func(o *config)

func (f optionFunc) apply(o *config) {
	f(o)
}

// WithIn sets the input reader for the form.
func WithIn(input io.Reader) Option {
	return optionFunc(func(c *config) {
		c.input = &input
	})
}

// WithOut sets the output writer for the form.
func WithOut(out io.Writer) Option {
	return optionFunc(func(c *config) {
		c.output = &out
	})
}

// WithStyler sets the styler for the form.
func WithStyler(styler style.Style) Option {
	return optionFunc(func(c *config) {
		c.styler = styler
	})
}

// WithFS sets the filesystem for the form.
func WithFS(fs afero.Fs) Option {
	return optionFunc(func(c *config) {
		c.fs = fs
	})
}
