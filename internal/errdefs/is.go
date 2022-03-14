package errdefs

func IsInternal(err error) bool {
	_, ok := getImplementer(err).(Internal)
	return ok
}

type causer interface {
	Cause() error
}

func getImplementer(err error) error {
	switch e := err.(type) {
	case Internal:
		return err
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	}
}
