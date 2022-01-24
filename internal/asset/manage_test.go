package asset

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_Create(t *testing.T) {
	const (
		assetName  = "Asset-Name"
		assetImage = "Asset-Image"
	)
	var (
		asset Asset
		err   error
	)
	cfg := &apitypes.Asset{Name: assetName, Image: assetImage}

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return Create(txn, cfg)
	})
	assert.NilError(t, err, "create error not nil")
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(assetKey(assetName))
		assert.NilError(t, err, "get error")
		cp, err := item.ValueCopy(nil)
		return load(&asset, cp)
	})
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name(), cfg.Name, "name not correct")
	assert.Equal(t, asset.Image(), cfg.Image, "image not correct")
}

func TestStore_CreateInvalidArguments(t *testing.T) {
	const assetName = "Asset-Name"
	tests := []struct {
		name   string
		cfg    *apitypes.Asset
		errMsg string
	}{
		{
			name:   "nil cfg",
			cfg:    nil,
			errMsg: "'cfg' is nil",
		},
	}

	for _, test := range tests {
		cfg, errMsg := test.cfg, test.errMsg

		t.Run(
			test.name, func(t *testing.T) {
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				err = db.Update(func(txn *badger.Txn) error {
					return Create(txn, cfg)
				})
				assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
				assert.ErrorContains(t, err, errMsg)
				err = db.View(func(txn *badger.Txn) error {
					item, err := txn.Get(assetKey(assetName))
					assert.Assert(t, item == nil, "nil item")
					return err
				})
				assert.Assert(t, err == badger.ErrKeyNotFound)
			})
	}
}

func TestStore_CreateAlreadyExists(t *testing.T) {
	const (
		assetName  = "Asset-Name"
		assetImage = "Asset-Image"
	)
	var (
		asset Asset
		err   error
	)
	cfg := &apitypes.Asset{Name: assetName, Image: assetImage}

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	// First create
	err = db.Update(func(txn *badger.Txn) error {
		return Create(txn, cfg)
	})
	assert.NilError(t, err, "create error not nil")
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(assetKey(assetName))
		assert.NilError(t, err, "get error")
		cp, err := item.ValueCopy(nil)
		return load(&asset, cp)
	})
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name(), cfg.Name, "name not correct")
	assert.Equal(t, asset.Image(), cfg.Image, "image not correct")

	// Create with new image
	cfg.Image = fmt.Sprintf("%v-new", assetImage)
	err = db.Update(func(txn *badger.Txn) error {
		return Create(txn, cfg)
	})
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("asset '%v' already exists", cfg.Name))

	// Store should keep old asset
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(assetKey(assetName))
		assert.NilError(t, err, "get error")
		cp, err := item.ValueCopy(nil)
		return load(&asset, cp)
	})
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name(), cfg.Name, "name not correct")
	// Still should be old image as asset is not replaced
	assert.Equal(t, asset.Image(), assetImage, "image not correct")
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name  string
		query *apitypes.Asset
		// numbers to store
		stored []int
		// names of the expected assets
		expected []apitypes.AssetName
	}{
		{
			name:     "zero elements store, nil query",
			query:    nil,
			stored:   []int{},
			expected: []apitypes.AssetName{},
		},
		{
			name:     "zero elements store, some query",
			query:    &apitypes.Asset{Name: "some-name"},
			stored:   []int{},
			expected: []apitypes.AssetName{},
		},
		{
			name:     "one element stored, nil query",
			query:    nil,
			stored:   []int{0},
			expected: []apitypes.AssetName{testutil.AssetNameForNum(0)},
		},
		{
			name:   "multiple elements stored, nil query",
			query:  nil,
			stored: []int{0, 1, 2},
			expected: []apitypes.AssetName{
				testutil.AssetNameForNum(0),
				testutil.AssetNameForNum(1),
				testutil.AssetNameForNum(2),
			},
		},
		{
			name:     "multiple elements stored, matching name query",
			query:    &apitypes.Asset{Name: testutil.AssetNameForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []apitypes.AssetName{testutil.AssetNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching name query",
			query:    &apitypes.Asset{Name: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []apitypes.AssetName{},
		},
		{
			name:     "multiple elements stored, matching image query",
			query:    &apitypes.Asset{Image: testutil.AssetImageForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []apitypes.AssetName{testutil.AssetNameForNum(2)},
		},
		{
			query:    &apitypes.Asset{Image: "unknown-image"},
			stored:   []int{0, 1, 2},
			expected: []apitypes.AssetName{},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var received []*apitypes.Asset
				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				for _, n := range test.stored {
					err = db.Update(func(txn *badger.Txn) error {
						return Create(txn, assetForNum(n))
					})
					assert.NilError(t, err, "create asset error")
				}

				err = db.View(func(txn *badger.Txn) error {
					received, err = Get(txn, test.query)
					return err
				})
				assert.NilError(t, err, "get assets")
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[apitypes.AssetName]bool, 0)
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

func assetForNum(num int) *apitypes.Asset {
	return &apitypes.Asset{
		Name:  testutil.AssetNameForNum(num),
		Image: testutil.AssetImageForNum(num),
	}
}
