package errdefs

// AlreadyExists error signals that some resource already exists
type AlreadyExists interface {
	AlreadyExists()
}

// NotFound error signals that some resource does not exist
type NotFound interface {
	NotFound()
}

// Unknown error signals that the king of error is unknown
type Unknown interface {
	Unknown()
}
