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

func TestManager_CreateLink(t *testing.T) {
	tests := []struct {
		name     string
		req      *api.CreateLinkRequest
		expected *api.Link
	}{
		{
			name: "required parameters",
			req: &api.CreateLinkRequest{
				Name:          api.LinkName("link-0"),
				SourceStage:   api.StageName("stage-0"),
				TargetStage:   api.StageName("stage-1"),
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			expected: &api.Link{
				Name:          api.LinkName("link-0"),
				SourceStage:   api.StageName("stage-0"),
				SourceField:   "",
				TargetStage:   api.StageName("stage-1"),
				TargetField:   "",
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
		},
		{
			name: "all parameters",
			req: &api.CreateLinkRequest{
				Name:          api.LinkName("link-0"),
				SourceStage:   api.StageName("stage-0"),
				SourceField:   "source-field-0",
				TargetStage:   api.StageName("stage-1"),
				TargetField:   "target-field-0",
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			expected: &api.Link{
				Name:          api.LinkName("link-0"),
				SourceStage:   api.StageName("stage-0"),
				SourceField:   "source-field-0",
				TargetStage:   api.StageName("stage-1"),
				TargetField:   "target-field-0",
				Orchestration: api.OrchestrationName("orchestration-0"),
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
	db := kv.NewTestDb(t)
	defer db.Close()

	err := db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
			o := &api.Orchestration{Name: api.OrchestrationName("orchestration-0")}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			s := &api.Stage{
				Name:          api.StageName("stage-0"),
				Phase:         api.StagePending,
				Orchestration: api.OrchestrationName("orchestration-0"),
			}
			if err := helper.SaveStage(s); err != nil {
				return err
			}
			t := &api.Stage{
				Name:          api.StageName("stage-1"),
				Phase:         api.StagePending,
				Orchestration: api.OrchestrationName("orchestration-0"),
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
			createLink := CreateLinkWithTxn(txn)
			return createLink(req)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)
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
			helper := kv.NewTxnHelper(txn)
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
				TargetStage:   api.StageName("stage-1"),
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "empty source stage name",
		},
		{
			name: "empty target name",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   api.StageName("stage-0"),
				TargetStage:   "",
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "empty target stage name",
		},
		{
			name: "source stage not found",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "unknown",
				TargetStage:   api.StageName("stage-1"),
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "source stage 'unknown' not found",
		},
		{
			name: "target stage not found",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   api.StageName("stage-0"),
				TargetStage:   "unknown",
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "target stage 'unknown' not found",
		},
		{
			name: "source stage different orchestration",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "stage-different",
				TargetStage:   api.StageName("stage-1"),
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"orchestration for link '%s' is '%s' but source stage is registered in '%s'.",
				"some-link",
				api.OrchestrationName("orchestration-0"),
				"different-orchestration",
			),
		},
		{
			name: "target stage different orchestration",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   api.StageName("stage-0"),
				TargetStage:   "stage-different",
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			assertErrTypeFn: errdefs.IsFailedPrecondition,
			expectedErrMsg: fmt.Sprintf(
				"orchestration for link '%s' is '%s' but target stage is registered in '%s'.",
				"some-link",
				api.OrchestrationName("orchestration-0"),
				"different-orchestration",
			),
		},
		{
			name: "source stage not pending",
			req: &api.CreateLinkRequest{
				Name:          "some-link",
				SourceStage:   "not-pending",
				TargetStage:   api.StageName("stage-1"),
				Orchestration: api.OrchestrationName("orchestration-0"),
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
				SourceStage:   api.StageName("stage-0"),
				TargetStage:   "not-pending",
				Orchestration: api.OrchestrationName("orchestration-0"),
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
	db := kv.NewTestDb(t)
	defer db.Close()

	// Prepare tests
	// Number 0 has an orchestration and correct stages created to be used.
	err := db.Update(
		func(txn *badger.Txn) error {
			helper := kv.NewTxnHelper(txn)

			o := &api.Orchestration{Name: api.OrchestrationName("orchestration-0")}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}

			s := &api.Stage{
				Name:          api.StageName("stage-0"),
				Phase:         api.StagePending,
				Orchestration: api.OrchestrationName("orchestration-0"),
			}
			if err := helper.SaveStage(s); err != nil {
				return err
			}

			t := &api.Stage{
				Name:          api.StageName("stage-1"),
				Phase:         api.StagePending,
				Orchestration: api.OrchestrationName("orchestration-0"),
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
				Orchestration: api.OrchestrationName("orchestration-0"),
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
			createLink := CreateLinkWithTxn(txn)
			return createLink(req)
		},
	)
	assert.Assert(t, assertErrTypeFn(err), "wrong error type")
	assert.Equal(t, expectedErrMsg, err.Error(), "wrong error message")
}

func TestManager_GetMatchingLinks(t *testing.T) {
	tests := []struct {
		name   string
		req    *api.GetLinkRequest
		stored []*api.Link
		// names of the expected links
		expected []api.LinkName
	}{
		{
			name:     "zero elements stored, nil req",
			req:      nil,
			stored:   []*api.Link{},
			expected: []api.LinkName{},
		},
		{
			name:     "zero elements stored, some req",
			req:      &api.GetLinkRequest{Name: "some-name"},
			stored:   []*api.Link{},
			expected: []api.LinkName{},
		},
		{
			name: "one element stored, nil req",
			req:  nil,
			stored: []*api.Link{
				testLink(0),
			},
			expected: []api.LinkName{api.LinkName("link-0")},
		},
		{
			name: "one element stored, matching name req",
			req:  &api.GetLinkRequest{Name: api.LinkName("link-0")},
			stored: []*api.Link{
				testLink(0),
			},
			expected: []api.LinkName{api.LinkName("link-0")},
		},
		{
			name: "one element stored, non-matching name req",
			req:  &api.GetLinkRequest{Name: api.LinkName("link-1")},
			stored: []*api.Link{
				testLink(2),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, nil req",
			req:  nil,
			stored: []*api.Link{
				testLink(1),
				testLink(5),
				testLink(3),
			},
			expected: []api.LinkName{
				api.LinkName("link-1"),
				api.LinkName("link-3"),
				api.LinkName("link-5"),
			},
		},
		{
			name: "multiple elements stored, matching name req",
			req:  &api.GetLinkRequest{Name: api.LinkName("link-2")},
			stored: []*api.Link{
				testLink(3),
				testLink(1),
				testLink(2),
			},
			expected: []api.LinkName{api.LinkName("link-2")},
		},
		{
			name: "multiple elements stored, non-matching name req",
			req:  &api.GetLinkRequest{Name: api.LinkName("link-2")},
			stored: []*api.Link{
				testLink(0),
				testLink(3),
				testLink(1),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, matching source stage req",
			req:  &api.GetLinkRequest{SourceStage: api.StageName("stage-4")},
			stored: []*api.Link{
				testLink(3),
				testLink(4),
				testLink(2),
			},
			expected: []api.LinkName{api.LinkName("link-4")},
		},
		{
			name: "multiple elements stored, non-matching source stage req",
			req:  &api.GetLinkRequest{SourceStage: api.StageName("stage-4")},
			stored: []*api.Link{
				testLink(0),
				testLink(3),
				testLink(1),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, matching source field req",
			req:  &api.GetLinkRequest{SourceField: "source-field-1"},
			stored: []*api.Link{
				testLink(1),
				testLink(4),
				testLink(2),
			},
			expected: []api.LinkName{api.LinkName("link-1")},
		},
		{
			name: "multiple elements stored, non-matching source field req",
			req:  &api.GetLinkRequest{SourceField: "source-field-1"},
			stored: []*api.Link{
				testLink(0),
				testLink(3),
				testLink(2),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, matching target stage req",
			req:  &api.GetLinkRequest{TargetStage: api.StageName("stage-4")},
			stored: []*api.Link{
				testLink(3),
				testLink(4),
				testLink(2),
			},
			expected: []api.LinkName{api.LinkName("link-3")},
		},
		{
			name: "multiple elements stored, non-matching target stage req",
			req:  &api.GetLinkRequest{TargetStage: api.StageName("stage-4")},
			stored: []*api.Link{
				testLink(0),
				testLink(4),
				testLink(1),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, matching target field req",
			req:  &api.GetLinkRequest{TargetField: "target-field-3"},
			stored: []*api.Link{
				testLink(1),
				testLink(3),
				testLink(2),
			},
			expected: []api.LinkName{api.LinkName("link-3")},
		},
		{
			name: "multiple elements stored, non-matching target field req",
			req:  &api.GetLinkRequest{TargetField: "target-field-3"},
			stored: []*api.Link{
				testLink(0),
				testLink(1),
				testLink(2),
			},
			expected: []api.LinkName{},
		},
		{
			name: "multiple elements stored, matching orchestration req",
			req: &api.GetLinkRequest{
				Orchestration: api.OrchestrationName("orchestration-0"),
			},
			stored: []*api.Link{
				testLink(0),
				testLink(3),
				testLink(1),
			},
			expected: []api.LinkName{api.LinkName("link-0")},
		},
		{
			name: "multiple elements stored, non-matching orchestration req",
			req: &api.GetLinkRequest{
				Orchestration: api.OrchestrationName("orchestration-2"),
			},
			stored: []*api.Link{
				testLink(0),
				testLink(3),
				testLink(4),
			},
			expected: []api.LinkName{},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					err      error
					received []*api.Link
				)

				db := kv.NewTestDb(t)
				defer db.Close()

				for _, l := range test.stored {
					err = db.Update(
						func(txn *badger.Txn) error {
							return saveLinkAndDependencies(txn, l)
						},
					)
				}

				err = db.View(
					func(txn *badger.Txn) error {
						getLinks := GetLinksWithTxn(txn)
						received, err = getLinks(test.req)
						return err
					},
				)
				assert.NilError(t, err, "get orchestration")
				assert.Equal(t, len(test.expected), len(received))

				seen := make(map[api.LinkName]bool, 0)
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

func testLink(num int) *api.Link {
	return &api.Link{
		Name:        api.LinkName(fmt.Sprintf("link-%d", num)),
		SourceStage: api.StageName(fmt.Sprintf("stage-%d", num)),
		SourceField: fmt.Sprintf("source-field-%d", num),
		TargetStage: api.StageName(fmt.Sprintf("stage-%d", num+1)),
		TargetField: fmt.Sprintf("target-field-%d", num),
		Orchestration: api.OrchestrationName(
			fmt.Sprintf(
				"orchestration-%d",
				num,
			),
		),
	}
}

func saveLinkAndDependencies(txn *badger.Txn, l *api.Link) error {
	helper := kv.NewTxnHelper(txn)
	if !helper.ContainsOrchestration(l.Orchestration) {
		err := helper.SaveOrchestration(
			orchestrationForName(
				l.Orchestration,
				api.OrchestrationRunning,
			),
		)
		if err != nil {
			return err
		}
	}
	if !helper.ContainsStage(l.SourceStage) {
		err := helper.SaveStage(
			&api.Stage{
				Name:          l.SourceStage,
				Orchestration: l.Orchestration,
			},
		)
		if err != nil {
			return err
		}
	}
	if !helper.ContainsStage(l.TargetStage) {
		err := helper.SaveStage(
			&api.Stage{
				Name:          l.TargetStage,
				Orchestration: l.Orchestration,
			},
		)
		if err != nil {
			return err
		}
	}
	return helper.SaveLink(l)
}
