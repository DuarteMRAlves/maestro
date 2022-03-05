package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestSaveLinkWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		link     domain.Link
		expected []byte
	}{
		{
			name:     "required fields",
			link:     createLink(t, "some-name", "source", "target", true),
			expected: []byte("source;;target;"),
		},
		{
			name:     "all fields",
			link:     createLink(t, "some-name", "source", "target", false),
			expected: []byte("source;source-field;target;target-field"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					storedLink   []byte
					loadedSource domain.Stage
					loadedTarget domain.Stage
				)

				db := NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						store := SaveWithTxn(txn)
						result := store(test.link)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						linkItem, err := txn.Get(linkKey(test.link.Name()))
						if err != nil {
							return err
						}
						storedLink, err = linkItem.ValueCopy(nil)
						if err != nil {
							return err
						}
						loadStage := LoadStageWithTxn(txn)
						sourceRes := loadStage(test.link.Source().Stage().Name())
						if sourceRes.IsError() {
							return sourceRes.Error()
						}
						loadedSource = sourceRes.Unwrap()
						targetRes := loadStage(test.link.Target().Stage().Name())
						if targetRes.IsError() {
							return targetRes.Error()
						}
						loadedTarget = targetRes.Unwrap()
						return nil
					},
				)
				assert.Equal(t, len(test.expected), len(storedLink))
				for i, e := range test.expected {
					assert.Equal(t, e, storedLink[i])
				}
				assertEqualStage(t, test.link.Source().Stage(), loadedSource)
				assertEqualStage(t, test.link.Target().Stage(), loadedTarget)
			},
		)
	}
}

func TestLoadLinkWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		expected domain.Link
		stored   []byte
	}{
		{
			name:     "required fields",
			expected: createLink(t, "some-name", "source", "target", true),
			stored:   []byte("source;;target;"),
		},
		{
			name:     "all fields",
			expected: createLink(t, "some-name", "source", "target", false),
			stored:   []byte("source;source-field;target;target-field"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded domain.Link

				db := NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						storeStage := SaveStageWithTxn(txn)
						res := storeStage(test.expected.Source().Stage())
						if res.IsError() {
							return res.Error()
						}
						res = storeStage(test.expected.Target().Stage())
						if res.IsError() {
							return res.Error()
						}
						return txn.Set(
							linkKey(test.expected.Name()),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")
				err = db.View(
					func(txn *badger.Txn) error {
						load := LoadLinkWithTxn(txn)
						res := load(test.expected.Name())
						if !res.IsError() {
							loaded = res.Unwrap()
						}
						return res.Error()
					},
				)
				assert.NilError(t, err, "load error")
				fmt.Println(loaded)
				assert.Equal(
					t,
					test.expected.Name().Unwrap(),
					loaded.Name().Unwrap(),
				)
				assertEqualStage(
					t,
					test.expected.Source().Stage(),
					loaded.Source().Stage(),
				)
				assertEqualStage(
					t,
					test.expected.Target().Stage(),
					loaded.Target().Stage(),
				)
				if test.expected.Source().Field().Present() {
					assert.Equal(
						t,
						test.expected.Source().Field().Unwrap(),
						loaded.Source().Field().Unwrap(),
					)
				} else {
					assert.Assert(t, !loaded.Source().Field().Present())
				}
				if test.expected.Target().Field().Present() {
					assert.Equal(
						t,
						test.expected.Target().Field().Unwrap(),
						loaded.Target().Field().Unwrap(),
					)
				} else {
					assert.Assert(t, !loaded.Target().Field().Present())
				}
			},
		)
	}
}
