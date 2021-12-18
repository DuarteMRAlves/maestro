package orchestration

import (
	"gotest.tools/v3/assert"
	"testing"
)

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
