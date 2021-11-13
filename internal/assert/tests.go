package assert

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

	if actual != nil {
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

func NotDeepEqual(
	t *testing.T,
	expected interface{},
	actual interface{},
	msg string,
	msgArgs ...interface{},
) {

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf(
			"DeepEqual is true: expected='%v', actual='%v' (msg='%v')",
			expected,
			actual,
			fmt.Sprintf(msg, msgArgs...))
	}
}
