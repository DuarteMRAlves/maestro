package errdefs

// Internal error signals a severe error that occurred in the computation
type Internal interface {
	Internal()
}
