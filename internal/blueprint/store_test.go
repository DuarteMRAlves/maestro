package blueprint

import (
	"fmt"
	"gotest.tools/v3/assert"
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
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(config)
				assert.NilError(t, err, "create error")
				assert.Equal(t, 1, lenBlueprints(st), "store size")
				stored, ok := st.blueprints.Load(bpName)
				assert.Assert(t, ok, "blueprint exists")
				bp, ok := stored.(*Blueprint)
				assert.Assert(t, ok, "blueprint type assertion failed")
				assert.Equal(t, bpName, bp.Name, "correct name")
				assert.Equal(t, 0, len(bp.Stages), "empty Stages")
				assert.Equal(t, 0, len(bp.Links), "empty Links")
			})
	}
}

func TestStore_CreateIncorrect(t *testing.T) {
	tests := []struct {
		config *Blueprint
		errMsg string
	}{
		{nil, "nil config"},
		{
			&Blueprint{Name: bpName, Stages: []string{"Some Stage"}},
			"blueprint should not have Stages",
		},
		{
			&Blueprint{Name: bpName, Links: []string{"Some Link"}},
			"blueprint should not have Links",
		},
	}
	for _, inner := range tests {
		config, errMsg := inner.config, inner.errMsg
		testName := fmt.Sprintf("config=%v,errMsg='%v'", config, errMsg)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(config)
				assert.ErrorContains(t, err, errMsg)
				assert.Equal(t, 0, lenBlueprints(st), "store size")
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
