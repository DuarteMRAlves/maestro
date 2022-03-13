package orchestration

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_GetMatchingStages(t *testing.T) {
	tests := []struct {
		name   string
		req    *api.GetStageRequest
		stored []*api.Stage
		// names of the expected stages
		expected []api.StageName
	}{
		{
			name:     "zero elements stored, nil req",
			req:      nil,
			stored:   []*api.Stage{},
			expected: []api.StageName{},
		},
		{
			name:     "zero elements stored, some req",
			req:      &api.GetStageRequest{Name: "some-name"},
			stored:   []*api.Stage{},
			expected: []api.StageName{},
		},
		{
			name: "one element stored, nil req",
			req:  nil,
			stored: []*api.Stage{
				testStage(0, api.StageSucceeded),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "one element stored, matching name req",
			req:  &api.GetStageRequest{Name: api.StageName("stage-0")},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "one element stored, non-matching name req",
			req:  &api.GetStageRequest{Name: api.StageName("stage-1")},
			stored: []*api.Stage{
				testStage(2, api.StagePending),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, nil req",
			req:  nil,
			stored: []*api.Stage{
				testStage(1, api.StagePending),
				testStage(5, api.StageSucceeded),
				testStage(3, api.StageFailed),
			},
			expected: []api.StageName{
				api.StageName("stage-1"),
				api.StageName("stage-3"),
				api.StageName("stage-5"),
			},
		},
		{
			name: "multiple elements stored, matching name req",
			req:  &api.GetStageRequest{Name: api.StageName("stage-2")},
			stored: []*api.Stage{
				testStage(3, api.StageRunning),
				testStage(1, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-2")},
		},
		{
			name: "multiple elements stored, non-matching name req",
			req:  &api.GetStageRequest{Name: api.StageName("stage-2")},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(3, api.StagePending),
				testStage(1, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching phase req",
			req: &api.GetStageRequest{
				Phase: api.StageFailed,
			},
			stored: []*api.Stage{
				testStage(1, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(0, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "multiple elements stored, non-matching phase req",
			req: &api.GetStageRequest{
				Phase: api.StageSucceeded,
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(2, api.StagePending),
				testStage(1, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching service req",
			req: &api.GetStageRequest{
				Service: "service-2",
			},
			stored: []*api.Stage{
				testStage(1, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-2")},
		},
		{
			name: "multiple elements stored, non-matching service req",
			req: &api.GetStageRequest{
				Service: "service-4",
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(2, api.StagePending),
				testStage(1, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching rpc req",
			req: &api.GetStageRequest{
				Rpc: "rpc-0",
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "multiple elements stored, non-matching rpc req",
			req: &api.GetStageRequest{
				Rpc: "rpc-2",
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(3, api.StagePending),
				testStage(1, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching address req",
			req: &api.GetStageRequest{
				Address: "address-1",
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-1")},
		},
		{
			name: "multiple elements stored, non-matching address req",
			req: &api.GetStageRequest{
				Address: "address-1",
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(3, api.StagePending),
				testStage(2, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching orchestration req",
			req: &api.GetStageRequest{
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "multiple elements stored, non-matching orchestration req",
			req: &api.GetStageRequest{
				Orchestration: api.OrchestrationName("orchestration-2"),
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(3, api.StagePending),
				testStage(4, api.StageRunning),
			},
			expected: []api.StageName{},
		},
		{
			name: "multiple elements stored, matching asset req",
			req: &api.GetStageRequest{
				Asset: api.AssetName("asset-0"),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{api.StageName("stage-0")},
		},
		{
			name: "multiple elements stored, non-matching asset req",
			req: &api.GetStageRequest{
				Asset: api.AssetName("asset-2"),
			},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
				testStage(3, api.StagePending),
				testStage(1, api.StageRunning),
			},
			expected: []api.StageName{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					err      error
					received []*api.Stage
				)

				db := storage.NewTestDb(t)
				defer db.Close()

				for _, s := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							return saveStageAndDependencies(txn, s)
						},
					)
				}

				err = db.View(
					func(txn *badger.Txn) error {
						getStages := GetStagesWithTxn(txn)
						received, err = getStages(test.req)
						return err
					},
				)
				assert.NilError(t, err, "get orchestration")
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[api.StageName]bool, 0)
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

func testStage(num int, phase api.StagePhase) *api.Stage {
	return &api.Stage{
		Name:    api.StageName(fmt.Sprintf("stage-%d", num)),
		Phase:   phase,
		Service: fmt.Sprintf("service-%d", num),
		Rpc:     fmt.Sprintf("rpc-%d", num),
		Address: fmt.Sprintf("address-%d", num),
		Orchestration: api.OrchestrationName(
			fmt.Sprintf(
				"orchestration-%d",
				num,
			),
		),
		Asset: api.AssetName(fmt.Sprintf("asset-%d", num)),
	}
}

func saveStageAndDependencies(txn *badger.Txn, s *api.Stage) error {
	helper := storage.NewTxnHelper(txn)
	if !helper.ContainsOrchestration(s.Orchestration) {
		err := helper.SaveOrchestration(
			orchestrationForName(
				s.Orchestration,
				api.OrchestrationRunning,
			),
		)
		if err != nil {
			return err
		}
	}
	if !helper.ContainsAsset(s.Asset) {
		err := helper.SaveAsset(&api.Asset{Name: s.Asset})
		if err != nil {
			return err
		}
	}
	return helper.SaveStage(s)
}
