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

func TestTxnHelper_SaveStage(t *testing.T) {
	tests := []struct {
		name     string
		stage    *api.Stage
		expected []byte
	}{
		{
			name: "default stage",
			stage: &api.Stage{
				Name:          "some-name",
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "",
				Orchestration: "default",
				Asset:         "",
			},
			expected: []byte("some-name;" + string(api.StagePending) + ";;;;default;"),
		},
		{
			name: "non default stage",
			stage: &api.Stage{
				Name:          "some-name",
				Phase:         api.StageRunning,
				Service:       "SomeService",
				Rpc:           "SomeRpc",
				Address:       "SomeAddress",
				Orchestration: "some-orchestration",
				Asset:         "some-asset",
			},
			expected: []byte("some-name;" + string(api.StageRunning) + ";SomeService;SomeRpc;SomeAddress;some-orchestration;some-asset"),
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
						return helper.SaveStage(test.stage)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(stageKey(test.stage.Name))
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

func TestTxnHelper_LoadStage(t *testing.T) {
	tests := []struct {
		name     string
		expected *api.Stage
		stored   []byte
	}{
		{
			name:   "default stage",
			stored: []byte("some-name;" + string(api.StagePending) + ";;;;default;"),
			expected: &api.Stage{
				Name:          "some-name",
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "",
				Orchestration: "default",
				Asset:         "",
			},
		},
		{
			name:   "non default stage",
			stored: []byte("some-name;" + string(api.StageRunning) + ";SomeService;SomeRpc;SomeAddress;some-orchestration;some-asset"),
			expected: &api.Stage{
				Name:          "some-name",
				Phase:         api.StageRunning,
				Service:       "SomeService",
				Rpc:           "SomeRpc",
				Address:       "SomeAddress",
				Orchestration: "some-orchestration",
				Asset:         "some-asset",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded api.Stage
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							stageKey(test.expected.Name),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						helper := NewTxnHelper(txn)
						return helper.LoadStage(
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

func TestTxnHelper_SaveLink(t *testing.T) {
	tests := []struct {
		name     string
		link     *api.Link
		expected []byte
	}{
		{
			name: "default link",
			link: &api.Link{
				Name:          "some-name",
				SourceStage:   "",
				SourceField:   "",
				TargetStage:   "",
				TargetField:   "",
				Orchestration: defaultOrchestrationName,
			},
			expected: []byte("some-name;;;;;" + defaultOrchestrationName),
		},
		{
			name: "non default link",
			link: &api.Link{
				Name:          "some-name",
				SourceStage:   "source-stage",
				SourceField:   "SourceField",
				TargetStage:   "target-stage",
				TargetField:   "TargetField",
				Orchestration: "some-orchestration",
			},
			expected: []byte("some-name;source-stage;SourceField;target-stage;TargetField;some-orchestration"),
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
						return helper.SaveLink(test.link)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(linkKey(test.link.Name))
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

func TestTxnHelper_LoadLink(t *testing.T) {
	tests := []struct {
		name     string
		expected *api.Link
		stored   []byte
	}{
		{
			name:   "default link",
			stored: []byte("some-name;;;;;" + defaultOrchestrationName),
			expected: &api.Link{
				Name:          "some-name",
				SourceStage:   "",
				SourceField:   "",
				TargetStage:   "",
				TargetField:   "",
				Orchestration: defaultOrchestrationName,
			},
		},
		{
			name:   "non default link",
			stored: []byte("some-name;source-stage;SourceField;target-stage;TargetField;some-orchestration"),
			expected: &api.Link{
				Name:          "some-name",
				SourceStage:   "source-stage",
				SourceField:   "SourceField",
				TargetStage:   "target-stage",
				TargetField:   "TargetField",
				Orchestration: "some-orchestration",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded api.Link
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							linkKey(test.expected.Name),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						helper := NewTxnHelper(txn)
						return helper.LoadLink(
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

func TestTxnHelper_SaveAsset(t *testing.T) {
	tests := []struct {
		name     string
		asset    *api.Asset
		expected []byte
	}{
		{
			name: "default asset",
			asset: &api.Asset{
				Name:  "some-name",
				Image: "",
			},
			expected: []byte("some-name;"),
		},
		{
			name: "non default stage",
			asset: &api.Asset{
				Name:  "some-name",
				Image: "some-image",
			},
			expected: []byte("some-name;some-image"),
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
						return helper.SaveAsset(test.asset)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(assetKey(test.asset.Name))
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

func TestTxnHelper_LoadAsset(t *testing.T) {
	tests := []struct {
		name     string
		expected *api.Asset
		stored   []byte
	}{
		{
			name:   "default asset",
			stored: []byte("some-name;"),
			expected: &api.Asset{
				Name:  "some-name",
				Image: "",
			},
		},
		{
			name:   "non default stage",
			stored: []byte("some-name;some-image"),
			expected: &api.Asset{
				Name:  "some-name",
				Image: "some-image",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded api.Asset
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							assetKey(test.expected.Name),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						helper := NewTxnHelper(txn)
						return helper.LoadAsset(
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
