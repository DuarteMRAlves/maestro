package domain

type AssetName interface{ Unwrap() string }

type Image interface{ Unwrap() string }

type OptionalImage interface {
	Unwrap() Image
	Present() bool
}

type Asset interface {
	Name() AssetName
	Image() OptionalImage
}

type StageName interface{ Unwrap() string }

type Service interface{ Unwrap() string }

type OptionalService interface {
	Unwrap() Service
	Present() bool
}

type Method interface{ Unwrap() string }

type OptionalMethod interface {
	Unwrap() Method
	Present() bool
}

type Address interface{ Unwrap() string }

type MethodContext interface {
	Address() Address
	Service() OptionalService
	Method() OptionalMethod
}

type Stage interface {
	Name() StageName
	MethodContext() MethodContext
}

type LinkName interface{ Unwrap() string }

type MessageField interface{ Unwrap() string }

type OptionalMessageField interface {
	Unwrap() MessageField
	Present() bool
}

type OrchestrationName interface{ Unwrap() string }
