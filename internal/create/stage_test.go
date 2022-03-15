package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateStage(t *testing.T) {
	tests := []struct {
		name              string
		stageName         internal.StageName
		methodCtx         internal.MethodContext
		orchName          internal.OrchestrationName
		expStage          internal.Stage
		loadOrchestration internal.Orchestration
		expOrchestration  internal.Orchestration
	}{
		{
			name:      "required fields",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewEmptyService(),
				internal.NewEmptyMethod(),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expStage: createStage(
				t,
				"some-name",
				"orchestration",
				true,
			),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
			expOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"some-name"},
				[]string{},
			),
		},
		{
			name:      "all fields",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewPresentService(internal.NewService("some-service")),
				internal.NewPresentMethod(internal.NewMethod("some-method")),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expStage: createStage(
				t,
				"some-name",
				"orchestration",
				false,
			),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
			expOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"some-name"},
				[]string{},
			),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)

				err := createFn(test.stageName, test.methodCtx, test.orchName)
				assert.NilError(t, err)

				assert.Equal(t, 1, len(stageStore.Stages))
				s, exists := stageStore.Stages[test.expStage.Name()]
				assert.Assert(t, exists)
				assertEqualStage(t, test.expStage, s)

				assert.Equal(t, 1, len(orchStore.Orchs))
				o, exists := orchStore.Orchs[test.expOrchestration.Name()]
				assert.Assert(t, exists)
				assertEqualOrchestration(t, test.expOrchestration, o)
			},
		)
	}
}

func TestCreateStage_Err(t *testing.T) {
	tests := []struct {
		name              string
		stageName         internal.StageName
		methodCtx         internal.MethodContext
		orchName          internal.OrchestrationName
		isError           error
		loadOrchestration internal.Orchestration
	}{
		{
			name:      "empty name",
			stageName: createStageName(t, ""),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewEmptyService(),
				internal.NewEmptyMethod(),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			isError:  EmptyStageName,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name:      "empty address",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress(""),
				internal.NewEmptyService(),
				internal.NewEmptyMethod(),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			isError:  EmptyAddress,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name:      "empty service",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewPresentService(internal.NewService("")),
				internal.NewEmptyMethod(),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			isError:  EmptyService,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name:      "empty method",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewEmptyService(),
				internal.NewPresentMethod(internal.NewMethod("")),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			isError:  EmptyMethod,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name:      "empty orchestration",
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewEmptyService(),
				internal.NewEmptyMethod(),
			),
			orchName: createOrchestrationName(t, ""),
			isError:  EmptyOrchestrationName,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)
				err := createFn(test.stageName, test.methodCtx, test.orchName)
				assert.Assert(t, err != nil)
				assert.Assert(t, errors.Is(err, test.isError))

				assert.Equal(t, 0, len(stageStore.Stages))
			},
		)
	}
}

func TestCreateStage_AlreadyExists(t *testing.T) {
	stageName := createStageName(t, "some-name")
	methodCtx := internal.NewMethodContext(
		internal.NewAddress("some-address"),
		internal.NewPresentService(internal.NewService("some-service")),
		internal.NewPresentMethod(internal.NewMethod("some-method")),
	)
	orchName := createOrchestrationName(t, "orchestration")

	expStage := createStage(t, "some-name", "orchestration", false)
	storedOrchestration := createOrchestration(t, "orchestration", nil, nil)
	expOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"some-name"},
		[]string{},
	)

	stageStore := mock.StageStorage{
		Stages: map[internal.StageName]internal.Stage{},
	}

	orchStore := mock.OrchestrationStorage{
		Orchs: map[internal.OrchestrationName]internal.Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := Stage(stageStore, orchStore)

	err := createFn(stageName, methodCtx, orchName)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(stageStore.Stages))
	s, exists := stageStore.Stages[expStage.Name()]
	assert.Assert(t, exists)
	assertEqualStage(t, expStage, s)

	assert.Equal(t, 1, len(orchStore.Orchs))
	o, exists := orchStore.Orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)

	err = createFn(stageName, methodCtx, orchName)
	assert.Assert(t, err != nil)
	var alreadyExists *internal.AlreadyExists
	assert.Assert(t, errors.As(err, &alreadyExists))
	assert.Equal(t, "stage", alreadyExists.Type)
	assert.Equal(t, stageName.Unwrap(), alreadyExists.Ident)

	assert.Equal(t, 1, len(stageStore.Stages))
	s, exists = stageStore.Stages[expStage.Name()]
	assert.Assert(t, exists)
	assertEqualStage(t, expStage, s)

	assert.Equal(t, 1, len(orchStore.Orchs))
	o, exists = orchStore.Orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

func createStageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	assert.NilError(t, err, "create stage name %s", name)
	return stageName
}

func createStage(
	t *testing.T,
	stageName, orchName string,
	requiredOnly bool,
) internal.Stage {
	name, err := internal.NewStageName(stageName)
	assert.NilError(t, err, "create name for stage %s", stageName)
	address := internal.NewAddress("some-address")
	orchestration, err := internal.NewOrchestrationName(orchName)
	assert.NilError(t, err, "create orchestration for stage %s", stageName)
	serviceOpt := internal.NewEmptyService()
	methodOpt := internal.NewEmptyMethod()
	if !requiredOnly {
		serviceOpt = internal.NewPresentService(internal.NewService("some-service"))
		method := internal.NewMethod("some-method")
		methodOpt = internal.NewPresentMethod(method)
	}
	ctx := internal.NewMethodContext(address, serviceOpt, methodOpt)
	return internal.NewStage(name, ctx, orchestration)
}

func assertEqualStage(
	t *testing.T,
	expected internal.Stage,
	actual internal.Stage,
) {
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
	expected internal.MethodContext,
	actual internal.MethodContext,
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
