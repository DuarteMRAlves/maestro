package asset

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"sync"
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

		t.Run(testName, func(t *testing.T) {
			st := &store{assets: map[identifier.Id]*Asset{}, lock: sync.RWMutex{}}

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
	desc *Asset) {

	assert.DeepEqual(t, IdSize, assetId.Size(), "asset identifier size")
	assert.DeepEqual(t, 1, len(st.assets), "store size")
	asset, ok := st.assets[assetId]
	assert.IsTrue(t, ok, "asset exists")
	assert.DeepEqual(t, desc.Name, asset.Name, "correct names")
	assert.DeepEqual(t, assetId, asset.Id, "correct identifier")
}

func assertIncorrectCreate(
	t *testing.T,
	st *store,
	assetId identifier.Id) {
	assert.DeepEqual(t, 0, len(st.assets), "store size")
	assert.DeepEqual(t, identifier.Empty(), assetId, "correct identifier")
}
