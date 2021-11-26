package testing

import (
	"fmt"
	"reflect"
	"testing"
)

func IsTrue(t *testing.T, actual bool, msg string, msgArgs ...interface{}) {
	if !actual {
		t.Fatalf("Value is false ('%v')", fmt.Sprintf(msg, msgArgs...))
	}
}

func IsNil(
	t *testing.T,
	actual interface{},
	msg string,
	msgArgs ...interface{},
) {
	if !isNil(actual) {
		t.Fatalf(
			"Value not nil: '%v' (msg='%v')",
			actual,
			fmt.Sprintf(msg, msgArgs...))
	}
}

func DeepEqual(
	t *testing.T,
	expected interface{},
	actual interface{},
	msg string,
	msgArgs ...interface{},
) {

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(
			"DeepEqual is false: expected='%v', actual='%v' (msg='%v')",
			expected,
			actual,
			fmt.Sprintf(msg, msgArgs...))
	}
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
