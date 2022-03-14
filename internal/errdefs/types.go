package errdefs

// InvalidArgument signals the client specified an invalid argument value
type InvalidArgument interface {
	InvalidArgument()
}

// Internal error signals a severe error that occurred in the computation
type Internal interface {
	Internal()
}
