package identifier

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	idSize := 10
	gen := GenForSize(idSize)
	numIdsToGen := 10
	ids := make([]Id, 0, numIdsToGen)
	for i := 0; i < numIdsToGen; i++ {
		id, err := gen()
		assert.IsNil(t, err, "creating id")
		for _, otherId := range ids {
			assert.NotDeepEqual(
				t,
				id.Val,
				otherId.Val,
				"id equal to previous id")
		}
		ids = append(ids, id)
	}
}
