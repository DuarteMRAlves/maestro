package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateLink(t *testing.T) {
	tests := []struct {
		name              string
		req               LinkRequest
		expLink           internal.Link
		loadOrchestration internal.Orchestration
		expOrchestration  internal.Orchestration
		storedStages      []internal.Stage
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
			storedStages: []internal.Stage{
				createStage(t, "source", "orchestration", true),
				createStage(t, "target", "orchestration", true),
			},
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
			storedStages: []internal.Stage{
				createStage(t, "source", "orchestration", false),
				createStage(t, "target", "orchestration", false),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				linkStore := mockLinkStorage{links: map[internal.LinkName]internal.Link{}}

				stageStore := mockStageStorage{
					stages: map[internal.StageName]internal.Stage{},
				}
				for _, s := range test.storedStages {
					stageStore.stages[s.Name()] = s
				}

				orchStore := mockOrchestrationStorage{
					orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateLink(linkStore, stageStore, orchStore)
				res := createFn(test.req)

				assert.Assert(t, !res.Err.Present())

				assert.Equal(t, 1, len(linkStore.links))
				l, exists := linkStore.links[test.expLink.Name()]
				assert.Assert(t, exists)
				assertEqualLink(t, test.expLink, l)

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
	storedStages := []internal.Stage{
		createStage(t, "source", "orchestration", false),
		createStage(t, "target", "orchestration", false),
	}

	linkStore := mockLinkStorage{links: map[internal.LinkName]internal.Link{}}

	stageStore := mockStageStorage{
		stages: map[internal.StageName]internal.Stage{},
	}
	for _, s := range storedStages {
		stageStore.stages[s.Name()] = s
	}

	orchStore := mockOrchestrationStorage{
		orchs: map[internal.OrchestrationName]internal.Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := CreateLink(linkStore, stageStore, orchStore)
	res := createFn(req)

	assert.Assert(t, !res.Err.Present())

	assert.Equal(t, 1, len(linkStore.links))
	l, exists := linkStore.links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)

	assert.Equal(t, 1, len(orchStore.orchs))
	o, exists := orchStore.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)

	res = createFn(req)

	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "link", alreadyExists.Type)
	assert.Equal(t, req.Name, alreadyExists.Ident)
	assert.Equal(t, 1, len(linkStore.links))
	l, exists = linkStore.links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)

	assert.Equal(t, 1, len(orchStore.orchs))
	o, exists = orchStore.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

type mockLinkStorage struct {
	links map[internal.LinkName]internal.Link
}

func (m mockLinkStorage) Save(l internal.Link) error {
	m.links[l.Name()] = l
	return nil
}

func (m mockLinkStorage) Load(name internal.LinkName) (internal.Link, error) {
	l, exists := m.links[name]
	if !exists {
		err := &internal.NotFound{Type: "link", Ident: name.Unwrap()}
		return internal.Link{}, err
	}
	return l, nil
}

func createLink(
	t *testing.T,
	linkName, orchestrationName string,
	requiredOnly bool,
) internal.Link {
	name, err := internal.NewLinkName(linkName)
	assert.NilError(t, err, "create name for link %s", linkName)

	sourceStage, err := internal.NewStageName("source")
	assert.NilError(t, err, "create source stage for link %s", linkName)
	sourceFieldOpt := internal.NewEmptyMessageField()
	if !requiredOnly {
		sourceField := internal.NewMessageField("source-field")
		sourceFieldOpt = internal.NewPresentMessageField(sourceField)
	}
	sourceEndpoint := internal.NewLinkEndpoint(sourceStage, sourceFieldOpt)

	targetStage, err := internal.NewStageName("target")
	assert.NilError(t, err, "create target stage for link %s", linkName)
	targetFieldOpt := internal.NewEmptyMessageField()
	if !requiredOnly {
		targetField := internal.NewMessageField("target-field")
		targetFieldOpt = internal.NewPresentMessageField(targetField)
	}
	targetEndpoint := internal.NewLinkEndpoint(targetStage, targetFieldOpt)

	orchestration, err := internal.NewOrchestrationName(orchestrationName)
	assert.NilError(t, err, "create orchestration for link %s", linkName)

	return internal.NewLink(name, sourceEndpoint, targetEndpoint, orchestration)
}

func assertEqualLink(t *testing.T, expected, actual internal.Link) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assertEqualEndpoint(t, expected.Source(), actual.Source())
	assertEqualEndpoint(t, expected.Target(), actual.Target())
	assert.Equal(t, expected.Orchestration(), actual.Orchestration())
}

func assertEqualEndpoint(t *testing.T, expected, actual internal.LinkEndpoint) {
	assert.Equal(t, expected.Stage().Unwrap(), actual.Stage().Unwrap())
	assert.Equal(t, expected.Field().Present(), actual.Field().Present())
	if expected.Field().Present() {
		assert.Equal(t, expected.Field().Unwrap(), actual.Field().Unwrap())
	}
}
