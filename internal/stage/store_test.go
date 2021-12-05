package stage

import (
	"fmt"
	"gotest.tools/v3/assert"
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
			Name:    stageName,
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
				assert.Assert(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				assert.NilError(t, err, "create error")
				assert.Equal(t, 1, lenStages(st), "store size")
				stored, ok := st.stages.Load(cfg.Name)
				assert.Assert(t, ok, "stage exists")
				s, ok := stored.(*Stage)
				assert.Assert(t, ok, "stage type assertion failed")
				assert.Equal(t, cfg.Name, s.Name, "correct name")
				assert.Equal(t, cfg.Asset, s.Asset, "correct asset")
				assert.Equal(t, cfg.Service, s.Service, "correct service")
				assert.Equal(t, cfg.Method, s.Method, "correct method")
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
