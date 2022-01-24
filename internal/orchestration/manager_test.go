package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_Create(t *testing.T) {
	const name apitypes.OrchestrationName = "Orchestration-Name"
	var (
		orchestration Orchestration
		err           error
	)
	cfg := &apitypes.Orchestration{Name: name}

	m, ok := NewManager().(*manager)
	assert.Assert(t, ok, "type assertion failed for manager")

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return m.CreateOrchestration(txn, cfg)
	})
	assert.NilError(t, err, "create error not nil")

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(orchestrationKey(name))
		assert.NilError(t, err, "get error")
		cp, err := item.ValueCopy(nil)
		return loadOrchestration(&orchestration, cp)
	})
	assert.NilError(t, err, "load error")
	assert.Equal(t, orchestration.Name(), cfg.Name, "name not correct")
	phase := orchestration.phase
	assert.Equal(t, phase, apitypes.OrchestrationPending, "phase not correct")
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
				var received []*apitypes.Orchestration

				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				m := NewManager()

				for _, o := range test.stored {
					err = db.Update(func(txn *badger.Txn) error {
						return persistOrchestration(txn, o)
					})
				}

				err = db.View(func(txn *badger.Txn) error {
					received, err = m.GetMatchingOrchestration(txn, test.query)
					return err
				})
				assert.NilError(t, err, "get orchestration")
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
