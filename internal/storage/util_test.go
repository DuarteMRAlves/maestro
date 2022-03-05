package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

func createAsset(
	t *testing.T,
	assetName string,
	requiredOnly bool,
) domain.Asset {
	name, err := domain.NewAssetName(assetName)
	assert.NilError(t, err, "create name for asset %s", assetName)
	imgOpt := domain.NewEmptyImage()
	if !requiredOnly {
		img, err := domain.NewImage("some-image")
		assert.NilError(t, err, "create image for asset %s", assetName)
		imgOpt = domain.NewPresentImage(img)
	}
	return domain.NewAsset(name, imgOpt)
}

func createStage(t *testing.T, name string, requiredOnly bool) domain.Stage {
	var (
		serviceOpt = domain.NewEmptyService()
		methodOpt  = domain.NewEmptyMethod()
	)
	stageName, err := domain.NewStageName(name)
	assert.NilError(t, err, "create name for stage %s", name)

	addr, err := domain.NewAddress(fmt.Sprintf("address-%s", name))
	assert.NilError(t, err, "create address for stage %s", name)

	if !requiredOnly {
		service, err := domain.NewService(fmt.Sprintf("service-%s", name))
		assert.NilError(t, err, "create service for stage %s", name)
		serviceOpt = domain.NewPresentService(service)
		method, err := domain.NewMethod(fmt.Sprintf("method-%s", name))
		assert.NilError(t, err, "create method for stage %s", name)
		methodOpt = domain.NewPresentMethod(method)
	}

	ctx := domain.NewMethodContext(addr, serviceOpt, methodOpt)
	return domain.NewStage(stageName, ctx)
}

func createLink(
	t *testing.T,
	linkName, sourceName, targetName string,
	requiredOnly bool,
) domain.Link {
	name, err := domain.NewLinkName(linkName)
	assert.NilError(t, err, "create name for link %s", linkName)

	source := createStage(t, sourceName, requiredOnly)
	target := createStage(t, targetName, requiredOnly)
	sourceFieldOpt := domain.NewEmptyMessageField()
	targetFieldOpt := domain.NewEmptyMessageField()

	if !requiredOnly {
		sourceField, err := domain.NewMessageField("source-field")
		assert.NilError(t, err, "create source field for link %s", linkName)
		sourceFieldOpt = domain.NewPresentMessageField(sourceField)
		targetField, err := domain.NewMessageField("target-field")
		assert.NilError(t, err, "create target field for link %s", linkName)
		targetFieldOpt = domain.NewPresentMessageField(targetField)
	}

	sourceEndpoint := domain.NewLinkEndpoint(source, sourceFieldOpt)
	targetEndpoint := domain.NewLinkEndpoint(target, targetFieldOpt)

	return domain.NewLink(name, sourceEndpoint, targetEndpoint)
}

func createOrchestration(
	t *testing.T,
	orchName string,
	stageNames, linkNames []string,
	requiredOnly bool,
) domain.Orchestration {
	// The ith link is from the ith stage to the (i+1)th stage
	assert.Equal(t, len(stageNames), len(linkNames)+1)
	name, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create name for orchestration %s", orchName)

	stages := make([]domain.Stage, 0, len(stageNames))
	for _, n := range stageNames {
		stages = append(stages, createStage(t, n, requiredOnly))
	}
	links := make([]domain.Link, 0, len(linkNames))
	for i, n := range linkNames {
		l := createLink(t, n, stageNames[i], stageNames[i+1], requiredOnly)
		links = append(links, l)
	}

	return domain.NewOrchestration(name, stages, links)
}

func assertEqualStage(t *testing.T, expected, actual domain.Stage) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assert.Equal(
		t,
		expected.MethodContext().Address().Unwrap(),
		actual.MethodContext().Address().Unwrap(),
	)
	if expected.MethodContext().Service().Present() {
		assert.Equal(
			t,
			expected.MethodContext().Service().Unwrap().Unwrap(),
			actual.MethodContext().Service().Unwrap().Unwrap(),
		)
	}
	if expected.MethodContext().Method().Present() {
		assert.Equal(
			t,
			expected.MethodContext().Method().Unwrap().Unwrap(),
			actual.MethodContext().Method().Unwrap().Unwrap(),
		)
	}
}

func assertEqualLink(t *testing.T, expected, actual domain.Link) {
	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
	assertEqualLinkEndpoint(t, expected.Source(), actual.Source())
	assertEqualLinkEndpoint(t, expected.Target(), actual.Target())
}

func assertEqualLinkEndpoint(
	t *testing.T,
	expected, actual domain.LinkEndpoint,
) {
	assertEqualStage(t, expected.Stage(), actual.Stage())
	assert.Equal(t, expected.Field().Present(), actual.Field().Present())
	if expected.Field().Present() {
		assert.Equal(t, expected.Field().Unwrap(), actual.Field().Unwrap())
	}
}
