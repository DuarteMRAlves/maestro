package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestration(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Create(storage)

	err := createFn(req)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(storage.orchs))

	o, exists := storage.orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)
}

func TestCreateOrchestration_Err(t *testing.T) {
	tests := []struct {
		name    string
		req     OrchestrationRequest
		isError error
	}{
		{
			name:    "empty name",
			req:     OrchestrationRequest{Name: ""},
			isError: EmptyOrchestrationName,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				storage := mockOrchestrationStorage{
					orchs: map[internal.OrchestrationName]internal.Orchestration{},
				}

				createFn := Create(storage)
				err := createFn(test.req)
				assert.Assert(t, err != nil)
				assert.Assert(t, errors.Is(err, test.isError))

				assert.Equal(t, 0, len(storage.orchs))
			},
		)
	}
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Create(storage)

	err := createFn(req)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(storage.orchs))

	o, exists := storage.orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)

	err = createFn(req)
	assert.Assert(t, err != nil)
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "orchestration", alreadyExists.Type)
	assert.Equal(t, req.Name, alreadyExists.Ident)
	assert.Equal(t, 1, len(storage.orchs))

	o, exists = storage.orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)
}

type mockOrchestrationStorage struct {
	orchs map[internal.OrchestrationName]internal.Orchestration
}

func (m mockOrchestrationStorage) Save(o internal.Orchestration) error {
	m.orchs[o.Name()] = o
	return nil
}

func (m mockOrchestrationStorage) Load(name internal.OrchestrationName) (
	internal.Orchestration,
	error,
) {
	o, exists := m.orchs[name]
	if !exists {
		err := &internal.NotFound{Type: "orchestration", Ident: name.Unwrap()}
		return internal.Orchestration{}, err
	}
	return o, nil
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
