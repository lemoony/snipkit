package assistant

// Option configures a Manager.
type Option interface {
	apply(a *assistantImpl)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(a *assistantImpl)

func (f optionFunc) apply(provider *assistantImpl) {
	f(provider)
}

func withClientProvider(provider ClientProvider) Option {
	return optionFunc(func(a *assistantImpl) {
		a.provider = provider
	})
}
