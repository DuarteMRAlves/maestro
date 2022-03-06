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
	expected := createOrchestration(t, "some-name")
	existsCallCount := 0
	saveCallCount := 0
	existsFn := existsOrchestrationFn(expected.Name(), &existsCallCount)
	saveFn := saveOrchestrationFn(t, expected, &saveCallCount)
	createFn := CreateOrchestration(existsFn, saveFn)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsCallCount, 1)
	assert.Equal(t, saveCallCount, 1)
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createOrchestration(t, "some-name")
	existsCallCount := 0
	saveCallCount := 0
	existsFn := existsOrchestrationFn(expected.Name(), &existsCallCount)
	saveFn := saveOrchestrationFn(t, expected, &saveCallCount)
	createFn := CreateOrchestration(existsFn, saveFn)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsCallCount, 1)
	assert.Equal(t, saveCallCount, 1)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("orchestration '%v' already exists", req.Name),
	)
	assert.Equal(t, existsCallCount, 2)
	// Should not call save
	assert.Equal(t, saveCallCount, 1)
}

func existsOrchestrationFn(
	expected domain.OrchestrationName,
	callCount *int,
) ExistsOrchestration {
	return func(name domain.OrchestrationName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > 1)
	}
}

func saveOrchestrationFn(
	t *testing.T,
	expected domain.Orchestration,
	callCount *int,
) SaveOrchestration {
	return func(actual domain.Orchestration) domain.OrchestrationResult {
		*callCount++
		assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
		return domain.SomeOrchestration(actual)
	}
}

func createOrchestration(t *testing.T, orchName string) domain.Orchestration {
	name, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create name for orchestration %s", orchName)
	return domain.NewOrchestration(name, []domain.Stage{}, []domain.Link{})
}
