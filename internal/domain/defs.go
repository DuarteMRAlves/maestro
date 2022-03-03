package domain

type AssetName interface{ Unwrap() string }

type Image interface{ Unwrap() string }

type Asset interface {
	Name() AssetName
	Image() Image
}

type StageName interface{ Unwrap() string }

type Service interface{ Unwrap() string }

type OptionalService interface {
	Unwrap() (Service, bool)
	Present() bool
	Empty() bool
}

type Method interface{ Unwrap() string }

type OptionalMethod interface {
	Unwrap() Method
	Present() bool
	Empty() bool
}

type Address interface{ Unwrap() string }

type Stage interface {
	Name() StageName
	Service() OptionalService
	Method() OptionalMethod
	Address() Address
}

type LinkName interface{ Unwrap() string }

type MessageField interface{ Unwrap() string }

type OptionalMessageField interface {
	Unwrap() MessageField
	Present() bool
	Empty() bool
}

type LinkEndpoint interface {
	Stage() StageName
	Field() OptionalMessageField
}

type Link interface {
	Name() LinkName
	Source() LinkEndpoint
	Target() LinkEndpoint
}

type OrchestrationName interface{ Unwrap() string }

type Orchestration interface {
	Name() OrchestrationName
	Stages() []Stage
	Links() []Link
}
