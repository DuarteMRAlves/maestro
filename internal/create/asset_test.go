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
				existsCallCount := 0
				saveCallCount := 0
				existsFn := existsAssetFn(
					test.expected.Name(),
					&existsCallCount,
				)
				saveFn := saveAssetFn(t, test.expected, &saveCallCount)
				createFn := CreateAsset(existsFn, saveFn)
				res := createFn(test.req)
				assert.Assert(t, !res.Err.Present())
				assert.Equal(t, existsCallCount, 1)
				assert.Equal(t, saveCallCount, 1)
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
	existsCallCount := 0
	saveCallCount := 0
	existsFn := existsAssetFn(expected.Name(), &existsCallCount)
	saveFn := saveAssetFn(t, expected, &saveCallCount)
	createFn := CreateAsset(existsFn, saveFn)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsCallCount, 1)
	assert.Equal(t, saveCallCount, 1)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("asset '%v' already exists", req.Name),
	)
	assert.Equal(t, existsCallCount, 2)
	// Should not call save
	assert.Equal(t, saveCallCount, 1)
}

func existsAssetFn(expected domain.AssetName, callCount *int) ExistsAsset {
	return func(name domain.AssetName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > 1)
	}
}

func saveAssetFn(
	t *testing.T,
	expected domain.Asset,
	callCount *int,
) SaveAsset {
	return func(actual domain.Asset) domain.AssetResult {
		*callCount++
		assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
		assert.Equal(t, expected.Image().Present(), actual.Image().Present())
		if expected.Image().Present() {
			assert.Equal(t, expected.Image().Unwrap(), actual.Image().Unwrap())
		}
		return domain.SomeAsset(actual)
	}
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
