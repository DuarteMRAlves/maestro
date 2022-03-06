package create

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

func createEmptyOrchestration(t *testing.T, orchName string) Orchestration {
	name, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create name for orchestration %s", orchName)
	return NewOrchestration(name, []domain.StageName{}, []domain.LinkName{})
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

func createStage(
	t *testing.T,
	stageName, orchName string,
	requiredOnly bool,
) Stage {
	name, err := domain.NewStageName(stageName)
	assert.NilError(t, err, "create name for stage some-name")
	address, err := domain.NewAddress("some-address")
	assert.NilError(t, err, "create address for stage some-name")
	orchestration, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create orchestration for stage some-name")
	serviceOpt := domain.NewEmptyService()
	methodOpt := domain.NewEmptyMethod()
	if !requiredOnly {
		service, err := domain.NewService("some-service")
		assert.NilError(t, err, "create service for stage some-name")
		serviceOpt = domain.NewPresentService(service)
		method, err := domain.NewMethod("some-method")
		assert.NilError(t, err, "create method for stage some-name")
		methodOpt = domain.NewPresentMethod(method)
	}
	ctx := domain.NewMethodContext(address, serviceOpt, methodOpt)
	return NewStage(name, ctx, orchestration)
}
