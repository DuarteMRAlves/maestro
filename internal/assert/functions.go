package assert

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"reflect"
)

func ArgNotNil(v interface{}, name string) (bool, error) {
	if isNil(v) {
		return false, errdefs.InvalidArgumentWithMsg("'%v' is nil", name)
	}
	return true, nil
}

func ArgStatus(b bool, msg string, msgArgs ...interface{}) (bool, error) {
	if !b {
		return false, errdefs.InvalidArgumentWithMsg(msg, msgArgs...)
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

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(v).IsNil()
	}
	return false
}
