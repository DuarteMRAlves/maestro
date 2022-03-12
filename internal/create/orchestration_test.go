package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestration(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Create(storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())

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
				res := createFn(test.req)
				assert.Assert(t, res.Err.Present())

				assert.Equal(t, 0, len(storage.orchs))

				err := res.Err.Unwrap()
				assert.Assert(t, errors.Is(err, test.isError))
			},
		)
	}
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Create(storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, 1, len(storage.orchs))

	o, exists := storage.orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("orchestration '%v' already exists", req.Name),
	)
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
