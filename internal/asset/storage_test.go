package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestTxnHelper_SaveAsset(t *testing.T) {
	tests := []struct {
		name     string
		asset    domain.Asset
		expected []byte
	}{
		{
			name: "default asset",
			asset: &asset{
				name:  assetName("some-name"),
				image: image(""),
			},
			expected: []byte("some-name;"),
		},
		{
			name: "non default stage",
			asset: &asset{
				name:  assetName("some-name"),
				image: image("some-image"),
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
						item, err := txn.Get(kv.AssetKey(test.asset.Name()))
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
		expected domain.Asset
		stored   []byte
	}{
		{
			name:   "default asset",
			stored: []byte("some-name;"),
			expected: &asset{
				name:  assetName("some-name"),
				image: image(""),
			},
		},
		{
			name:   "non default stage",
			stored: []byte("some-name;some-image"),
			expected: &asset{
				name:  assetName("some-name"),
				image: image("some-image"),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var loaded domain.Asset
				db := kv.NewTestDb(t)
				defer db.Close()

				err := db.Update(
					func(txn *badger.Txn) error {
						return txn.Set(
							kv.AssetKey(test.expected.Name()),
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
					test.expected.Image().Unwrap(),
					loaded.Image().Unwrap(),
				)
			},
		)
	}
}
