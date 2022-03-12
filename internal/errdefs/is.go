package errdefs

func IsAlreadyExists(err error) bool {
	_, ok := getImplementer(err).(AlreadyExists)
	return ok
}

func IsInvalidArgument(err error) bool {
	_, ok := getImplementer(err).(InvalidArgument)
	return ok
}

func IsFailedPrecondition(err error) bool {
	_, ok := getImplementer(err).(FailedPrecondition)
	return ok
}

func IsUnavailable(err error) bool {
	_, ok := getImplementer(err).(Unavailable)
	return ok
}

func IsInternal(err error) bool {
	_, ok := getImplementer(err).(Internal)
	return ok
}

func IsUnknown(err error) bool {
	_, ok := getImplementer(err).(Unknown)
	return ok
}

type causer interface {
	Cause() error
}

func getImplementer(err error) error {
	switch e := err.(type) {
	case
		AlreadyExists,
		InvalidArgument,
		FailedPrecondition,
		Unavailable,
		Internal,
		Unknown:
		return err
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	}
}
