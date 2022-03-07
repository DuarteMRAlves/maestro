package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestration(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createEmptyOrchestration(t, "some-name")
	existsCallCount := 0
	saveCallCount := 0
	existsFn := existsOrchestrationFn(expected.Name(), &existsCallCount, 1)
	saveFn := saveOrchestrationFn(t, expected, &saveCallCount)
	createFn := CreateOrchestration(existsFn, saveFn)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsCallCount, 1)
	assert.Equal(t, saveCallCount, 1)
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	req := OrchestrationRequest{Name: "some-name"}
	expected := createEmptyOrchestration(t, "some-name")
	existsCallCount := 0
	saveCallCount := 0
	existsFn := existsOrchestrationFn(expected.Name(), &existsCallCount, 1)
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
