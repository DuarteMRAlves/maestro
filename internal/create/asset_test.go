package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name      string
		assetName internal.AssetName
		image     internal.OptionalImage
		expected  internal.Asset
	}{
		{
			name:      "required fields",
			assetName: createAssetName(t, "some-name"),
			expected:  createAsset(t, "some-name", true),
		},
		{
			name:      "all fields",
			assetName: createAssetName(t, "some-name"),
			image:     internal.NewPresentImage(internal.NewImage("some-image")),
			expected:  createAsset(t, "some-name", false),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				storage := mock.AssetStorage{
					Assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				err := createFn(test.assetName, test.image)
				assert.NilError(t, err)

				assert.Equal(t, 1, len(storage.Assets))

				asset, exists := storage.Assets[test.expected.Name()]
				assert.Assert(t, exists)
				assertEqualAsset(t, test.expected, asset)
			},
		)
	}
}

func TestCreateAsset_Err(t *testing.T) {
	tests := []struct {
		name      string
		assetName internal.AssetName
		image     internal.OptionalImage
		isError   error
	}{
		{
			name:      "empty name",
			assetName: createAssetName(t, ""),
			isError:   EmptyAssetName,
		},
		{
			name:      "empty image",
			assetName: createAssetName(t, "some-name"),
			image:     internal.NewPresentImage(internal.NewImage("")),
			isError:   EmptyImageName,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				storage := mock.AssetStorage{
					Assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				err := createFn(test.assetName, test.image)
				assert.Assert(t, err != nil)
				assert.Assert(t, errors.Is(err, test.isError))

				assert.Equal(t, 0, len(storage.Assets))
			},
		)
	}
}

func TestCreateAsset_AlreadyExists(t *testing.T) {
	name := "some-name"
	assetName := createAssetName(t, name)
	image1 := internal.NewEmptyImage()
	image2 := internal.NewPresentImage(internal.NewImage("some-image"))
	expected := createAsset(t, name, true)
	storage := mock.AssetStorage{
		Assets: map[internal.AssetName]internal.Asset{},
	}

	createFn := Asset(storage)

	err := createFn(assetName, image1)
	assert.NilError(t, err, "first create")
	assert.Equal(t, 1, len(storage.Assets))
	asset, exists := storage.Assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)

	err = createFn(assetName, image2)
	assert.Assert(t, err != nil)
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "asset", alreadyExists.Type)
	assert.Equal(t, name, alreadyExists.Ident)
	assert.Equal(t, 1, len(storage.Assets))
	asset, exists = storage.Assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)
}

func createAssetName(t *testing.T, assetName string) internal.AssetName {
	name, err := internal.NewAssetName(assetName)
	assert.NilError(t, err, "create asset name %s", assetName)
	return name
}

func createAsset(
	t *testing.T,
	assetName string,
	requiredOnly bool,
) internal.Asset {
	name, err := internal.NewAssetName(assetName)
	assert.NilError(t, err, "create name for asset %s", assetName)
	imgOpt := internal.NewEmptyImage()
	if !requiredOnly {
		img := internal.NewImage("some-image")
		imgOpt = internal.NewPresentImage(img)
	}
	return internal.NewAsset(name, imgOpt)
}

func assertEqualAsset(t *testing.T, expected, actual internal.Asset) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assert.Equal(t, expected.Image().Present(), actual.Image().Present())
	if expected.Image().Present() {
		assert.Equal(
			t,
			expected.Image().Unwrap().Unwrap(),
			actual.Image().Unwrap().Unwrap(),
		)
	}
}
