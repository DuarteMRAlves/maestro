package asset

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_CreateCorrect(t *testing.T) {
	const (
		assetName  = "Asset-Name"
		assetImage = "Asset-Image"
	)
	tests := []struct {
		name   string
		config *Asset
	}{
		{
			name:   "non default params",
			config: &Asset{Name: assetName, Image: assetImage},
		},
		{
			name:   "default params",
			config: &Asset{Name: "", Image: ""},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				config := test.config

				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				e := st.Create(config)
				assert.NilError(t, e, "error not nil")
				assert.Equal(t, 1, lenAssets(st), "store size")
				stored, ok := st.assets.Load(config.Name)
				assert.Assert(t, ok, "config does not exist")
				asset, ok := stored.(*Asset)
				assert.Assert(t, ok, "config type assertion failed")
				assert.Equal(t, asset.Name, config.Name, "name not correct")
				assert.Equal(t, asset.Image, config.Image, "image not correct")
			})
	}
}

func TestStore_CreateInvalidArguments(t *testing.T) {
	const assetName = "Asset-Name"
	tests := []struct {
		name   string
		config *Asset
		errMsg string
	}{
		{
			name:   "nil config",
			config: nil,
			errMsg: "'config' is nil",
		},
	}

	for _, test := range tests {
		config, errMsg := test.config, test.errMsg

		t.Run(
			test.name, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(config)
				assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
				assert.ErrorContains(t, err, errMsg)
				assert.Equal(t, 0, lenAssets(st), "store size")
				_, ok = st.assets.Load(assetName)
				assert.Assert(t, !ok, "config does not exist")
			})
	}
}

func TestStore_CreateAlreadyExists(t *testing.T) {
	const (
		assetName  apitypes.AssetName = "Asset-Name"
		assetImage                    = "Asset-Image"
	)
	config := &Asset{
		Name:  assetName,
		Image: assetImage,
	}
	st, ok := NewStore().(*store)
	assert.Assert(t, ok, "type assertion failed for store")

	// First create should go well
	err := st.Create(config)
	assert.NilError(t, err, "error not nil")
	assert.Equal(t, 1, lenAssets(st), "store size")
	stored, ok := st.assets.Load(assetName)
	assert.Assert(t, ok, "asset does not exist")
	asset, ok := stored.(*Asset)
	assert.Assert(t, ok, "asset type assertion failed")
	assert.Equal(t, assetName, asset.Name, "name not correct")
	assert.Equal(t, assetImage, asset.Image, "image not correct")

	// Create new image
	config.Image = fmt.Sprintf("%v-new", assetImage)
	err = st.Create(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("asset '%v' already exists", config.Name))
	// Store should keep old asset
	assert.Equal(t, 1, lenAssets(st), "store size")
	stored, ok = st.assets.Load(assetName)
	assert.Assert(t, ok, "asset does not exist")
	asset, ok = stored.(*Asset)
	assert.Assert(t, ok, "asset type assertion failed")
	assert.Equal(t, assetName, asset.Name, "name not correct")
	// Still should be old image as asset is not replaced
	assert.Equal(t, assetImage, asset.Image, "image not correct")
}

func lenAssets(st *store) int {
	count := 0
	st.assets.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
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
			name:     "multiple elements stored, non-matching image query",
			query:    &apitypes.Asset{Image: "unknown-image"},
			stored:   []int{0, 1, 2},
			expected: []apitypes.AssetName{},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				st := NewStore()

				for _, n := range test.stored {
					err := st.Create(assetForNum(n))
					assert.NilError(t, err, "create asset error")
				}

				received := st.Get(test.query)
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

func assetForNum(num int) *Asset {
	return &Asset{
		Name:  testutil.AssetNameForNum(num),
		Image: testutil.AssetImageForNum(num),
	}
}
