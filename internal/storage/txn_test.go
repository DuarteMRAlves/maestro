package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestTxnHelper_SaveOrchestration(t *testing.T) {
	tests := []struct {
		name          string
		orchestration *api.Orchestration
		expected      []byte
	}{
		{
			name: "default orchestration",
			orchestration: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{},
				Links:  []api.LinkName{},
			},
			expected: []byte("default;" + string(api.OrchestrationPending) + ";;"),
		},
		{
			name: "non default orchestration",
			orchestration: &api.Orchestration{
				Name:   "some-name",
				Phase:  api.OrchestrationFailed,
				Stages: []api.StageName{},
				Links:  []api.LinkName{},
			},
			expected: []byte("some-name;" + string(api.OrchestrationFailed) + ";;"),
		},
		{
			name: "orchestration with stages",
			orchestration: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{"stage-1", "stage-2"},
				Links:  []api.LinkName{},
			},
			expected: []byte("default;" + string(api.OrchestrationPending) + ";stage-1,stage-2;"),
		},
		{
			name: "orchestration with links",
			orchestration: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{},
				Links:  []api.LinkName{"link-1", "link-2"},
			},
			expected: []byte("default;" + string(api.OrchestrationPending) + ";;link-1,link-2"),
		},
		{
			name: "orchestration with stages and links",
			orchestration: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{"stage-1", "stage-2", "stage-3"},
				Links:  []api.LinkName{"link-1", "link-2"},
			},
			expected: []byte("default;" + string(api.OrchestrationPending) + ";stage-1,stage-2,stage-3;link-1,link-2"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var stored []byte

				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(
					func(txn *badger.Txn) error {
						helper := NewTxnHelper(txn)
						return helper.SaveOrchestration(test.orchestration)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(orchestrationKey(test.orchestration.Name))
						if err != nil {
							return err
						}
						stored, err = item.ValueCopy(nil)
						return err
					},
				)
				assert.Equal(t, len(test.expected), len(stored), "stored size")
				for i, e := range test.expected {
					assert.Equal(t, e, stored[i], "stored not equal")
				}
			},
		)
	}
}

func TestTxnHelper_LoadOrchestration(t *testing.T) {
	tests := []struct {
		name     string
		expected *api.Orchestration
		stored   []byte
	}{
		{
			name:   "default orchestration",
			stored: []byte("default;" + string(api.OrchestrationPending) + ";;"),
			expected: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{},
				Links:  []api.LinkName{},
			},
		},
		{
			name:   "non default orchestration",
			stored: []byte("some-name;" + string(api.OrchestrationFailed) + ";;"),
			expected: &api.Orchestration{
				Name:   "some-name",
				Phase:  api.OrchestrationFailed,
				Stages: []api.StageName{},
				Links:  []api.LinkName{},
			},
		},
		{
			name:   "orchestration with stages",
			stored: []byte("default;" + string(api.OrchestrationPending) + ";stage-1,stage-2;"),
			expected: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{"stage-1", "stage-2"},
				Links:  []api.LinkName{},
			},
		},
		{
			name:   "orchestration with links",
			stored: []byte("default;" + string(api.OrchestrationPending) + ";;link-1,link-2"),
			expected: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{},
				Links:  []api.LinkName{"link-1", "link-2"},
			},
		},
		{
			name:   "orchestration with stages and links",
			stored: []byte("default;" + string(api.OrchestrationPending) + ";stage-1,stage-2,stage-3;link-1,link-2"),
			expected: &api.Orchestration{
				Name:   "default",
				Phase:  api.OrchestrationPending,
				Stages: []api.StageName{"stage-1", "stage-2", "stage-3"},
				Links:  []api.LinkName{"link-1", "link-2"},
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded api.Orchestration
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							orchestrationKey(test.expected.Name),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						helper := NewTxnHelper(txn)
						return helper.LoadOrchestration(
							&loaded,
							test.expected.Name,
						)
					},
				)
				assert.NilError(t, err, "load error")
				assert.DeepEqual(t, test.expected, &loaded)
			},
		)
	}
}
