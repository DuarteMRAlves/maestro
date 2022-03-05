package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestSaveAssetWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		asset    domain.Asset
		expected []byte
	}{
		{
			name:     "default asset",
			asset:    createAsset(t, "some-name", true),
			expected: []byte("some-name;"),
		},
		{
			name:     "non default stage",
			asset:    createAsset(t, "some-name", false),
			expected: []byte("some-name;some-image"),
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
						store := SaveAssetWithTxn(txn)
						result := store(test.asset)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(assetKey(test.asset.Name()))
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

func TestLoadAssetWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		expected domain.Asset
		stored   []byte
	}{
		{
			name:     "default asset",
			stored:   []byte("some-name;"),
			expected: createAsset(t, "some-name", true),
		},
		{
			name:     "non default stage",
			stored:   []byte("some-name;some-image"),
			expected: createAsset(t, "some-name", false),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded domain.Asset
				db := NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							assetKey(test.expected.Name()),
							test.stored,
						)
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						load := LoadAssetWithTxn(txn)
						res := load(test.expected.Name())
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
				assert.Equal(
					t,
					test.expected.Image().Present(),
					loaded.Image().Present(),
				)
				if test.expected.Image().Present() {
					assert.Equal(
						t,
						test.expected.Image().Unwrap().Unwrap(),
						loaded.Image().Unwrap().Unwrap(),
					)
				}
			},
		)
	}
}
