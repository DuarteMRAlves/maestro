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
