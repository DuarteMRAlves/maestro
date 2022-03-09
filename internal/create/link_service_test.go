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
				existsStageCount := 0

				possibleStages := []domain.StageName{
					test.expLink.Source().Stage(),
					test.expLink.Target().Stage(),
				}
				existsStage := existsOneOfStageFn(
					possibleStages,
					&existsStageCount,
				)

				linkStore := mockLinkStorage{links: map[domain.LinkName]Link{}}

				orchStore := mockOrchestrationStorage{
					orchs: map[domain.OrchestrationName]Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateLink(
					linkStore,
					existsStage,
					orchStore,
				)
				res := createFn(test.req)

				assert.Assert(t, !res.Err.Present())

				assert.Equal(t, 1, len(linkStore.links))
				l, exists := linkStore.links[test.expLink.Name()]
				assert.Assert(t, exists)
				assertEqualLink(t, test.expLink, l)

				// Two because of source and target
				assert.Equal(t, existsStageCount, 2)

				assert.Equal(t, 1, len(orchStore.orchs))
				o, exists := orchStore.orchs[test.expOrchestration.Name()]
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

	existsStageCount := 0

	possibleStages := []domain.StageName{
		expLink.Source().Stage(),
		expLink.Target().Stage(),
	}
	existsStage := existsOneOfStageFn(possibleStages, &existsStageCount)

	linkStore := mockLinkStorage{links: map[domain.LinkName]Link{}}

	storage := mockOrchestrationStorage{
		orchs: map[domain.OrchestrationName]Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := CreateLink(linkStore, existsStage, storage)
	res := createFn(req)

	assert.Assert(t, !res.Err.Present())

	assert.Equal(t, 1, len(linkStore.links))
	l, exists := linkStore.links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)

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
	assert.Equal(t, 1, len(linkStore.links))
	l, exists = linkStore.links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)
	// Two because of source and target
	assert.Equal(t, existsStageCount, 2)

	assert.Equal(t, 1, len(storage.orchs))
	o, exists = storage.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

type mockLinkStorage struct {
	links map[domain.LinkName]Link
}

func (m mockLinkStorage) Save(l Link) LinkResult {
	m.links[l.Name()] = l
	return SomeLink(l)
}

func (m mockLinkStorage) Load(name domain.LinkName) LinkResult {
	l, exists := m.links[name]
	if !exists {
		err := errdefs.NotFoundWithMsg("link not found: %s", l.Name())
		return ErrLink(err)
	}
	return SomeLink(l)
}

func (m mockLinkStorage) Verify(name domain.LinkName) bool {
	_, exists := m.links[name]
	return exists
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

func createLink(
	t *testing.T,
	linkName, orchestrationName string,
	requiredOnly bool,
) Link {
	name, err := domain.NewLinkName(linkName)
	assert.NilError(t, err, "create name for link %s", linkName)

	sourceStage, err := domain.NewStageName("source")
	assert.NilError(t, err, "create source stage for link %s", linkName)
	sourceFieldOpt := domain.NewEmptyMessageField()
	if !requiredOnly {
		sourceField, err := domain.NewMessageField("source-field")
		assert.NilError(t, err, "create source field for link %s", linkName)
		sourceFieldOpt = domain.NewPresentMessageField(sourceField)
	}
	sourceEndpoint := NewLinkEndpoint(sourceStage, sourceFieldOpt)

	targetStage, err := domain.NewStageName("target")
	assert.NilError(t, err, "create target stage for link %s", linkName)
	targetFieldOpt := domain.NewEmptyMessageField()
	if !requiredOnly {
		targetField, err := domain.NewMessageField("target-field")
		assert.NilError(t, err, "create target field for link %s", linkName)
		targetFieldOpt = domain.NewPresentMessageField(targetField)
	}
	targetEndpoint := NewLinkEndpoint(targetStage, targetFieldOpt)

	orchestration, err := domain.NewOrchestrationName(orchestrationName)
	assert.NilError(t, err, "create orchestration for link %s", linkName)

	return NewLink(name, sourceEndpoint, targetEndpoint, orchestration)
}

func assertEqualLink(t *testing.T, expected, actual Link) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assertEqualEndpoint(t, expected.Source(), actual.Source())
	assertEqualEndpoint(t, expected.Target(), actual.Target())
	assert.Equal(t, expected.Orchestration(), actual.Orchestration())
}

func assertEqualEndpoint(t *testing.T, expected, actual LinkEndpoint) {
	assert.Equal(t, expected.Stage().Unwrap(), actual.Stage().Unwrap())
	assert.Equal(t, expected.Field().Present(), actual.Field().Present())
	if expected.Field().Present() {
		assert.Equal(t, expected.Field().Unwrap(), actual.Field().Unwrap())
	}
}
