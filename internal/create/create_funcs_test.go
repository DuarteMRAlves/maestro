package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"gotest.tools/v3/assert"
	"testing"
)

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
