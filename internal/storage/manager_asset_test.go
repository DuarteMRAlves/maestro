package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_CreateAsset(t *testing.T) {
	const (
		assetName  = "Asset-Name"
		assetImage = "Asset-Image"
	)
	var (
		asset api.Asset
		err   error
	)
	cfg := &api.CreateAssetRequest{Name: assetName, Image: assetImage}

	m := NewManager(rpc.NewManager())

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateAsset(txn, cfg)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get(assetKey(assetName))
			assert.NilError(t, err, "get error")
			cp, err := item.ValueCopy(nil)
			return loadAsset(&asset, cp)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name, cfg.Name, "name not correct")
	assert.Equal(t, asset.Image, cfg.Image, "image not correct")
}

func TestManager_CreateAsset_InvalidArguments(t *testing.T) {
	const assetName = "Asset-Name"

	var (
		req    *api.CreateAssetRequest = nil
		errMsg                         = "'req' is nil"
	)

	m := NewManager(rpc.NewManager())

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateAsset(txn, req)
		},
	)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, errMsg)
	err = db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get(assetKey(assetName))
			assert.Assert(t, item == nil, "nil item")
			return err
		},
	)
	assert.Assert(t, err == badger.ErrKeyNotFound)
}

func TestManager_CreateAsset_AlreadyExists(t *testing.T) {
	const (
		assetName  = "Asset-Name"
		assetImage = "Asset-Image"
	)
	var (
		asset api.Asset
		err   error
	)
	req := &api.CreateAssetRequest{Name: assetName, Image: assetImage}

	m := NewManager(rpc.NewManager())

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	// First create
	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateAsset(txn, req)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get(assetKey(assetName))
			assert.NilError(t, err, "get error")
			cp, err := item.ValueCopy(nil)
			return loadAsset(&asset, cp)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name, req.Name, "name not correct")
	assert.Equal(t, asset.Image, req.Image, "image not correct")

	// Create with new image
	req.Image = fmt.Sprintf("%v-new", assetImage)
	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateAsset(txn, req)
		},
	)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("asset '%v' already exists", req.Name),
	)

	// Store should keep old asset
	err = db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get(assetKey(assetName))
			assert.NilError(t, err, "get error")
			cp, err := item.ValueCopy(nil)
			return loadAsset(&asset, cp)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, asset.Name, req.Name, "name not correct")
	// Still should be old image as asset is not replaced
	assert.Equal(t, asset.Image, assetImage, "image not correct")
}

func TestManager_GetMatchingAssets(t *testing.T) {
	tests := []struct {
		name string
		req  *api.GetAssetRequest
		// numbers to store
		stored []int
		// names of the expected assets
		expected []api.AssetName
	}{
		{
			name:     "zero elements store, nil req",
			req:      nil,
			stored:   []int{},
			expected: []api.AssetName{},
		},
		{
			name:     "zero elements store, some req",
			req:      &api.GetAssetRequest{Name: "some-name"},
			stored:   []int{},
			expected: []api.AssetName{},
		},
		{
			name:     "one element stored, nil req",
			req:      nil,
			stored:   []int{0},
			expected: []api.AssetName{util.AssetNameForNum(0)},
		},
		{
			name:   "multiple elements stored, nil req",
			req:    nil,
			stored: []int{0, 1, 2},
			expected: []api.AssetName{
				util.AssetNameForNum(0),
				util.AssetNameForNum(1),
				util.AssetNameForNum(2),
			},
		},
		{
			name:     "multiple elements stored, matching name req",
			req:      &api.GetAssetRequest{Name: util.AssetNameForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{util.AssetNameForNum(2)},
		},
		{
			name:     "multiple elements stored, non-matching name req",
			req:      &api.GetAssetRequest{Name: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{},
		},
		{
			name:     "multiple elements stored, matching image req",
			req:      &api.GetAssetRequest{Image: util.AssetImageForNum(2)},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{util.AssetNameForNum(2)},
		},
		{
			req:      &api.GetAssetRequest{Image: "unknown-image"},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var received []*api.Asset

				m := NewManager(rpc.NewManager())

				db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()

				for _, n := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							return m.CreateAsset(txn, assetForNum(n))
						},
					)
					assert.NilError(t, err, "create asset error")
				}

				err = db.View(
					func(txn *badger.Txn) error {
						received, err = m.GetMatchingAssets(txn, test.req)
						return err
					},
				)
				assert.NilError(t, err, "get assets")
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[api.AssetName]bool, 0)
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

func assetForNum(num int) *api.CreateAssetRequest {
	return &api.CreateAssetRequest{
		Name:  util.AssetNameForNum(num),
		Image: util.AssetImageForNum(num),
	}
}
