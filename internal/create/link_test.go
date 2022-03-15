package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateLink(t *testing.T) {
	tests := []struct {
		name              string
		linkName          internal.LinkName
		source            internal.LinkEndpoint
		target            internal.LinkEndpoint
		orchName          internal.OrchestrationName
		expLink           internal.Link
		loadOrchestration internal.Orchestration
		expOrchestration  internal.Orchestration
		storedStages      []internal.Stage
	}{
		{
			name:     "required fields",
			linkName: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewEmptyMessageField(),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewEmptyMessageField(),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expLink: createLink(
				t,
				"some-name",
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
				createStage(t, "source", true),
				createStage(t, "target", true),
			},
		},
		{
			name:     "all fields",
			linkName: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewPresentMessageField(internal.NewMessageField("source-field")),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewPresentMessageField(internal.NewMessageField("target-field")),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expLink:  createLink(t, "some-name", false),
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
				createStage(t, "source", false),
				createStage(t, "target", false),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}
				for _, s := range test.storedStages {
					stageStore.Stages[s.Name()] = s
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := Link(linkStore, stageStore, orchStore)
				err := createFn(
					test.linkName,
					test.source,
					test.target,
					test.orchName,
				)

				assert.NilError(t, err)

				assert.Equal(t, 1, len(linkStore.Links))
				l, exists := linkStore.Links[test.expLink.Name()]
				assert.Assert(t, exists)
				assertEqualLink(t, test.expLink, l)

				assert.Equal(t, 1, len(orchStore.Orchs))
				o, exists := orchStore.Orchs[test.expOrchestration.Name()]
				assert.Assert(t, exists)
				assertEqualOrchestration(t, test.expOrchestration, o)
			},
		)
	}
}

func TestCreateLink_AlreadyExists(t *testing.T) {
	linkName := createLinkName(t, "some-name")
	source := internal.NewLinkEndpoint(
		createStageName(t, "source"),
		internal.NewPresentMessageField(internal.NewMessageField("source-field")),
	)
	target := internal.NewLinkEndpoint(
		createStageName(t, "target"),
		internal.NewPresentMessageField(internal.NewMessageField("target-field")),
	)
	orchName := createOrchestrationName(t, "orchestration")
	expLink := createLink(t, "some-name", false)
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
		createStage(t, "source", false),
		createStage(t, "target", false),
	}

	linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

	stageStore := mock.StageStorage{
		Stages: map[internal.StageName]internal.Stage{},
	}
	for _, s := range storedStages {
		stageStore.Stages[s.Name()] = s
	}

	orchStore := mock.OrchestrationStorage{
		Orchs: map[internal.OrchestrationName]internal.Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := Link(linkStore, stageStore, orchStore)
	err := createFn(linkName, source, target, orchName)

	assert.NilError(t, err)

	assert.Equal(t, 1, len(linkStore.Links))
	l, exists := linkStore.Links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)

	assert.Equal(t, 1, len(orchStore.Orchs))
	o, exists := orchStore.Orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)

	err = createFn(linkName, source, target, orchName)

	assert.Assert(t, err != nil)
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "link", alreadyExists.Type)
	assert.Equal(t, linkName.Unwrap(), alreadyExists.Ident)
	assert.Equal(t, 1, len(linkStore.Links))
	l, exists = linkStore.Links[expLink.Name()]
	assert.Assert(t, exists)
	assertEqualLink(t, expLink, l)

	assert.Equal(t, 1, len(orchStore.Orchs))
	o, exists = orchStore.Orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

func createLinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	assert.NilError(t, err, "create link name %s", name)
	return linkName
}

func createLink(
	t *testing.T,
	linkName string,
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

	return internal.NewLink(name, sourceEndpoint, targetEndpoint)
}

func assertEqualLink(t *testing.T, expected, actual internal.Link) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assertEqualEndpoint(t, expected.Source(), actual.Source())
	assertEqualEndpoint(t, expected.Target(), actual.Target())
}

func assertEqualEndpoint(t *testing.T, expected, actual internal.LinkEndpoint) {
	assert.Equal(t, expected.Stage().Unwrap(), actual.Stage().Unwrap())
	assert.Equal(t, expected.Field().Present(), actual.Field().Present())
	if expected.Field().Present() {
		assert.Equal(t, expected.Field().Unwrap(), actual.Field().Unwrap())
	}
}
