package stage

import (
	"fmt"
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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
				testing2.IsTrue(t, ok, "type assertion failed for store")

				err := st.Create(cfg)
				testing2.IsNil(t, err, "create error")
				testing2.DeepEqual(t, 1, lenStages(st), "store size")
				stored, ok := st.stages.Load(cfg.Name)
				testing2.IsTrue(t, ok, "stage exists")
				s, ok := stored.(*Stage)
				testing2.IsTrue(t, ok, "stage type assertion failed")
				testing2.DeepEqual(t, cfg.Name, s.Name, "correct name")
				testing2.DeepEqual(t, cfg.Asset, s.Asset, "correct asset")
				testing2.DeepEqual(t, cfg.Service, s.Service, "correct service")
				testing2.DeepEqual(t, cfg.Method, s.Method, "correct method")
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
