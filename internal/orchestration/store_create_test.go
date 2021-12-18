package orchestration

import (
	"gotest.tools/v3/assert"
	"testing"
)

const oName = "Orchestration Name"

func TestStore_CreateCorrect(t *testing.T) {
	tests := []struct {
		name   string
		config *Orchestration
	}{
		{
			name:   "test no links variable",
			config: &Orchestration{Name: oName},
		},
		{
			name:   "test nil links",
			config: &Orchestration{Name: oName, Links: nil},
		},
		{
			name:   "test empty links",
			config: &Orchestration{Name: oName, Links: []string{}},
		},
		{
			name:   "test links with one element",
			config: &Orchestration{Name: oName, Links: []string{"link-1"}},
		},
		{
			name: "test links with multiple elements",
			config: &Orchestration{
				Name:  oName,
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
				assert.Equal(t, 1, lenOrchestrations(st), "store size")

				stored, ok := st.orchestrations.Load(oName)
				assert.Assert(t, ok, "orchestration exists")

				o, ok := stored.(*Orchestration)
				assert.Assert(t, ok, "orchestration type assertion failed")
				assert.Equal(t, test.config.Name, o.Name, "correct name")
				if test.config.Links == nil {
					assert.DeepEqual(t, []string{}, o.Links)
				} else {
					assert.DeepEqual(t, test.config.Links, o.Links)
				}
			})
	}
}

func lenOrchestrations(st *store) int {
	count := 0
	st.orchestrations.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
}
