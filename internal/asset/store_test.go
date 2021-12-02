package asset

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStore_CreateCorrect(t *testing.T) {
	tests := []*Asset{
		{Name: assetName},
		{Name: assetName, Image: assetImage},
	}
	for _, a := range tests {
		testName := fmt.Sprintf("a='%v'", a)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				e := st.Create(a)
				assert.NilError(t, e, "error not nil")
				assert.Equal(t, 1, lenAssets(st), "store size")
				stored, ok := st.assets.Load(assetName)
				assert.Assert(t, ok, "asset does not exist")
				asset, ok := stored.(*Asset)
				assert.Assert(t, ok, "asset type assertion failed")
				assert.Equal(t, asset.Name, a.Name, "name not correct")
				assert.Equal(t, asset.Image, a.Image, "image not correct")
			})
	}
}

func TestStore_CreateInvalidArguments(t *testing.T) {
	tests := []struct {
		a      *Asset
		errMsg string
	}{
		{
			nil,
			"'config' is nil",
		},
		{
			&Asset{
				// No name will create empty string
				Image: assetImage,
			},
			"invalid name ''",
		},
		{
			&Asset{
				Name:  "",
				Image: assetImage,
			},
			"invalid name ''",
		},
		{
			&Asset{
				Name:  "invalid-name/",
				Image: assetImage,
			},
			"invalid name 'invalid-name/'",
		},
	}

	for _, inner := range tests {
		a, errMsg := inner.a, inner.errMsg
		testName := fmt.Sprintf("config=%v, errMsg=%v", a, errMsg)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(a)
				assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
				assert.ErrorContains(t, err, errMsg)
				assert.Equal(t, 0, lenAssets(st), "store size")
				_, ok = st.assets.Load(assetName)
				assert.Assert(t, !ok, "asset does not exist")
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
