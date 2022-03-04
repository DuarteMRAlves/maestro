package orchestration

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_CreateOrchestration(t *testing.T) {
	const name api.OrchestrationName = "Orchestration-Name"
	var (
		orchestration api.Orchestration
		err           error
	)
	req := &api.CreateOrchestrationRequest{Name: name}

	db := kv.NewTestDb(t)
	defer db.Close()

	err = db.Update(
		func(txn *badger.Txn) error {
			createOrchestration := CreateOrchestrationWithTxn(txn)
			return createOrchestration(req)
		},
	)
	assert.NilError(t, err, "create error not nil")

	err = db.View(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			return helper.LoadOrchestration(&orchestration, name)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, orchestration.Name, req.Name, "name not correct")
	phase := orchestration.Phase
	assert.Equal(t, phase, api.OrchestrationPending, "phase not correct")
}

func TestManager_GetMatchingOrchestrations(t *testing.T) {
	tests := []struct {
		name   string
		req    *api.GetOrchestrationRequest
		stored []*api.Orchestration
		// names of the expected orchestrations
		expected []api.OrchestrationName
	}{
		{
			name:     "zero elements stored, nil req",
			req:      nil,
			stored:   []*api.Orchestration{},
			expected: []api.OrchestrationName{},
		},
		{
			name:     "zero elements stored, some req",
			req:      &api.GetOrchestrationRequest{Name: "some-name"},
			stored:   []*api.Orchestration{},
			expected: []api.OrchestrationName{},
		},
		{
			name: "one element stored, nil req",
			req:  nil,
			stored: []*api.Orchestration{
				orchestrationForName("some-name", api.OrchestrationFailed),
			},
			expected: []api.OrchestrationName{"some-name"},
		},
		{
			name: "one element stored, matching name req",
			req:  &api.GetOrchestrationRequest{Name: "some-name"},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name",
					api.OrchestrationRunning,
				),
			},
			expected: []api.OrchestrationName{"some-name"},
		},
		{
			name: "one element stored, non-matching name req",
			req:  &api.GetOrchestrationRequest{Name: "unknown-name"},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name",
					api.OrchestrationPending,
				),
			},
			expected: []api.OrchestrationName{},
		},
		{
			name: "multiple elements stored, nil req",
			req:  nil,
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name-1",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-2",
					api.OrchestrationSucceeded,
				),
				orchestrationForName(
					"some-name-3",
					api.OrchestrationFailed,
				),
			},
			expected: []api.OrchestrationName{
				"some-name-1",
				"some-name-2",
				"some-name-3",
			},
		},
		{
			name: "multiple elements stored, matching name req",
			req:  &api.GetOrchestrationRequest{Name: "some-name-2"},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name-1",
					api.OrchestrationRunning,
				),
				orchestrationForName(
					"some-name-2",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-3",
					api.OrchestrationFailed,
				),
			},
			expected: []api.OrchestrationName{"some-name-2"},
		},
		{
			name: "multiple elements stored, non-matching name req",
			req:  &api.GetOrchestrationRequest{Name: "unknown-name"},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name-1",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-2",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-3",
					api.OrchestrationRunning,
				),
			},
			expected: []api.OrchestrationName{},
		},
		{
			name: "multiple elements stored, matching phase req",
			req: &api.GetOrchestrationRequest{
				Phase: api.OrchestrationFailed,
			},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name-1",
					api.OrchestrationRunning,
				),
				orchestrationForName(
					"some-name-2",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-3",
					api.OrchestrationFailed,
				),
			},
			expected: []api.OrchestrationName{"some-name-3"},
		},
		{
			name: "multiple elements stored, non-matching phase req",
			req: &api.GetOrchestrationRequest{
				Phase: api.OrchestrationSucceeded,
			},
			stored: []*api.Orchestration{
				orchestrationForName(
					"some-name-1",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-2",
					api.OrchestrationPending,
				),
				orchestrationForName(
					"some-name-3",
					api.OrchestrationRunning,
				),
			},
			expected: []api.OrchestrationName{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					err      error
					received []*api.Orchestration
				)

				db := kv.NewTestDb(t)
				defer db.Close()

				for _, o := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							helper := kv.NewTxnHelper(txn)
							return helper.SaveOrchestration(o)
						},
					)
				}

				err = db.View(
					func(txn *badger.Txn) error {
						getOrchestration := GetOrchestrationsWithTxn(txn)
						received, err = getOrchestration(test.req)
						return err
					},
				)
				assert.NilError(t, err, "get orchestration")
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[api.OrchestrationName]bool, 0)
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
			},
		)
	}
}

func orchestrationForName(
	name api.OrchestrationName,
	phase api.OrchestrationPhase,
) *api.Orchestration {
	return &api.Orchestration{
		Name:  name,
		Phase: phase,
	}
}
