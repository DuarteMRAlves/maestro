package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestration(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[domain.OrchestrationName]Orchestration{}}

	createFn := CreateOrchestration(storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())

	assert.Equal(t, 1, len(storage.orchs))

	o, exists := storage.orchs[expected.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expected, o)
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name", nil, nil)
	storage := mockOrchestrationStorage{orchs: map[domain.OrchestrationName]Orchestration{}}

	createFn := CreateOrchestration(storage)

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
	orchs map[domain.OrchestrationName]Orchestration
}

func (m mockOrchestrationStorage) Save(o Orchestration) OrchestrationResult {
	m.orchs[o.Name()] = o
	return SomeOrchestration(o)
}

func (m mockOrchestrationStorage) Load(name domain.OrchestrationName) OrchestrationResult {
	o, exists := m.orchs[name]
	if !exists {
		err := errdefs.NotFoundWithMsg("orchestration not found: %s", o.Name())
		return ErrOrchestration(err)
	}
	return SomeOrchestration(o)
}

func (m mockOrchestrationStorage) Verify(name domain.OrchestrationName) bool {
	_, exists := m.orchs[name]
	return exists
}

func createOrchestration(
	t *testing.T,
	orchName string,
	stages, links []string,
) Orchestration {
	name, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create name for orchestration %s", orchName)
	stageNames := make([]domain.StageName, 0, len(stages))
	for _, s := range stages {
		sName, err := domain.NewStageName(s)
		assert.NilError(t, err, "create stage for orchestration %s", orchName)
		stageNames = append(stageNames, sName)
	}
	linkNames := make([]domain.LinkName, 0, len(links))
	for _, l := range links {
		lName, err := domain.NewLinkName(l)
		assert.NilError(t, err, "create link for orchestration %s", orchName)
		linkNames = append(linkNames, lName)
	}
	return NewOrchestration(name, stageNames, linkNames)
}

func assertEqualOrchestration(t *testing.T, expected, actual Orchestration) {
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
