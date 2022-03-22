package execute

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestQueue(t *testing.T) {
	var (
		retrieved []state
		expected  []state
	)
	capacity := 10
	q := newQueue(capacity)

	for i := 0; i < 2*capacity; i++ {
		s := state{id: id(i), msg: nil}
		q.push(s)
		expected = append(expected, s)
	}

	// Only the last capacity elements should be kept.
	expected = expected[len(expected)-capacity:]

	for i := 0; i < capacity; i++ {
		retrieved = append(retrieved, q.pop())
	}

	cmpOps := cmp.AllowUnexported(state{})
	if diff := cmp.Diff(expected, retrieved, cmpOps); diff != "" {
		t.Fatalf("elements mismatch:\n%s", expected)
	}
}
