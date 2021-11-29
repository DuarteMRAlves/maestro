package asset

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/test"
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
				test.IsTrue(t, ok, "type assertion failed for store")

				e := st.Create(a)
				test.DeepEqual(t, nil, e, "error not nil")
				test.DeepEqual(t, 1, lenAssets(st), "store size")
				stored, ok := st.assets.Load(assetName)
				test.IsTrue(t, ok, "asset does not exist")
				asset, ok := stored.(*Asset)
				test.IsTrue(t, ok, "asset type assertion failed")
				test.DeepEqual(t, asset.Name, a.Name, "name not correct")
				test.DeepEqual(t, asset.Image, a.Image, "image not correct")
			})
	}
}

func TestStore_CreateInvalidArguments(t *testing.T) {
	tests := []struct {
		a   *Asset
		err error
	}{
		{
			nil,
			errdefs.InvalidArgumentWithMsg("'config' is nil"),
		},
		{
			&Asset{
				// No name will create empty string
				Image: assetImage,
			},
			errdefs.InvalidArgumentWithMsg("invalid name ''"),
		},
		{
			&Asset{
				Name:  "",
				Image: assetImage,
			},
			errdefs.InvalidArgumentWithMsg("invalid name ''"),
		},
		{
			&Asset{
				Name:  "invalid-name/",
				Image: assetImage,
			},
			errdefs.InvalidArgumentWithMsg("invalid name 'invalid-name/'"),
		},
	}

	for _, inner := range tests {
		a, err := inner.a, inner.err
		testName := fmt.Sprintf("config=%v, err=%v", a, err)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				test.IsTrue(t, ok, "type assertion failed for store")

				e := st.Create(a)
				test.DeepEqual(t, err, e, "expected error")
				test.DeepEqual(t, 0, lenAssets(st), "store size")
				_, ok = st.assets.Load(assetName)
				test.IsTrue(t, !ok, "asset does not exist")
			})
	}
}

func TestStore_CreateAlreadyExists(t *testing.T) {
	config := &Asset{
		Name:  assetName,
		Image: assetImage,
	}
	st, ok := NewStore().(*store)
	test.IsTrue(t, ok, "type assertion failed for store")

	// First create should go well
	e := st.Create(config)
	test.DeepEqual(t, nil, e, "error not nil")
	test.DeepEqual(t, 1, lenAssets(st), "store size")
	stored, ok := st.assets.Load(assetName)
	test.IsTrue(t, ok, "asset does not exist")
	asset, ok := stored.(*Asset)
	test.IsTrue(t, ok, "asset type assertion failed")
	test.DeepEqual(t, assetName, asset.Name, "name not correct")
	test.DeepEqual(t, assetImage, asset.Image, "image not correct")

	// Create new image
	config.Image = fmt.Sprintf("%v-new", assetImage)
	e = st.Create(config)
	err := errdefs.AlreadyExistsWithMsg(
		"asset '%v' already exists",
		config.Name)
	test.DeepEqual(t, err, e, "error no already exists")
	// Store should keep old asset
	test.DeepEqual(t, 1, lenAssets(st), "store size")
	stored, ok = st.assets.Load(assetName)
	test.IsTrue(t, ok, "asset does not exist")
	asset, ok = stored.(*Asset)
	test.IsTrue(t, ok, "asset type assertion failed")
	test.DeepEqual(t, assetName, asset.Name, "name not correct")
	// Still should be old image as asset is not replaced
	test.DeepEqual(t, assetImage, asset.Image, "image not correct")
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
