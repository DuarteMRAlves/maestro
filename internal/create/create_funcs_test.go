package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

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
	assert.NilError(t, err, "create name for stage %s", stageName)
	address, err := domain.NewAddress("some-address")
	assert.NilError(t, err, "create address for stage %s", stageName)
	orchestration, err := domain.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create orchestration for stage %s", stageName)
	serviceOpt := domain.NewEmptyService()
	methodOpt := domain.NewEmptyMethod()
	if !requiredOnly {
		service, err := domain.NewService("some-service")
		assert.NilError(t, err, "create service for stage %", stageName)
		serviceOpt = domain.NewPresentService(service)
		method, err := domain.NewMethod("some-method")
		assert.NilError(t, err, "create method for stage %s", stageName)
		methodOpt = domain.NewPresentMethod(method)
	}
	ctx := domain.NewMethodContext(address, serviceOpt, methodOpt)
	return NewStage(name, ctx, orchestration)
}

func createLink(
	t *testing.T,
	linkName, orchestrationName string,
	requiredOnly bool,
) Link {
	name, err := domain.NewLinkName(linkName)
	assert.NilError(t, err, "create name for link %s", linkName)

	sourceStage, err := domain.NewStageName("source")
	assert.NilError(t, err, "create source stage for link %s", linkName)
	sourceFieldOpt := domain.NewEmptyMessageField()
	if !requiredOnly {
		sourceField, err := domain.NewMessageField("source-field")
		assert.NilError(t, err, "create source field for link %s", linkName)
		sourceFieldOpt = domain.NewPresentMessageField(sourceField)
	}
	sourceEndpoint := NewLinkEndpoint(sourceStage, sourceFieldOpt)

	targetStage, err := domain.NewStageName("target")
	assert.NilError(t, err, "create target stage for link %s", linkName)
	targetFieldOpt := domain.NewEmptyMessageField()
	if !requiredOnly {
		targetField, err := domain.NewMessageField("target-field")
		assert.NilError(t, err, "create target field for link %s", linkName)
		targetFieldOpt = domain.NewPresentMessageField(targetField)
	}
	targetEndpoint := NewLinkEndpoint(targetStage, targetFieldOpt)

	orchestration, err := domain.NewOrchestrationName(orchestrationName)
	assert.NilError(t, err, "create orchestration for link %s", linkName)

	return NewLink(name, sourceEndpoint, targetEndpoint, orchestration)
}
