package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/create"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestSaveOrchestrationWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		orch     create.Orchestration
		expected []byte
	}{
		{
			name: "empty fields",
			orch: createCreateOrchestration(
				t,
				"some-name",
				[]string{},
				[]string{},
			),
			expected: []byte(";"),
		},
		{
			name: "all fields",
			orch: createCreateOrchestration(
				t,
				"some-name",
				[]string{"stage-1", "stage-2", "stage-3"},
				[]string{"link-1", "link-2"},
			),
			expected: []byte("stage-1,stage-2,stage-3;link-1,link-2"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var stored []byte

				db := NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						saveFn := SaveOrchestrationWithTxn(txn)
						result := saveFn(test.orch)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(orchestrationKey(test.orch.Name()))
						if err != nil {
							return err
						}
						stored, err = item.ValueCopy(nil)
						if err != nil {
							return err
						}
						return nil
					},
				)
				assert.Equal(t, len(test.expected), len(stored))
				for i, e := range test.expected {
					assert.Equal(t, e, stored[i])
				}
			},
		)
	}
}

func TestLoadOrchestrationWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		expected execute.Orchestration
		stored   []byte
	}{
		{
			name: "required fields",
			expected: createExecuteOrchestration(
				t,
				"some-name",
				[]string{"stage-1", "stage-2", "stage-3"},
				[]string{"link-1", "link-2"},
				true,
			),
			stored: []byte("stage-1,stage-2,stage-3;link-1,link-2"),
		},
		{
			name: "all fields",
			expected: createExecuteOrchestration(
				t,
				"some-name",
				[]string{"stage-1", "stage-2", "stage-3"},
				[]string{"link-1", "link-2"},
				false,
			),
			stored: []byte("stage-1,stage-2,stage-3;link-1,link-2"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded execute.Orchestration

				db := NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						saveStageFn := SaveStageWithTxn(txn)
						saveLinkFn := SaveLinkWithTxn(txn)

						for _, s := range test.expected.Stages() {
							res := saveStageFn(s)
							if res.IsError() {
								return res.Error()
							}
						}
						for _, l := range test.expected.Links() {
							res := saveLinkFn(l)
							if res.IsError() {
								return res.Error()
							}
						}
						return txn.Set(
							orchestrationKey(test.expected.Name()),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						loadFn := LoadOrchestrationWithTxn(txn)
						res := loadFn(test.expected.Name())
						if !res.IsError() {
							loaded = res.Unwrap()
						}
						return res.Error()
					},
				)
				assert.NilError(t, err, "load error")
				assert.Equal(
					t,
					test.expected.Name().Unwrap(),
					loaded.Name().Unwrap(),
				)

				loadedStages := loaded.Stages()
				loadedLinks := loaded.Links()

				assert.Equal(t, len(test.expected.Stages()), len(loadedStages))
				assert.Equal(t, len(test.expected.Links()), len(loadedLinks))

				for i, s := range test.expected.Stages() {
					assertEqualStage(t, s, loadedStages[i])
				}

				for i, l := range test.expected.Links() {
					assertEqualLink(t, l, loadedLinks[i])
				}
			},
		)
	}
}
