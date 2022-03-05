package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStoreAssetWithTxn(t *testing.T) {
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

func createAsset(
	t *testing.T,
	assetName string,
	requiredOnly bool,
) domain.Asset {
	name, err := domain.NewAssetName(assetName)
	assert.NilError(t, err, "create name for asset %s", assetName)

	if !requiredOnly {
		image, err := domain.NewImage("some-image")
		assert.NilError(t, err, "create image for asset %s", assetName)
		return domain.NewAssetWithImage(name, image)
	}
	return domain.NewAssetWithoutImage(name)
}
