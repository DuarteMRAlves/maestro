package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func existsStageFn(
	expected domain.StageName,
	callCount *int,
	threshold int,
) ExistsStage {
	return func(name domain.StageName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > threshold)
	}
}

func saveStageFn(t *testing.T, expected Stage, callCount *int) SaveStage {
	return func(s Stage) StageResult {
		*callCount++
		assertEqualStage(t, expected, s)
		return SomeStage(s)
	}
}

func assertEqualStage(t *testing.T, expected Stage, actual Stage) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assert.Equal(
		t,
		expected.Orchestration().Unwrap(),
		actual.Orchestration().Unwrap(),
	)
	assertEqualMethodContext(
		t,
		expected.MethodContext(),
		actual.MethodContext(),
	)
}

func assertEqualMethodContext(
	t *testing.T,
	expected domain.MethodContext,
	actual domain.MethodContext,
) {
	assert.Equal(t, expected.Address().Unwrap(), actual.Address().Unwrap())
	assert.Equal(t, expected.Service().Present(), actual.Service().Present())
	if expected.Service().Present() {
		assert.Equal(t, expected.Service().Unwrap(), actual.Service().Unwrap())
	}
	assert.Equal(t, expected.Method().Present(), actual.Method().Present())
	if expected.Method().Present() {
		assert.Equal(t, expected.Method().Present(), actual.Method().Present())
	}
}
