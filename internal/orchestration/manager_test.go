package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_CreateCorrect(t *testing.T) {
	const name apitypes.OrchestrationName = "Orchestration-Name"
	tests := []struct {
		name   string
		config *apitypes.Orchestration
	}{
		{
			name:   "test all fields",
			config: &apitypes.Orchestration{Name: name},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				st, ok := NewManager().(*manager)
				assert.Assert(t, ok, "type assertion failed for manager")

				err := st.CreateOrchestration(test.config)
				assert.NilError(t, err, "create error")
				assert.Equal(t, 1, lenOrchestrations(st), "manager size")

				stored, ok := st.orchestrations.Load(name)
				assert.Assert(t, ok, "orchestration exists")

				o, ok := stored.(*Orchestration)
				assert.Assert(t, ok, "orchestration type assertion failed")
				assert.Equal(t, test.config.Name, o.name, "correct name")
				assert.Equal(t, apitypes.OrchestrationPending, o.phase, "correct phase")
			})
	}
}

func lenOrchestrations(st *manager) int {
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
				st := NewManager()

				for _, o := range test.stored {
					st.CreateOrchestrationInternal(o)
				}

				received := st.GetMatchingOrchestration(test.query)
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[apitypes.OrchestrationName]bool, 0)
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

func orchestrationForName(
	name apitypes.OrchestrationName,
	phase apitypes.OrchestrationPhase,
) *Orchestration {
	return &Orchestration{
		name:  name,
		phase: phase,
	}
}
