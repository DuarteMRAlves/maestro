package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name     string
		req      AssetRequest
		expected internal.Asset
	}{
		{
			name: "required fields",
			req: AssetRequest{
				Name:  "some-name",
				Image: domain.NewEmptyString(),
			},
			expected: createAsset(t, "some-name", true),
		},
		{
			name: "all fields",
			req: AssetRequest{
				Name:  "some-name",
				Image: domain.NewPresentString("some-image"),
			},
			expected: createAsset(t, "some-name", false),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				storage := mockAssetStorage{
					assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				res := createFn(test.req)
				assert.Assert(t, !res.Err.Present())

				assert.Equal(t, 1, len(storage.assets))

				asset, exists := storage.assets[test.expected.Name()]
				assert.Assert(t, exists)
				assertEqualAsset(t, test.expected, asset)
			},
		)
	}
}

func TestCreateAsset_Err(t *testing.T) {
	tests := []struct {
		name    string
		req     AssetRequest
		isError error
	}{
		{
			name:    "empty name",
			req:     AssetRequest{Name: ""},
			isError: EmptyAssetName,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				storage := mockAssetStorage{
					assets: map[internal.AssetName]internal.Asset{},
				}

				createFn := Asset(storage)
				res := createFn(test.req)
				assert.Assert(t, res.Err.Present())

				assert.Equal(t, 0, len(storage.assets))

				err := res.Err.Unwrap()
				assert.Assert(t, errors.Is(err, test.isError))
			},
		)
	}
}

func TestCreateAsset_AlreadyExists(t *testing.T) {
	req := AssetRequest{
		Name:  "some-name",
		Image: domain.NewEmptyString(),
	}
	expected := createAsset(t, "some-name", true)
	storage := mockAssetStorage{
		assets: map[internal.AssetName]internal.Asset{},
	}

	createFn := Asset(storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, 1, len(storage.assets))
	asset, exists := storage.assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "asset", alreadyExists.Type)
	assert.Equal(t, req.Name, alreadyExists.Ident)
	assert.Equal(t, 1, len(storage.assets))
	asset, exists = storage.assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)
}

type mockAssetStorage struct {
	assets map[internal.AssetName]internal.Asset
}

func (m mockAssetStorage) Save(asset internal.Asset) error {
	m.assets[asset.Name()] = asset
	return nil
}

func (m mockAssetStorage) Load(name internal.AssetName) (
	internal.Asset,
	error,
) {
	asset, exists := m.assets[name]
	if !exists {
		err := &internal.NotFound{Type: "asset", Ident: name.Unwrap()}
		return internal.Asset{}, err
	}
	return asset, nil
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
