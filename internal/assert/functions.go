package assert

import (
	"fmt"
)

func ArgNotNil(v interface{}, name string) (bool, error) {
	if v == nil {
		return false, fmt.Errorf("%v is nil", name)
	}
	return true, nil
}

func Status(b bool, msg string, msgArgs ...interface{}) (bool, error) {
	if !b {
		return false, fmt.Errorf(
			"status assertion failed: '%v'",
			fmt.Sprintf(msg, msgArgs...))
	}
	return true, nil
}
