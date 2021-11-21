package errdefs

// AlreadyExists error signals that some resource already exists
type AlreadyExists interface {
	AlreadyExists()
}

// NotFound error signals that some resource does not exist
type NotFound interface {
	NotFound()
}

// InvalidArgument signals the client specified an invalid argument value
type InvalidArgument interface {
	InvalidArgument()
}

// Internal error signals a severe error that occurred in the computation
type Internal interface {
	Internal()
}

// Unknown error signals that the king of error is unknown
type Unknown interface {
	Unknown()
}
