package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateLink(t *testing.T) {
	tests := []struct {
		name              string
		req               LinkRequest
		expLink           Link
		loadOrchestration Orchestration
		expOrchestration  Orchestration
	}{
		{
			name: "required fields",
			req: LinkRequest{
				Name:          "some-name",
				SourceStage:   "source",
				SourceField:   domain.NewEmptyString(),
				TargetStage:   "target",
				TargetField:   domain.NewEmptyString(),
				Orchestration: "orchestration",
			},
			expLink: createLink(
				t,
				"some-name",
				"orchestration",
				true,
			),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{},
			),
			expOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{"some-name"},
			),
		},
		{
			name: "all fields",
			req: LinkRequest{
				Name:          "some-name",
				SourceStage:   "source",
				SourceField:   domain.NewPresentString("source-field"),
				TargetStage:   "target",
				TargetField:   domain.NewPresentString("target-field"),
				Orchestration: "orchestration",
			},
			expLink: createLink(
				t,
				"some-name",
				"orchestration",
				false,
			),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{},
			),
			expOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{"some-name"},
			),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				existsLinkCount := 0
				saveLinkCount := 0

				existsStageCount := 0

				existsLink := existsLinkFn(
					test.expLink.Name(),
					&existsLinkCount,
					1,
				)
				saveLink := saveLinkFn(t, test.expLink, &saveLinkCount)

				possibleStages := []domain.StageName{
					test.expLink.Source().Stage(),
					test.expLink.Target().Stage(),
				}
				existsStage := existsOneOfStageFn(
					possibleStages,
					&existsStageCount,
				)

				storage := mockOrchestrationStorage{
					orchs: map[domain.OrchestrationName]Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateLink(
					existsLink,
					saveLink,
					existsStage,
					storage,
				)
				res := createFn(test.req)

				assert.Assert(t, !res.Err.Present())
				assert.Equal(t, existsLinkCount, 1)
				assert.Equal(t, saveLinkCount, 1)
				// Two because of source and target
				assert.Equal(t, existsStageCount, 2)

				assert.Equal(t, 1, len(storage.orchs))
				o, exists := storage.orchs[test.expOrchestration.Name()]
				assert.Assert(t, exists)
				assertEqualOrchestration(t, test.expOrchestration, o)
			},
		)
	}
}

func TestCreateLink_AlreadyExists(t *testing.T) {
	req := LinkRequest{
		Name:          "some-name",
		SourceStage:   "source",
		SourceField:   domain.NewPresentString("source-field"),
		TargetStage:   "target",
		TargetField:   domain.NewPresentString("target-field"),
		Orchestration: "orchestration",
	}
	expLink := createLink(t, "some-name", "orchestration", false)
	storedOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"source", "target"},
		[]string{},
	)
	expOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"source", "target"},
		[]string{"some-name"},
	)

	existsLinkCount := 0
	saveLinkCount := 0

	existsStageCount := 0

	existsLink := existsLinkFn(expLink.Name(), &existsLinkCount, 1)
	saveLink := saveLinkFn(t, expLink, &saveLinkCount)

	possibleStages := []domain.StageName{
		expLink.Source().Stage(),
		expLink.Target().Stage(),
	}
	existsStage := existsOneOfStageFn(possibleStages, &existsStageCount)

	storage := mockOrchestrationStorage{
		orchs: map[domain.OrchestrationName]Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := CreateLink(existsLink, saveLink, existsStage, storage)
	res := createFn(req)

	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsLinkCount, 1)
	assert.Equal(t, saveLinkCount, 1)
	// Two because of source and target
	assert.Equal(t, existsStageCount, 2)

	assert.Equal(t, 1, len(storage.orchs))
	o, exists := storage.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)

	res = createFn(req)

	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("link '%v' already exists", req.Name),
	)
	assert.Equal(t, existsLinkCount, 2)
	assert.Equal(t, saveLinkCount, 1)
	// Two because of source and target
	assert.Equal(t, existsStageCount, 2)

	assert.Equal(t, 1, len(storage.orchs))
	o, exists = storage.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

func existsOneOfStageFn(
	expected []domain.StageName,
	callCount *int,
) ExistsStage {
	return func(name domain.StageName) bool {
		*callCount++
		for _, s := range expected {
			if s.Unwrap() == name.Unwrap() {
				return true
			}
		}
		return false
	}
}
