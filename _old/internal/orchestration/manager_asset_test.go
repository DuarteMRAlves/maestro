package orchestration

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

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
			expected: []api.AssetName{api.AssetName("asset-0")},
		},
		{
			name:   "multiple elements stored, nil req",
			req:    nil,
			stored: []int{0, 1, 2},
			expected: []api.AssetName{
				api.AssetName("asset-0"),
				api.AssetName("asset-1"),
				api.AssetName("asset-2"),
			},
		},
		{
			name:     "multiple elements stored, matching name req",
			req:      &api.GetAssetRequest{Name: api.AssetName("asset-2")},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{api.AssetName("asset-2")},
		},
		{
			name:     "multiple elements stored, non-matching name req",
			req:      &api.GetAssetRequest{Name: "unknown-name"},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{},
		},
		{
			name:     "multiple elements stored, matching image req",
			req:      &api.GetAssetRequest{Image: "image-2"},
			stored:   []int{0, 1, 2},
			expected: []api.AssetName{api.AssetName("asset-2")},
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
				var (
					received []*api.Asset
					err      error
				)

				db := storage.NewTestDb(t)
				defer db.Close()

				for _, n := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							createAsset := CreateAssetWithTxn(txn)
							return createAsset(assetForNum(n))
						},
					)
					assert.NilError(t, err, "create asset error")
				}

				err = db.View(
					func(txn *badger.Txn) error {
						getAssets := GetAssetsWithTxn(txn)
						received, err = getAssets(test.req)
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
		Name:  api.AssetName(fmt.Sprintf("asset-%d", num)),
		Image: fmt.Sprintf("image-%d", num),
	}
}