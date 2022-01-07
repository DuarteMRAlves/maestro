package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_CreateCorrect(t *testing.T) {
	const (
		name  apitypes.OrchestrationName = "Orchestration Name"
		phase                            = apitypes.OrchestrationRunning
	)
	tests := []struct {
		name   string
		config *Orchestration
	}{
		{
			name:   "test no links variable",
			config: &Orchestration{name: name, phase: phase},
		},
		{
			name:   "test nil links",
			config: &Orchestration{name: name, phase: phase, links: nil},
		},
		{
			name: "test empty links",
			config: &Orchestration{
				name:  name,
				phase: phase,
				links: []apitypes.LinkName{},
			},
		},
		{
			name: "test links with one element",
			config: &Orchestration{
				name:  name,
				phase: phase,
				links: []apitypes.LinkName{"link-1"},
			},
		},
		{
			name: "test links with multiple elements",
			config: &Orchestration{
				name:  name,
				phase: phase,
				links: []apitypes.LinkName{"link-1", "link-2", "link-3"},
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

				stored, ok := st.orchestrations.Load(name)
				assert.Assert(t, ok, "orchestration exists")

				o, ok := stored.(*Orchestration)
				assert.Assert(t, ok, "orchestration type assertion failed")
				assert.Equal(t, test.config.name, o.name, "correct name")
				assert.Equal(t, test.config.phase, o.phase, "correct phase")
				if test.config.links == nil {
					assert.DeepEqual(t, []apitypes.LinkName{}, o.links)
				} else {
					assert.DeepEqual(t, test.config.links, o.links)
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
		query  *apitypes.Orchestration
		stored []*Orchestration
		// names of the expected orchestrations
		expected []apitypes.OrchestrationName
	}{
		{
			name:     "zero elements stored, nil query",
			query:    nil,
			stored:   []*Orchestration{},
			expected: []apitypes.OrchestrationName{},
		},
		{
			name:     "zero elements stored, some query",
			query:    &apitypes.Orchestration{Name: "some-name"},
			stored:   []*Orchestration{},
			expected: []apitypes.OrchestrationName{},
		},
		{
			name:  "one element stored, nil query",
			query: nil,
			stored: []*Orchestration{
				orchestrationForName("some-name", apitypes.OrchestrationFailed),
			},
			expected: []apitypes.OrchestrationName{"some-name"},
		},
		{
			name:  "one element stored, matching name query",
			query: &apitypes.Orchestration{Name: "some-name"},
			stored: []*Orchestration{
				orchestrationForName("some-name", apitypes.OrchestrationRunning),
			},
			expected: []apitypes.OrchestrationName{"some-name"},
		},
		{
			name:  "one element stored, non-matching name query",
			query: &apitypes.Orchestration{Name: "unknown-name"},
			stored: []*Orchestration{
				orchestrationForName("some-name", apitypes.OrchestrationPending),
			},
			expected: []apitypes.OrchestrationName{},
		},
		{
			name:  "multiple elements stored, nil query",
			query: nil,
			stored: []*Orchestration{
				orchestrationForName("some-name-1", apitypes.OrchestrationPending),
				orchestrationForName("some-name-2", apitypes.OrchestrationSucceeded),
				orchestrationForName("some-name-3", apitypes.OrchestrationFailed),
			},
			expected: []apitypes.OrchestrationName{
				"some-name-1",
				"some-name-2",
				"some-name-3",
			},
		},
		{
			name:  "multiple elements stored, matching name query",
			query: &apitypes.Orchestration{Name: "some-name-2"},
			stored: []*Orchestration{
				orchestrationForName("some-name-1", apitypes.OrchestrationRunning),
				orchestrationForName("some-name-2", apitypes.OrchestrationPending),
				orchestrationForName("some-name-3", apitypes.OrchestrationFailed),
			},
			expected: []apitypes.OrchestrationName{"some-name-2"},
		},
		{
			name:  "multiple elements stored, non-matching name query",
			query: &apitypes.Orchestration{Name: "unknown-name"},
			stored: []*Orchestration{
				orchestrationForName("some-name-1", apitypes.OrchestrationPending),
				orchestrationForName("some-name-2", apitypes.OrchestrationPending),
				orchestrationForName("some-name-3", apitypes.OrchestrationRunning),
			},
			expected: []apitypes.OrchestrationName{},
		},
		{
			name:  "multiple elements stored, matching phase query",
			query: &apitypes.Orchestration{Phase: apitypes.OrchestrationFailed},
			stored: []*Orchestration{
				orchestrationForName("some-name-1", apitypes.OrchestrationRunning),
				orchestrationForName("some-name-2", apitypes.OrchestrationPending),
				orchestrationForName("some-name-3", apitypes.OrchestrationFailed),
			},
			expected: []apitypes.OrchestrationName{"some-name-3"},
		},
		{
			name:  "multiple elements stored, non-matching phase query",
			query: &apitypes.Orchestration{Phase: apitypes.OrchestrationSucceeded},
			stored: []*Orchestration{
				orchestrationForName("some-name-1", apitypes.OrchestrationPending),
				orchestrationForName("some-name-2", apitypes.OrchestrationPending),
				orchestrationForName("some-name-3", apitypes.OrchestrationRunning),
			},
			expected: []apitypes.OrchestrationName{},
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

				seen := make(map[apitypes.OrchestrationName]bool, 0)
				for _, e := range test.expected {
					seen[e] = false
				}

				for _, r := range received {
					alreadySeen, exists := seen[r.name]
					assert.Assert(t, exists, "element should be expected")
					// Elements can't be seen twice
					assert.Assert(t, !alreadySeen, "element already seen")
					seen[r.name] = true
				}

				for _, e := range test.expected {
					// All elements should be seen
					assert.Assert(t, seen[e], "element not seen")
				}
			})
	}
}

func orchestrationForName(
	name apitypes.OrchestrationName,
	phase apitypes.OrchestrationPhase,
) *Orchestration {
	return &Orchestration{
		name:  name,
		phase: phase,
		links: []apitypes.LinkName{"link-1", "link-2"},
	}
}
