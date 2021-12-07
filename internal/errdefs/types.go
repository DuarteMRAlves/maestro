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

// FailedPrecondition signals the system is not in a state capable of executing
// the desired operation. The client should wait for the system to be fixed
// before retrying.
type FailedPrecondition interface {
	FailedPrecondition()
}

// Unavailable signals some subsystem is currently unavailable. The client
// should wait and retry later.
type Unavailable interface {
	Unavailable()
}

// Internal error signals a severe error that occurred in the computation
type Internal interface {
	Internal()
}

// Unknown error signals that the king of error is unknown
type Unknown interface {
	Unknown()
}
