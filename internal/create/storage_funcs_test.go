package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func existsOrchestrationFn(
	expected domain.OrchestrationName,
	callCount *int,
	threshold int,
) ExistsOrchestration {
	return func(name domain.OrchestrationName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > threshold)
	}
}

func loadOrchestrationFn(
	t *testing.T,
	expected Orchestration,
	callCount *int,
) LoadOrchestration {
	return func(name domain.OrchestrationName) OrchestrationResult {
		*callCount++
		assert.Equal(t, expected.Name().Unwrap(), name.Unwrap())
		return SomeOrchestration(expected)
	}
}

func saveOrchestrationFn(
	t *testing.T,
	expected Orchestration,
	callCount *int,
) SaveOrchestration {
	return func(actual Orchestration) OrchestrationResult {
		*callCount++
		assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
		return SomeOrchestration(actual)
	}
}

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

func existsLinkFn(
	expected domain.LinkName,
	callCount *int,
	threshold int,
) ExistsLink {
	return func(name domain.LinkName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > threshold)
	}
}

func saveLinkFn(t *testing.T, expected Link, callCount *int) SaveLink {
	return func(l Link) LinkResult {
		*callCount++
		assertEqualLink(t, expected, l)
		return SomeLink(l)
	}
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
