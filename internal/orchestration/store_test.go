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

func TestStore_Get_Correct(t *testing.T) {
	tests := []struct {
		name   string
		query  *Orchestration
		stored []*Orchestration
		// names of the expected orchestrations
		expected []string
	}{
		{
			name:     "zero elements stored, nil query",
			query:    nil,
			stored:   []*Orchestration{},
			expected: []string{},
		},
		{
			name:     "zero elements stored, some query",
			query:    &Orchestration{Name: "some-name"},
			stored:   []*Orchestration{},
			expected: []string{},
		},
		{
			name:     "one element stored, nil query",
			query:    nil,
			stored:   []*Orchestration{orchestrationForName("some-name")},
			expected: []string{"some-name"},
		},
		{
			name:     "one element stored, matching query",
			query:    &Orchestration{Name: "some-name"},
			stored:   []*Orchestration{orchestrationForName("some-name")},
			expected: []string{"some-name"},
		},
		{
			name:     "one element stored, non-matching query",
			query:    &Orchestration{Name: "unknown-name"},
			stored:   []*Orchestration{orchestrationForName("some-name")},
			expected: []string{},
		},
		{
			name:  "multiple elements stored, nil query",
			query: nil,
			stored: []*Orchestration{
				orchestrationForName("some-name-1"),
				orchestrationForName("some-name-2"),
				orchestrationForName("some-name-3"),
			},
			expected: []string{"some-name-1", "some-name-2", "some-name-3"},
		},
		{
			name:  "multiple elements stored, matching query",
			query: &Orchestration{Name: "some-name-2"},
			stored: []*Orchestration{
				orchestrationForName("some-name-1"),
				orchestrationForName("some-name-2"),
				orchestrationForName("some-name-3"),
			},
			expected: []string{"some-name-2"},
		},
		{
			name:  "multiple elements stored, non-matching query",
			query: &Orchestration{Name: "unknown-name"},
			stored: []*Orchestration{
				orchestrationForName("some-name-1"),
				orchestrationForName("some-name-2"),
				orchestrationForName("some-name-3"),
			},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				st := NewStore()

				for _, o := range test.stored {
					err := st.Create(o)
					assert.Assert(t, err, "create asset error")
				}

				received := st.Get(test.query)
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[string]bool, 0)
				for _, e := range test.expected {
					seen[e] = false
				}

				for _, r := range received {
					alreadySeen, exists := seen[r.Name]
					assert.Assert(t, exists, "element should be expected")
					// Elements can't be seen twice
					assert.Assert(t, !alreadySeen, "element already seen")
					seen[r.Name] = true
				}

				for _, e := range test.expected {
					// All elements should be seen
					assert.Assert(t, seen[e], "element not seen")
				}
			})
	}
}

func orchestrationForName(name string) *Orchestration {
	return &Orchestration{
		Name:  name,
		Links: []string{"link-1", "link-2"},
	}
}