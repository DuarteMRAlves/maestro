package identifier

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRand(t *testing.T) {
	tests := []struct {
		size int
		err  error
	}{
		{10, nil},
		{0, nil},
		{-1, fmt.Errorf("size is greater or equal than 0: %d", -1)},
	}
	for _, inner := range tests {
		size, err := inner.size, inner.err
		testName := fmt.Sprintf("size='%v',err='%v'", size, err)

		t.Run(testName, func(t *testing.T) {
			res, e := Rand(size)
			if !reflect.DeepEqual(err, e) {
				t.Fatalf("wrong error: expected '%v' but got '%v'", err, e)
			}
			if e != nil && res != Empty() {
				t.Fatalf("res!=nil and e!=nil: res='%v', e='%v'", res, e)
			}
			if res != Empty() && !reflect.DeepEqual(res.Size(), size) {
				t.Fatalf(
					"wrong size: expected '%v' but got '%v'",
					size,
					res.Size())
			}
		})
	}
}
