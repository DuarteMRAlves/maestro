package storage

import (
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

// func createStage(t *testing.T, name string, requiredOnly bool) domain.Stage {
// 	var (
// 		serviceOpt = domain.NewEmptyService()
// 		methodOpt  = domain.NewEmptyMethod()
// 	)
// 	stageName, err := domain.NewStageName(name)
// 	assert.NilError(t, err, "create name for stage %s", name)
//
// 	addr, err := domain.NewAddress(fmt.Sprintf("address-%s", name))
// 	assert.NilError(t, err, "create address for stage %s", name)
//
// 	if !requiredOnly {
// 		service, err := domain.NewService(fmt.Sprintf("service-%s", name))
// 		assert.NilError(t, err, "create service for stage %s", name)
// 		serviceOpt = domain.NewPresentService(service)
// 		method, err := domain.NewMethod(fmt.Sprintf("method-%s", name))
// 		assert.NilError(t, err, "create method for stage %s", name)
// 		methodOpt = domain.NewPresentMethod(method)
// 	}
//
// 	ctx := domain.NewMethodContext(addr, serviceOpt, methodOpt)
// 	return create.NewStage(stageName, ctx)
// }
//
//
// func createLink(
// 	t *testing.T,
// 	linkName, sourceName, targetName string,
// 	requiredOnly bool,
// ) execute.Link {
// 	name, err := domain.NewLinkName(linkName)
// 	assert.NilError(t, err, "create name for link %s", linkName)
//
// 	source := createStage(t, sourceName, requiredOnly)
// 	target := createStage(t, targetName, requiredOnly)
// 	sourceFieldOpt := domain.NewEmptyMessageField()
// 	targetFieldOpt := domain.NewEmptyMessageField()
//
// 	if !requiredOnly {
// 		sourceField, err := domain.NewMessageField("source-field")
// 		assert.NilError(t, err, "create source field for link %s", linkName)
// 		sourceFieldOpt = domain.NewPresentMessageField(sourceField)
// 		targetField, err := domain.NewMessageField("target-field")
// 		assert.NilError(t, err, "create target field for link %s", linkName)
// 		targetFieldOpt = domain.NewPresentMessageField(targetField)
// 	}
//
// 	sourceEndpoint := execute.NewLinkEndpoint(source, sourceFieldOpt)
// 	targetEndpoint := execute.NewLinkEndpoint(target, targetFieldOpt)
//
// 	return execute.NewLink(name, sourceEndpoint, targetEndpoint)
// }
//
// func createCreateOrchestration(
// 	t *testing.T,
// 	orchName string,
// 	stageNames, linkNames []string,
// ) create.Orchestration {
// 	// The ith link is from the ith stage to the (i+1)th stage or both are empty.
// 	assert.Assert(
// 		t,
// 		len(stageNames) == len(linkNames)+1 ||
// 			(len(stageNames) == 0 && len(linkNames) == 0),
// 	)
// 	name, err := domain.NewOrchestrationName(orchName)
// 	assert.NilError(t, err, "create name for orchestration %s", orchName)
//
// 	stages := make([]domain.StageName, 0, len(stageNames))
// 	for _, n := range stageNames {
// 		stageName, err := domain.NewStageName(n)
// 		assert.NilError(t, err, "create stage for orchestration %s", orchName)
// 		stages = append(stages, stageName)
// 	}
//
// 	links := make([]domain.LinkName, 0, len(linkNames))
// 	for _, n := range linkNames {
// 		linkName, err := domain.NewLinkName(n)
// 		assert.NilError(t, err, "create link for orchestration %s", orchName)
// 		links = append(links, linkName)
// 	}
// 	return create.NewOrchestration(name, stages, links)
// }
//
// func createExecuteOrchestration(
// 	t *testing.T,
// 	orchName string,
// 	stageNames, linkNames []string,
// 	requiredOnly bool,
// ) execute.Orchestration {
// 	// The ith link is from the ith stage to the (i+1)th stage
// 	assert.Equal(t, len(stageNames), len(linkNames)+1)
// 	name, err := domain.NewOrchestrationName(orchName)
// 	assert.NilError(t, err, "create name for orchestration %s", orchName)
//
// 	stages := make([]domain.Stage, 0, len(stageNames))
// 	for _, n := range stageNames {
// 		stages = append(stages, createStage(t, n, requiredOnly))
// 	}
// 	links := make([]execute.Link, 0, len(linkNames))
// 	for i, n := range linkNames {
// 		l := createLink(t, n, stageNames[i], stageNames[i+1], requiredOnly)
// 		links = append(links, l)
// 	}
//
// 	return execute.NewOrchestration(name, stages, links)
// }
//
// func assertEqualStage(t *testing.T, expected, actual domain.Stage) {
// 	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
// 	assert.Equal(
// 		t,
// 		expected.MethodContext().Address().Unwrap(),
// 		actual.MethodContext().Address().Unwrap(),
// 	)
// 	if expected.MethodContext().Service().Present() {
// 		assert.Equal(
// 			t,
// 			expected.MethodContext().Service().Unwrap().Unwrap(),
// 			actual.MethodContext().Service().Unwrap().Unwrap(),
// 		)
// 	}
// 	if expected.MethodContext().Method().Present() {
// 		assert.Equal(
// 			t,
// 			expected.MethodContext().Method().Unwrap().Unwrap(),
// 			actual.MethodContext().Method().Unwrap().Unwrap(),
// 		)
// 	}
// }
//
// func assertEqualLink(t *testing.T, expected, actual execute.Link) {
// 	assert.Equal(t, expected.Name().Unwrap(), actual.Name().Unwrap())
// 	assertEqualLinkEndpoint(t, expected.Source(), actual.Source())
// 	assertEqualLinkEndpoint(t, expected.Target(), actual.Target())
// }
//
// func assertEqualLinkEndpoint(
// 	t *testing.T,
// 	expected, actual execute.LinkEndpoint,
// ) {
// 	assertEqualStage(t, expected.Stage(), actual.Stage())
// 	assert.Equal(t, expected.Field().Present(), actual.Field().Present())
// 	if expected.Field().Present() {
// 		assert.Equal(t, expected.Field().Unwrap(), actual.Field().Unwrap())
// 	}
// }
