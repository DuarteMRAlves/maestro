package blueprint

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_Get_Correct(t *testing.T) {
	tests := []struct {
		name   string
		query  *Blueprint
		stored []*Blueprint
		// names of the expected blueprints
		expected []string
	}{
		{
			name:     "zero elements stored, nil query",
			query:    nil,
			stored:   []*Blueprint{},
			expected: []string{},
		},
		{
			name:     "zero elements stored, some query",
			query:    &Blueprint{Name: "some-name"},
			stored:   []*Blueprint{},
			expected: []string{},
		},
		{
			name:     "one element stored, nil query",
			query:    nil,
			stored:   []*Blueprint{blueprintForName("some-name")},
			expected: []string{"some-name"},
		},
		{
			name:     "one element stored, matching query",
			query:    &Blueprint{Name: "some-name"},
			stored:   []*Blueprint{blueprintForName("some-name")},
			expected: []string{"some-name"},
		},
		{
			name:     "one element stored, non-matching query",
			query:    &Blueprint{Name: "unknown-name"},
			stored:   []*Blueprint{blueprintForName("some-name")},
			expected: []string{},
		},
		{
			name:  "multiple elements stored, nil query",
			query: nil,
			stored: []*Blueprint{
				blueprintForName("some-name-1"),
				blueprintForName("some-name-2"),
				blueprintForName("some-name-3"),
			},
			expected: []string{"some-name-1", "some-name-2", "some-name-3"},
		},
		{
			name:  "multiple elements stored, matching query",
			query: &Blueprint{Name: "some-name-2"},
			stored: []*Blueprint{
				blueprintForName("some-name-1"),
				blueprintForName("some-name-2"),
				blueprintForName("some-name-3"),
			},
			expected: []string{"some-name-2"},
		},
		{
			name:  "multiple elements stored, non-matching query",
			query: &Blueprint{Name: "unknown-name"},
			stored: []*Blueprint{
				blueprintForName("some-name-1"),
				blueprintForName("some-name-2"),
				blueprintForName("some-name-3"),
			},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				st := NewStore()

				for _, bp := range test.stored {
					err := st.Create(bp)
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

func blueprintForName(name string) *Blueprint {
	return &Blueprint{
		Name:  name,
		Links: []string{"link-1", "link-2"},
	}
}
