package form

import "io"

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
