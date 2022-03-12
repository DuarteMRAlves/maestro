package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateStage(t *testing.T) {
	tests := []struct {
		name              string
		req               StageRequest
		expStage          internal.Stage
		loadOrchestration internal.Orchestration
		expOrchestration  internal.Orchestration
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
					stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mockOrchestrationStorage{
					orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)

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
		isError           error
		loadOrchestration internal.Orchestration
	}{
		{
			name:    "empty name",
			req:     StageRequest{Name: ""},
			isError: EmptyStageName,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name:    "empty address",
			req:     StageRequest{Name: "some-name", Address: ""},
			isError: EmptyAddress,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name: "empty service",
			req: StageRequest{
				Name:    "some-name",
				Address: "some-address",
				Service: domain.NewPresentString(""),
			},
			isError: EmptyService,
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
		},
		{
			name: "empty method",
			req: StageRequest{
				Name:    "some-name",
				Address: "some-address",
				Service: domain.NewEmptyString(),
				Method:  domain.NewPresentString(""),
			},
			isError: EmptyMethod,
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
					stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mockOrchestrationStorage{
					orchs: map[internal.OrchestrationName]internal.Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)
				res := createFn(test.req)
				assert.Assert(t, res.Err.Present())

				assert.Equal(t, 0, len(stageStore.stages))

				err := res.Err.Unwrap()
				assert.Assert(t, errors.Is(err, test.isError))
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
		stages: map[internal.StageName]internal.Stage{},
	}

	orchStore := mockOrchestrationStorage{
		orchs: map[internal.OrchestrationName]internal.Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := Stage(stageStore, orchStore)

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
	stages map[internal.StageName]internal.Stage
}

func (m mockStageStorage) Save(s internal.Stage) error {
	m.stages[s.Name()] = s
	return nil
}

func (m mockStageStorage) Load(name internal.StageName) (
	internal.Stage,
	error,
) {
	s, exists := m.stages[name]
	if !exists {
		err := &internal.NotFound{Type: "stage", Ident: name.Unwrap()}
		return internal.Stage{}, err
	}
	return s, nil
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
