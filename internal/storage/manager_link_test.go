package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_CreateLink(t *testing.T) {
	tests := []struct {
		name     string
		req      *api.CreateLinkRequest
		expected *api.Link
	}{
		{
			name: "required parameters",
			req: &api.CreateLinkRequest{
				Name:        util.LinkNameForNum(0),
				SourceStage: util.LinkSourceStageForNum(0),
				TargetStage: util.LinkTargetStageForNum(0),
			},
			expected: &api.Link{
				Name:          util.LinkNameForNum(0),
				SourceStage:   util.LinkSourceStageForNum(0),
				SourceField:   "",
				TargetStage:   util.LinkTargetStageForNum(0),
				TargetField:   "",
				Orchestration: defaultOrchestrationName,
			},
		},
		{
			name: "all parameters",
			req: &api.CreateLinkRequest{
				Name:          util.LinkNameForNum(0),
				SourceStage:   util.LinkSourceStageForNum(0),
				SourceField:   util.LinkSourceFieldForNum(0),
				TargetStage:   util.LinkTargetStageForNum(0),
				TargetField:   util.LinkTargetFieldForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
			expected: &api.Link{
				Name:          util.LinkNameForNum(0),
				SourceStage:   util.LinkSourceStageForNum(0),
				SourceField:   util.LinkSourceFieldForNum(0),
				TargetStage:   util.LinkTargetStageForNum(0),
				TargetField:   util.LinkTargetFieldForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateLink(t, test.req, test.expected)
			},
		)
	}
}

func testCreateLink(
	t *testing.T,
	req *api.CreateLinkRequest,
	expected *api.Link,
) {
	var (
		stored        api.Link
		orchestration api.Orchestration
	)
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	err = db.Update(
		func(txn *badger.Txn) error {
			orchestrationName := defaultOrchestrationName
			if req.Orchestration != "" {
				orchestrationName = req.Orchestration
			}
			helper := NewTxnHelper(txn)
			o := &api.Orchestration{Name: util.OrchestrationNameForNum(0)}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			s := &api.Stage{
				Name:          util.LinkSourceStageForNum(0),
				Phase:         api.StagePending,
				Orchestration: orchestrationName,
			}
			if err := helper.SaveStage(s); err != nil {
				return err
			}
			t := &api.Stage{
				Name:          util.LinkTargetStageForNum(0),
				Phase:         api.StagePending,
				Orchestration: orchestrationName,
			}
			if err := helper.SaveStage(t); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateLink(txn, req)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			helper := TxnHelper{txn: txn}
			return helper.LoadLink(&stored, req.Name)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, expected.Name, stored.Name)
	assert.Equal(t, expected.SourceStage, stored.SourceStage)
	assert.Equal(t, expected.SourceField, stored.SourceField)
	assert.Equal(t, expected.TargetStage, stored.TargetStage)
	assert.Equal(t, expected.TargetField, stored.TargetField)
	assert.Equal(t, expected.Orchestration, stored.Orchestration)

	err = db.View(
		func(txn *badger.Txn) error {
			helper := TxnHelper{txn: txn}
			return helper.LoadOrchestration(
				&orchestration,
				stored.Orchestration,
			)
		},
	)
	found := false
	for _, l := range orchestration.Links {
		if l == stored.Name {
			found = true
		}
	}
	assert.Assert(t, found, "link is not in orchestration")
}

func TestManager_CreateLink_Error(t *testing.T) {
	tests := []struct {
		name            string
		req             *api.CreateLinkRequest
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
			req:             &api.CreateLinkRequest{Name: ""},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name ''",
		},
		{
			name:            "invalid name",
			req:             &api.CreateLinkRequest{Name: "some//name"},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name 'some//name'",
		},
		{
			name:            "link already exists",
			req:             &api.CreateLinkRequest{Name: "duplicate"},
			assertErrTypeFn: errdefs.IsAlreadyExists,
			expectedErrMsg:  "link 'duplicate' already exists",
		},
		{
			name: "orchestration not found",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				Orchestration: "unknown",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "orchestration 'unknown' not found",
		},
		{
			name: "empty source name",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "",
				TargetStage:   util.LinkTargetStageForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "empty source stage name",
		},
		{
			name: "empty target name",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   util.LinkSourceStageForNum(0),
				TargetStage:   "",
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "empty target stage name",
		},
		{
			name: "source stage not found",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "unknown",
				TargetStage:   util.LinkTargetStageForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "source stage 'unknown' not found",
		},
		{
			name: "target stage not found",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   util.LinkSourceStageForNum(0),
				TargetStage:   "unknown",
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "target stage 'unknown' not found",
		},
		{
			name: "source stage different orchestration",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "stage-different",
				TargetStage:   util.LinkTargetStageForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"orchestration for link '%s' is '%s' but source stage is registered in '%s'.",
				"some-link",
				util.OrchestrationNameForNum(0),
				"different-orchestration",
			),
		},
		{
			name: "target stage different orchestration",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   util.LinkSourceStageForNum(0),
				TargetStage:   "stage-different",
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"orchestration for link '%s' is '%s' but target stage is registered in '%s'.",
				"some-link",
				util.OrchestrationNameForNum(0),
				"different-orchestration",
			),
		},
		{
			name: "source stage not pending",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "not-pending",
				TargetStage:   util.LinkTargetStageForNum(0),
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"source stage is not in Pending phase for link '%s'.",
				"some-link",
			),
		},
		{
			name: "target stage not pending",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   util.LinkSourceStageForNum(0),
				TargetStage:   "not-pending",
				Orchestration: util.OrchestrationNameForNum(0),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"target stage is not in Pending phase for link '%s'.",
				"some-link",
			),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateLinkError(
					t,
					test.req,
					test.assertErrTypeFn,
					test.expectedErrMsg,
				)
			},
		)
	}
}

func testCreateLinkError(
	t *testing.T,
	req *api.CreateLinkRequest,
	assertErrTypeFn func(error) bool,
	expectedErrMsg string,
) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	// Prepare tests
	// Number 0 has an orchestration and correct stages created to be used.
	err = db.Update(
		func(txn *badger.Txn) error {
			helper := NewTxnHelper(txn)

			o := &api.Orchestration{Name: util.OrchestrationNameForNum(0)}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}

			s := &api.Stage{
				Name:          util.LinkSourceStageForNum(0),
				Phase:         api.StagePending,
				Orchestration: util.OrchestrationNameForNum(0),
			}
			if err := helper.SaveStage(s); err != nil {
				return err
			}

			t := &api.Stage{
				Name:          util.LinkTargetStageForNum(0),
				Phase:         api.StagePending,
				Orchestration: util.OrchestrationNameForNum(0),
			}
			if err := helper.SaveStage(t); err != nil {
				return err
			}

			differentOrchestration := &api.Stage{
				Name:          "stage-different",
				Phase:         api.StagePending,
				Orchestration: "different-orchestration",
			}
			if err := helper.SaveStage(differentOrchestration); err != nil {
				return err
			}
			notPending := &api.Stage{
				Name:          "not-pending",
				Phase:         api.StageSucceeded,
				Orchestration: util.OrchestrationNameForNum(0),
			}
			if err := helper.SaveStage(notPending); err != nil {
				return err
			}

			l := &api.Link{Name: "duplicate"}
			if err := helper.SaveLink(l); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateLink(txn, req)
		},
	)
	assert.Assert(t, assertErrTypeFn(err), "wrong error type")
	assert.Equal(t, expectedErrMsg, err.Error(), "wrong error message")
}
