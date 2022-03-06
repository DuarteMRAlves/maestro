package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func existsAssetFn(expected domain.AssetName, callCount *int) ExistsAsset {
	return func(name domain.AssetName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > 1)
	}
}

func saveAssetFn(
	t *testing.T,
	expected domain.Asset,
	callCount *int,
) SaveAsset {
	return func(actual domain.Asset) domain.AssetResult {
		*callCount++
		assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
		assert.Equal(t, expected.Image().Present(), actual.Image().Present())
		if expected.Image().Present() {
			assert.Equal(t, expected.Image().Unwrap(), actual.Image().Unwrap())
		}
		return domain.SomeAsset(actual)
	}
}

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

func existsStageFn(expected domain.StageName, callCount *int) ExistsStage {
	return func(name domain.StageName) bool {
		*callCount++
		return expected.Unwrap() == name.Unwrap() && (*callCount > 1)
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
