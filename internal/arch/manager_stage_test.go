package arch

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
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
				Name: name,
			},
			expected: &api.Stage{
				Name:    name,
				Phase:   api.StagePending,
				Service: "",
				Rpc:     "",
				Address: fmt.Sprintf(
					"%s:%d",
					defaultStageHost,
					defaultStagePort,
				),
				Orchestration: api.OrchestrationName("default"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "custom address",
			req: &api.CreateStageRequest{
				Name:    name,
				Address: "some-address",
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-address",
				Orchestration: api.OrchestrationName("default"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "custom host and port",
			req: &api.CreateStageRequest{
				Name: name,
				Host: "some-host",
				Port: 12345,
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-host:12345",
				Orchestration: api.OrchestrationName("default"),
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
				Orchestration: util.OrchestrationNameForNum(0),
				Asset:         util.AssetNameForNum(0),
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "Service",
				Rpc:           "Rpc",
				Address:       "Address",
				Orchestration: util.OrchestrationNameForNum(0),
				Asset:         util.AssetNameForNum(0),
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
		stored        api.Stage
		orchestration api.Orchestration
	)
	db := kv.NewTestDb(t)
	defer db.Close()

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	err = db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			a := &api.Asset{Name: util.AssetNameForNum(0)}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Orchestration{Name: util.OrchestrationNameForNum(0)}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateStage(txn, req)
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
			name:            "empty name",
			req:             &api.CreateStageRequest{Name: ""},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name ''",
		},
		{
			name:            "invalid name",
			req:             &api.CreateStageRequest{Name: "some#name"},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name 'some#name'",
		},
		{
			name:            "stage already exists",
			req:             &api.CreateStageRequest{Name: "duplicate"},
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
				Name:  "some-stage",
				Asset: "unknown",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "asset 'unknown' not found",
		},
		{
			name: "stage and host specified",
			req: &api.CreateStageRequest{
				Name:    "some-stage",
				Address: "some-address",
				Host:    "some-host",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "Cannot simultaneously specify address and host for stage",
		},
		{
			name: "stage and port specified",
			req: &api.CreateStageRequest{
				Name:    "some-stage",
				Address: "some-address",
				Port:    23456,
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

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	err = db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			a := &api.Asset{Name: util.AssetNameForNum(0)}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Stage{Name: util.StageNameForNum(0)}
			if err := helper.SaveStage(o); err != nil {
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
			return m.CreateStage(txn, req)
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
			expected: []api.StageName{util.StageNameForNum(0)},
		},
		{
			name: "one element stored, matching name req",
			req:  &api.GetStageRequest{Name: util.StageNameForNum(0)},
			stored: []*api.Stage{
				testStage(0, api.StagePending),
			},
			expected: []api.StageName{util.StageNameForNum(0)},
		},
		{
			name: "one element stored, non-matching name req",
			req:  &api.GetStageRequest{Name: util.StageNameForNum(1)},
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
				util.StageNameForNum(1),
				util.StageNameForNum(3),
				util.StageNameForNum(5),
			},
		},
		{
			name: "multiple elements stored, matching name req",
			req:  &api.GetStageRequest{Name: util.StageNameForNum(2)},
			stored: []*api.Stage{
				testStage(3, api.StageRunning),
				testStage(1, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(2)},
		},
		{
			name: "multiple elements stored, non-matching name req",
			req:  &api.GetStageRequest{Name: util.StageNameForNum(2)},
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
			expected: []api.StageName{util.StageNameForNum(0)},
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
				Service: util.StageServiceForNum(2),
			},
			stored: []*api.Stage{
				testStage(1, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(2)},
		},
		{
			name: "multiple elements stored, non-matching service req",
			req: &api.GetStageRequest{
				Service: util.StageServiceForNum(4),
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
				Rpc: util.StageRpcForNum(0),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(2, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(0)},
		},
		{
			name: "multiple elements stored, non-matching rpc req",
			req: &api.GetStageRequest{
				Rpc: util.StageRpcForNum(2),
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
				Address: util.StageAddressForNum(1),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(1)},
		},
		{
			name: "multiple elements stored, non-matching address req",
			req: &api.GetStageRequest{
				Address: util.StageAddressForNum(1),
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
				Orchestration: util.OrchestrationNameForNum(0),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(0)},
		},
		{
			name: "multiple elements stored, non-matching orchestration req",
			req: &api.GetStageRequest{
				Orchestration: util.OrchestrationNameForNum(2),
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
				Asset: util.AssetNameForNum(0),
			},
			stored: []*api.Stage{
				testStage(0, api.StageRunning),
				testStage(3, api.StagePending),
				testStage(1, api.StageFailed),
			},
			expected: []api.StageName{util.StageNameForNum(0)},
		},
		{
			name: "multiple elements stored, non-matching asset req",
			req: &api.GetStageRequest{
				Asset: util.AssetNameForNum(2),
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
				var received []*api.Stage

				db := kv.NewTestDb(t)
				defer db.Close()

				m, err := NewManager(NewTestContext(db))
				assert.NilError(t, err, "manager creation")

				for _, s := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							return saveStageAndDependencies(txn, s)
						},
					)
				}

				err = db.View(
					func(txn *badger.Txn) error {
						received, err = m.GetMatchingStage(
							txn,
							test.req,
						)
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
		Name:          util.StageNameForNum(num),
		Phase:         phase,
		Service:       util.StageServiceForNum(num),
		Rpc:           util.StageRpcForNum(num),
		Address:       util.StageAddressForNum(num),
		Orchestration: util.OrchestrationNameForNum(num),
		Asset:         util.AssetNameForNum(num),
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
