package arch

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_CreateStage(t *testing.T) {
	const name api.StageName = "Stage-Name"
	tests := []struct {
		name     string
		req      *api.CreateStageRequest
		expected *api.Stage
	}{
		{
			name: "required parameters",
			req: &api.CreateStageRequest{
				Name:          name,
				Orchestration: api.OrchestrationName("orchestration-0"),
				Address:       "some-address",
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-address",
				Orchestration: api.OrchestrationName("orchestration-0"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "custom host and port",
			req: &api.CreateStageRequest{
				Name:          name,
				Host:          "some-host",
				Port:          12345,
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-host:12345",
				Orchestration: api.OrchestrationName("orchestration-0"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "all parameters",
			req: &api.CreateStageRequest{
				Name:          name,
				Service:       "Service",
				Rpc:           "Rpc",
				Address:       "Address",
				Host:          "",
				Port:          0,
				Orchestration: api.OrchestrationName("orchestration-0"),
				Asset:         api.AssetName("asset-0"),
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "Service",
				Rpc:           "Rpc",
				Address:       "Address",
				Orchestration: api.OrchestrationName("orchestration-0"),
				Asset:         api.AssetName("asset-0"),
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateStage(t, test.req, test.expected)
			},
		)
	}
}

func testCreateStage(
	t *testing.T,
	req *api.CreateStageRequest,
	expected *api.Stage,
) {
	var (
		err           error
		stored        api.Stage
		orchestration api.Orchestration
	)
	db := kv.NewTestDb(t)
	defer db.Close()

	err = db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			a := &api.Asset{Name: api.AssetName("asset-0")}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Orchestration{Name: api.OrchestrationName("orchestration-0")}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			createStage := CreateStageWithTxn(txn)
			return createStage(req)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			return helper.LoadStage(&stored, req.Name)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, expected.Name, stored.Name, "name not equal")
	assert.Equal(t, expected.Phase, stored.Phase, "phase not equal")
	assert.Equal(t, expected.Service, stored.Service, "service not equal")
	assert.Equal(t, expected.Rpc, stored.Rpc, "rpc not equal")
	assert.Equal(t, expected.Address, stored.Address, "address not equal")
	assert.Equal(
		t,
		expected.Orchestration,
		stored.Orchestration,
		"orchestration not equal",
	)
	assert.Equal(t, expected.Asset, stored.Asset, "asset not equal")

	err = db.View(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			return helper.LoadOrchestration(
				&orchestration,
				stored.Orchestration,
			)
		},
	)
	found := false
	for _, s := range orchestration.Stages {
		if s == stored.Name {
			found = true
		}
	}
	assert.Assert(t, found, "stage is not in orchestration")
}

func TestManager_CreateStage_Error(t *testing.T) {
	tests := []struct {
		name            string
		req             *api.CreateStageRequest
		assertErrTypeFn func(error) bool
		expectedErrMsg  string
	}{
		{
			name:            "nil config",
			req:             nil,
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "'req' is nil",
		},
		{
			name: "empty name",
			req: &api.CreateStageRequest{
				Name:          "",
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name ''",
		},
		{
			name: "invalid name",
			req: &api.CreateStageRequest{
				Name:          "some#name",
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name 'some#name'",
		},
		{
			name: "stage already exists",
			req: &api.CreateStageRequest{
				Name:          "duplicate",
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsAlreadyExists,
			expectedErrMsg:  "stage 'duplicate' already exists",
		},
		{
			name: "orchestration not found",
			req: &api.CreateStageRequest{
				Name:          "some-stage",
				Orchestration: "unknown",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "orchestration 'unknown' not found",
		},
		{
			name: "asset not found",
			req: &api.CreateStageRequest{
				Name:          "some-stage",
				Asset:         "unknown",
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "asset 'unknown' not found",
		},
		{
			name: "stage and host specified",
			req: &api.CreateStageRequest{
				Name:          "some-stage",
				Address:       "some-address",
				Host:          "some-host",
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "Cannot simultaneously specify address and host for stage",
		},
		{
			name: "stage and port specified",
			req: &api.CreateStageRequest{
				Name:          "some-stage",
				Address:       "some-address",
				Port:          23456,
				Orchestration: "orchestration-0",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "Cannot simultaneously specify address and port for stage",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateStageError(
					t,
					test.req,
					test.assertErrTypeFn,
					test.expectedErrMsg,
				)
			},
		)
	}
}

func testCreateStageError(
	t *testing.T,
	req *api.CreateStageRequest,
	assertErrTypeFn func(error) bool,
	expectedErrMsg string,
) {
	db := kv.NewTestDb(t)
	defer db.Close()

	err := db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			a := &api.Asset{Name: api.AssetName("asset-0")}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Orchestration{
				Name: api.OrchestrationName("orchestration-0"),
			}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			s := &api.Stage{Name: "duplicate"}
			if err := helper.SaveStage(s); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			createStage := CreateStageWithTxn(txn)
			return createStage(req)
		},
	)
	assert.Assert(t, assertErrTypeFn(err), "wrong error type")
	assert.Equal(t, expectedErrMsg, err.Error(), "wrong error message")
}

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

				db := kv.NewTestDb(t)
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
	helper := kv.NewTxnHelper(txn)
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
