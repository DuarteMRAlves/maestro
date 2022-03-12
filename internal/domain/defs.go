package domain

type Address interface {
	Address()
	Unwrap() string
}

type MethodContext interface {
	Address() Address
	Service() OptionalService
	Method() OptionalMethod
}

type LinkName interface {
	LinkName()
	Unwrap() string
}

type MessageField interface {
	MessageField()
	Unwrap() string
}

type OptionalMessageField interface {
	Unwrap() MessageField
	Present() bool
}
