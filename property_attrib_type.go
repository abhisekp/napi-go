package napi

type PropertyAttributes int

const (
	Default           PropertyAttributes = 0
	Writable                             = 1 << 0
	Enumerable                           = 1 << 1
	Configurable                         = 1 << 2
	Static                               = 1 << 10
	DefaultMethod                        = Writable | Configurable
	DefaultJSProperty                    = Writable | Enumerable | Configurable
)
