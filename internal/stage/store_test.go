package stage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

const (
	stageName    = "stage-name"
	stageAsset   = "asset-name"
	stageService = "ServiceName"
	stageMethod  = "MethodName"
)

func TestStore_Create(t *testing.T) {
	tests := []*Stage{
		{
			Name:    stageName,
			Asset:   stageAsset,
			Service: stageService,
			Method:  stageMethod,
		},
		{
			Name:    "",
			Asset:   "",
			Service: "",
			Method:  "",
		},
	}

	for _, cfg := range tests {
		testName := fmt.Sprintf("cfg=%v", cfg)

		t.Run(
			testName,
			func(t *testing.T) {
				st, ok := NewStore().(*store)
				assert.IsTrue(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				assert.IsNil(t, err, "create error")
				assert.DeepEqual(t, 1, lenStages(st), "store size")
				stored, ok := st.stages.Load(cfg.Name)
				assert.IsTrue(t, ok, "stage exists")
				s, ok := stored.(*Stage)
				assert.IsTrue(t, ok, "stage type assertion failed")
				assert.DeepEqual(t, cfg.Name, s.Name, "correct name")
				assert.DeepEqual(t, cfg.Asset, s.Asset, "correct asset")
				assert.DeepEqual(t, cfg.Service, s.Service, "correct service")
				assert.DeepEqual(t, cfg.Method, s.Method, "correct method")
			})
	}
}

func lenStages(st *store) int {
	count := 0
	st.stages.Range(
		func(key, value interface{}) bool {
			count += 1
			return true
		})
	return count
}
