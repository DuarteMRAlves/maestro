package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name     string
		req      AssetRequest
		expected domain.Asset
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
					assets: map[domain.AssetName]domain.Asset{},
				}

				createFn := CreateAsset(storage)
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

func TestCreateAsset_AlreadyExists(t *testing.T) {
	req := AssetRequest{
		Name:  "some-name",
		Image: domain.NewEmptyString(),
	}
	expected := createAsset(t, "some-name", true)
	storage := mockAssetStorage{
		assets: map[domain.AssetName]domain.Asset{},
	}

	createFn := CreateAsset(storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, 1, len(storage.assets))
	asset, exists := storage.assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("asset '%v' already exists", req.Name),
	)
	assert.Equal(t, 1, len(storage.assets))
	asset, exists = storage.assets[expected.Name()]
	assert.Assert(t, exists)
	assertEqualAsset(t, expected, asset)
}

type mockAssetStorage struct {
	assets map[domain.AssetName]domain.Asset
}

func (m mockAssetStorage) Save(asset domain.Asset) domain.AssetResult {
	m.assets[asset.Name()] = asset
	return domain.SomeAsset(asset)
}

func (m mockAssetStorage) Load(name domain.AssetName) domain.AssetResult {
	asset, exists := m.assets[name]
	if !exists {
		err := errdefs.NotFoundWithMsg("asset not found: %s", name)
		return domain.ErrAsset(err)
	}
	return domain.SomeAsset(asset)
}

func createAsset(
	t *testing.T,
	assetName string,
	requiredOnly bool,
) domain.Asset {
	name, err := domain.NewAssetName(assetName)
	assert.NilError(t, err, "create name for asset %s", assetName)
	imgOpt := domain.NewEmptyImage()
	if !requiredOnly {
		img, err := domain.NewImage("some-image")
		assert.NilError(t, err, "create image for asset %s", assetName)
		imgOpt = domain.NewPresentImage(img)
	}
	return domain.NewAsset(name, imgOpt)
}

func assertEqualAsset(t *testing.T, expected, actual domain.Asset) {
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
