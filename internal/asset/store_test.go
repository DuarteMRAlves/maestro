package asset

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestStore_Create(t *testing.T) {
	prevId, _ := identifier.Rand(IdSize)
	tests := []struct {
		desc *Asset
		err  error
	}{
		{&Asset{Name: "asset name"}, nil},
		{&Asset{Id: identifier.Empty(), Name: "asset name"}, nil},
		{
			&Asset{Id: prevId, Name: "asset name"},
			errors.New("asset identifier should not be defined"),
		},
	}

	for _, inner := range tests {
		desc, err := inner.desc, inner.err
		testName := fmt.Sprintf("desc='%v',err='%v'", desc, err)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.IsTrue(t, ok, "type assertion failed for store")

				assetId, e := st.Create(desc)
				assert.DeepEqual(t, err, e, "expected error")
				if err == nil {
					assertCorrectCreate(t, st, assetId, desc)
				} else {
					assertIncorrectCreate(t, st, assetId)
				}
			})
	}
}

func assertCorrectCreate(
	t *testing.T,
	st *store,
	assetId identifier.Id,
	desc *Asset,
) {

	assert.DeepEqual(t, IdSize, assetId.Size(), "asset identifier size")
	assert.DeepEqual(t, 1, lenAssets(st), "store size")
	stored, ok := st.assets.Load(assetId)
	assert.IsTrue(t, ok, "asset exists")
	asset, ok := stored.(*Asset)
	assert.IsTrue(t, ok, "asset type assertion failed")
	assert.DeepEqual(t, desc.Name, asset.Name, "correct names")
	assert.DeepEqual(t, assetId, asset.Id, "correct identifier")
}

func assertIncorrectCreate(
	t *testing.T,
	st *store,
	assetId identifier.Id,
) {
	assert.DeepEqual(t, 0, lenAssets(st), "store size")
	assert.DeepEqual(t, identifier.Empty(), assetId, "correct identifier")
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
