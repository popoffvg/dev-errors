package errors

var (
	opts = options{
		withStack:  true,
		withFields: true,
		printer:    &defaultPrinter{},
	}
)

type (
	option func(opts options)
	hook   func(e error) error

	options struct {
		withStack  bool
		withFields bool
		printer    printer
		hook       hook
	}
)

func SetOptions(newOpts ...option) {
	for _, opt := range newOpts {
		opt(opts)
	}
}

func EnableStack(val bool) option {
	return func(opts options) {
		opts.withStack = val
	}
}

func EnableFields(val bool) option {
	return func(opts options) {
		opts.withFields = val
	}
}

func Verbose() option {
	return func(opts options) {
		opts.printer = &verbosePrinter{}
	}
}

func WithHook(h hook) option {
	return func(opts options) {
		opts.hook = h
	}
}
