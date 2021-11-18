package asset

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

func TestStore_Create(t *testing.T) {
	tests := []struct {
		desc *Asset
		err  error
	}{
		{&Asset{Name: assetName}, nil},
	}

	for _, inner := range tests {
		desc, err := inner.desc, inner.err
		testName := fmt.Sprintf("desc='%v',err='%v'", desc, err)

		t.Run(
			testName, func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.IsTrue(t, ok, "type assertion failed for store")

				e := st.Create(desc)
				assert.DeepEqual(t, err, e, "expected error")
				if err == nil {
					assertCorrectCreate(t, st, desc)
				} else {
					assertIncorrectCreate(t, st)
				}
			})
	}
}

func assertCorrectCreate(
	t *testing.T,
	st *store,
	desc *Asset,
) {

	assert.DeepEqual(t, 1, lenAssets(st), "store size")
	stored, ok := st.assets.Load(assetName)
	assert.IsTrue(t, ok, "asset exists")
	asset, ok := stored.(*Asset)
	assert.IsTrue(t, ok, "asset type assertion failed")
	assert.DeepEqual(t, desc.Name, asset.Name, "correct names")
}

func assertIncorrectCreate(
	t *testing.T,
	st *store,
) {
	assert.DeepEqual(t, 0, lenAssets(st), "store size")
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
