package domain

type AssetName interface {
	AssetName()
	Unwrap() string
}

type Image interface {
	Image()
	Unwrap() string
}

type OptionalImage interface {
	Unwrap() Image
	Present() bool
}

type Asset interface {
	Name() AssetName
	Image() OptionalImage
}

type StageName interface {
	StageName()
	Unwrap() string
}

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

type Stage interface {
	Name() StageName
	MethodContext() MethodContext
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

type OrchestrationName interface {
	OrchestrationName()
	Unwrap() string
}
