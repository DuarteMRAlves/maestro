package domain

type Service interface {
	Service()
	Unwrap() string
}

type OptionalService interface {
	Unwrap() Service
	Present() bool
}

type Method interface {
	Method()
	Unwrap() string
}

type OptionalMethod interface {
	Unwrap() Method
	Present() bool
}

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
