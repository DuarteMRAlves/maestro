package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateStage(t *testing.T) {
	tests := []struct {
		name              string
		req               StageRequest
		expStage          Stage
		loadOrchestration Orchestration
		expOrchestration  Orchestration
	}{
		{
			name: "required fields",
			req: StageRequest{
				Name:          "some-name",
				Address:       "some-address",
				Service:       domain.NewEmptyString(),
				Method:        domain.NewEmptyString(),
				Orchestration: "orchestration",
			},
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
			name: "all fields",
			req: StageRequest{
				Name:          "some-name",
				Address:       "some-address",
				Service:       domain.NewPresentString("some-service"),
				Method:        domain.NewPresentString("some-method"),
				Orchestration: "orchestration",
			},
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
				stageStore := mockStageStorage{
					stages: map[domain.StageName]Stage{},
				}

				orchStore := mockOrchestrationStorage{
					orchs: map[domain.OrchestrationName]Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateStage(stageStore, orchStore)

				res := createFn(test.req)
				assert.Assert(t, !res.Err.Present())

				assert.Equal(t, 1, len(stageStore.stages))
				s, exists := stageStore.stages[test.expStage.Name()]
				assert.Assert(t, exists)
				assertEqualStage(t, test.expStage, s)

				assert.Equal(t, 1, len(orchStore.orchs))
				o, exists := orchStore.orchs[test.expOrchestration.Name()]
				assert.Assert(t, exists)
				assertEqualOrchestration(t, test.expOrchestration, o)
			},
		)
	}
}

func TestCreateStage_Err(t *testing.T) {
	tests := []struct {
		name              string
		req               StageRequest
		assertErrTypeFn   func(error) bool
		expectedErrMsg    string
		loadOrchestration Orchestration
	}{
		{
			name:            "empty name",
			req:             StageRequest{Name: ""},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "empty stage name",
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
				stageStore := mockStageStorage{
					stages: map[domain.StageName]Stage{},
				}

				orchStore := mockOrchestrationStorage{
					orchs: map[domain.OrchestrationName]Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateStage(stageStore, orchStore)
				res := createFn(test.req)
				assert.Assert(t, res.Err.Present())

				assert.Equal(t, 0, len(stageStore.stages))

				err := res.Err.Unwrap()
				assert.Assert(t, test.assertErrTypeFn(err))
				assert.Error(t, err, test.expectedErrMsg)
			},
		)
	}
}

func TestCreateStage_AlreadyExists(t *testing.T) {
	req := StageRequest{
		Name:          "some-name",
		Address:       "some-address",
		Service:       domain.NewPresentString("some-service"),
		Method:        domain.NewPresentString("some-method"),
		Orchestration: "orchestration",
	}
	expStage := createStage(t, "some-name", "orchestration", false)
	storedOrchestration := createOrchestration(t, "orchestration", nil, nil)
	expOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"some-name"},
		[]string{},
	)

	stageStore := mockStageStorage{
		stages: map[domain.StageName]Stage{},
	}

	orchStore := mockOrchestrationStorage{
		orchs: map[domain.OrchestrationName]Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := CreateStage(stageStore, orchStore)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())

	assert.Equal(t, 1, len(stageStore.stages))
	s, exists := stageStore.stages[expStage.Name()]
	assert.Assert(t, exists)
	assertEqualStage(t, expStage, s)

	assert.Equal(t, 1, len(orchStore.orchs))
	o, exists := orchStore.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)

	res = createFn(req)
	assert.Assert(t, res.Err.Present())
	err := res.Err.Unwrap()
	assert.Assert(t, errdefs.IsAlreadyExists(err), "err type")
	assert.ErrorContains(
		t,
		err,
		fmt.Sprintf("stage '%v' already exists", req.Name),
	)

	assert.Equal(t, 1, len(stageStore.stages))
	s, exists = stageStore.stages[expStage.Name()]
	assert.Assert(t, exists)
	assertEqualStage(t, expStage, s)

	assert.Equal(t, 1, len(orchStore.orchs))
	o, exists = orchStore.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}

type mockStageStorage struct {
	stages map[domain.StageName]Stage
}

func (m mockStageStorage) Save(s Stage) StageResult {
	m.stages[s.Name()] = s
	return SomeStage(s)
}

func (m mockStageStorage) Load(name domain.StageName) StageResult {
	s, exists := m.stages[name]
	if !exists {
		err := errdefs.NotFoundWithMsg("stage not found: %s", name)
		return ErrStage(err)
	}
	return SomeStage(s)
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
