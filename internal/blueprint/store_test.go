package blueprint

import (
	"errors"
	"fmt"
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
	"testing"
)

const bpName = "Blueprint Name"

func TestStore_CreateCorrect(t *testing.T) {
	tests := []*Blueprint{
		{Name: bpName},
		{Name: bpName, Stages: nil},
		{Name: bpName, Stages: []string{}},
		{Name: bpName, Links: nil},
		{Name: bpName, Links: []string{}},
	}
	for _, config := range tests {
		testName := fmt.Sprintf("config=%v", config)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				testing2.IsTrue(t, ok, "type assertion failed for store")

				err := st.Create(config)
				testing2.IsNil(t, err, "create error")
				testing2.DeepEqual(t, 1, lenBlueprints(st), "store size")
				stored, ok := st.blueprints.Load(bpName)
				testing2.IsTrue(t, ok, "blueprint exists")
				bp, ok := stored.(*Blueprint)
				testing2.IsTrue(t, ok, "blueprint type assertion failed")
				testing2.DeepEqual(t, bpName, bp.Name, "correct name")
				testing2.DeepEqual(t, 0, len(bp.Stages), "empty Stages")
				testing2.DeepEqual(t, 0, len(bp.Links), "empty Links")
			})
	}
}

func TestStore_CreateIncorrect(t *testing.T) {
	tests := []struct {
		config *Blueprint
		err    error
	}{
		{nil, errors.New("nil config")},
		{
			&Blueprint{Name: bpName, Stages: []string{"Some Stage"}},
			errors.New("blueprint should not have Stages"),
		},
		{
			&Blueprint{Name: bpName, Links: []string{"Some Link"}},
			errors.New("blueprint should not have Links"),
		},
	}
	for _, inner := range tests {
		config, err := inner.config, inner.err
		testName := fmt.Sprintf("config=%v,err='%v'", config, err)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				testing2.IsTrue(t, ok, "type assertion failed for store")

				e := st.Create(config)
				testing2.DeepEqual(t, err, e, "expected error")
				testing2.DeepEqual(t, 0, lenBlueprints(st), "store size")
			})
	}
}

func lenBlueprints(st *store) int {
	count := 0
	st.blueprints.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
}
