package blueprint

import (
	"gotest.tools/v3/assert"
	"testing"
)

const bpName = "Blueprint Name"

func TestStore_CreateCorrect(t *testing.T) {
	tests := []struct {
		name   string
		config *Blueprint
	}{
		{
			name:   "test no links variable",
			config: &Blueprint{Name: bpName},
		},
		{
			name:   "test nil links",
			config: &Blueprint{Name: bpName, Links: nil},
		},
		{
			name:   "test empty links",
			config: &Blueprint{Name: bpName, Links: []string{}},
		},
		{
			name:   "test links with one element",
			config: &Blueprint{Name: bpName, Links: []string{"link-1"}},
		},
		{
			name: "test links with multiple elements",
			config: &Blueprint{
				Name:  bpName,
				Links: []string{"link-1", "link-2", "link-3"},
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(test.config)
				assert.NilError(t, err, "create error")
				assert.Equal(t, 1, lenBlueprints(st), "store size")

				stored, ok := st.blueprints.Load(bpName)
				assert.Assert(t, ok, "blueprint exists")

				bp, ok := stored.(*Blueprint)
				assert.Assert(t, ok, "blueprint type assertion failed")
				assert.Equal(t, test.config.Name, bp.Name, "correct name")
				if test.config.Links == nil {
					assert.DeepEqual(t, []string{}, bp.Links)
				} else {
					assert.DeepEqual(t, test.config.Links, bp.Links)
				}
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
