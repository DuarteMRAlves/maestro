package asset

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_CreateCorrect(t *testing.T) {
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
