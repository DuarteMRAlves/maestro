package orchestration

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

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

				db := storage.NewTestDb(t)
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
	helper := storage.NewTxnHelper(txn)
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
