package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestration(t *testing.T) {
	name := "some-name"
	orchName := createOrchestrationName(t, name)
	expected := createOrchestration(t, name, nil, nil)
	storage := mock.OrchestrationStorage{Orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Orchestration(storage)

	err := createFn(orchName)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(storage.Orchs))

	o, exists := storage.Orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)
}

func TestCreateOrchestration_Err(t *testing.T) {
	tests := map[string]struct {
		name    internal.OrchestrationName
		isError error
	}{
		"empty name": {
			name:    createOrchestrationName(t, ""),
			isError: EmptyOrchestrationName,
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				storage := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{},
				}

				createFn := Orchestration(storage)
				err := createFn(tc.name)
				assert.Assert(t, err != nil)
				assert.Assert(t, errors.Is(err, tc.isError))

				assert.Equal(t, 0, len(storage.Orchs))
			},
		)
	}
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	name := "some-name"
	orchName := createOrchestrationName(t, name)
	expected := createOrchestration(t, name, nil, nil)
	storage := mock.OrchestrationStorage{Orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Orchestration(storage)

	err := createFn(orchName)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(storage.Orchs))

	o, exists := storage.Orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)

	err = createFn(orchName)
	assert.Assert(t, err != nil)
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "orchestration", alreadyExists.Type)
	assert.Equal(t, name, alreadyExists.Ident)
	assert.Equal(t, 1, len(storage.Orchs))

	o, exists = storage.Orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)
}

func createOrchestrationName(
	t *testing.T,
	orchName string,
) internal.OrchestrationName {
	name, err := internal.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create orchestration name %s", orchName)
	return name
}

func createOrchestration(
	t *testing.T,
	orchName string,
	stages, links []string,
) internal.Orchestration {
	name, err := internal.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create name for orchestration %s", orchName)
	stageNames := make([]internal.StageName, 0, len(stages))
	for _, s := range stages {
		sName, err := internal.NewStageName(s)
		assert.NilError(t, err, "create stage for orchestration %s", orchName)
		stageNames = append(stageNames, sName)
	}
	linkNames := make([]internal.LinkName, 0, len(links))
	for _, l := range links {
		lName, err := internal.NewLinkName(l)
		assert.NilError(t, err, "create link for orchestration %s", orchName)
		linkNames = append(linkNames, lName)
	}
	return internal.NewOrchestration(name, stageNames, linkNames)
}

func assertEqualOrchestration(
	t *testing.T,
	expected, actual internal.Orchestration,
) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())

	assert.Equal(t, len(expected.Stages()), len(actual.Stages()))
	for i, s := range expected.Stages() {
		assert.Equal(t, s.Unwrap(), actual.Stages()[i].Unwrap())
	}

	assert.Equal(t, len(expected.Links()), len(actual.Links()))
	for i, l := range expected.Links() {
		assert.Equal(t, l.Unwrap(), actual.Links()[i].Unwrap())
	}
}
