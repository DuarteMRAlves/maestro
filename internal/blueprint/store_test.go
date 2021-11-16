package blueprint

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

const bpName = "Blueprint Name"

func TestStore_CreateCorrect(t *testing.T) {
	tests := []*Blueprint{
		{Name: bpName},
		{Id: identifier.Empty(), Name: bpName},
		{Name: bpName, Stages: nil},
		{Name: bpName, Stages: []*Stage{}},
		{Name: bpName, Links: nil},
		{Name: bpName, Links: []*Link{}},
	}
	for _, config := range tests {
		testName := fmt.Sprintf("config=%v", config)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.IsTrue(t, ok, "type assertion failed for store")

				bpId, err := st.Create(config)
				assert.IsNil(t, err, "create error")
				assert.DeepEqual(t, IdSize, bpId.Size(), "blueprint id size")
				assert.DeepEqual(t, 1, lenBlueprints(st), "store size")
				stored, ok := st.blueprints.Load(bpId)
				assert.IsTrue(t, ok, "blueprint exists")
				bp, ok := stored.(*Blueprint)
				assert.IsTrue(t, ok, "blueprint type assertion failed")
				assert.DeepEqual(t, bpId, bp.Id, "correct id")
				assert.DeepEqual(t, bpName, bp.Name, "correct name")
				assert.DeepEqual(t, 0, len(bp.Stages), "empty Stages")
				assert.DeepEqual(t, 0, len(bp.Links), "empty Links")
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
			&Blueprint{Id: identifier.Id{Val: "OSNG132VSG"}, Name: bpName},
			errors.New("blueprint identifier should not be defined"),
		},
		{
			&Blueprint{Name: bpName, Stages: []*Stage{nil}},
			errors.New("blueprint should not have Stages"),
		},
		{
			&Blueprint{Name: bpName, Links: []*Link{nil}},
			errors.New("blueprint should not have Links"),
		},
	}
	for _, inner := range tests {
		config, err := inner.config, inner.err
		testName := fmt.Sprintf("config=%v,err='%v'", config, err)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.IsTrue(t, ok, "type assertion failed for store")

				bpId, e := st.Create(config)
				assert.DeepEqual(t, identifier.Empty(), bpId, "empty id")
				assert.DeepEqual(t, err, e, "expected error")
				assert.DeepEqual(t, 0, lenBlueprints(st), "store size")
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
