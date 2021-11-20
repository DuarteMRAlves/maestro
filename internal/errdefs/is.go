package errdefs

func IsAlreadyExists(err error) bool {
	_, ok := getImplementer(err).(AlreadyExists)
	return ok
}

func IsNotFound(err error) bool {
	_, ok := getImplementer(err).(NotFound)
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
	case AlreadyExists, NotFound, Unknown:
		return err
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	}
}
