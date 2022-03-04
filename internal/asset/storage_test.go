package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStoreAssetWithTxn(t *testing.T) {
	tests := []struct {
		name     string
		asset    types.Asset
		expected []byte
	}{
		{
			name: "default asset",
			asset: &asset{
				name:  assetName("some-name"),
				image: emptyImage{},
			},
			expected: []byte("some-name;"),
		},
		{
			name: "non default stage",
			asset: &asset{
				name:  assetName("some-name"),
				image: presentImage{image("some-image")},
			},
			expected: []byte("some-name;some-image"),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var stored []byte

				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						store := StoreAssetWithTxn(txn)
						result := store(test.asset)
						return result.Error()
					},
				)
				assert.NilError(t, err, "save error")

				err = db.View(
					func(txn *badger.Txn) error {
						item, err := txn.Get(kvKey(test.asset.Name()))
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
		expected types.Asset
		stored   []byte
	}{
		{
			name:   "default asset",
			stored: []byte("some-name;"),
			expected: &asset{
				name:  assetName("some-name"),
				image: emptyImage{},
			},
		},
		{
			name:   "non default stage",
			stored: []byte("some-name;some-image"),
			expected: &asset{
				name:  assetName("some-name"),
				image: presentImage{image("some-image")},
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded types.Asset
				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							kvKey(test.expected.Name()),
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
