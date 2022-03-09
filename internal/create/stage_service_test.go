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
				existsStageCount := 0
				saveStageCount := 0

				existsStage := existsStageFn(
					test.expStage.Name(),
					&existsStageCount,
					1,
				)
				saveStage := saveStageFn(t, test.expStage, &saveStageCount)

				storage := mockOrchestrationStorage{
					orchs: map[domain.OrchestrationName]Orchestration{
						test.loadOrchestration.Name(): test.loadOrchestration,
					},
				}

				createFn := CreateStage(existsStage, saveStage, storage)

				res := createFn(test.req)
				assert.Assert(t, !res.Err.Present())
				assert.Equal(t, existsStageCount, 1)
				assert.Equal(t, saveStageCount, 1)

				assert.Equal(t, 1, len(storage.orchs))
				o, exists := storage.orchs[test.expOrchestration.Name()]
				assert.Assert(t, exists)
				assertEqualOrchestration(t, test.expOrchestration, o)
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
	existsStageCount := 0
	saveStageCount := 0

	existsStage := existsStageFn(expStage.Name(), &existsStageCount, 1)
	saveStage := saveStageFn(t, expStage, &saveStageCount)

	storage := mockOrchestrationStorage{
		orchs: map[domain.OrchestrationName]Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := CreateStage(existsStage, saveStage, storage)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsStageCount, 1)
	assert.Equal(t, saveStageCount, 1)

	assert.Equal(t, 1, len(storage.orchs))
	o, exists := storage.orchs[expOrchestration.Name()]
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
	assert.Equal(t, existsStageCount, 2)
	// Stage should not be saved
	assert.Equal(t, saveStageCount, 1)

	assert.Equal(t, 1, len(storage.orchs))
	o, exists = storage.orchs[expOrchestration.Name()]
	assert.Assert(t, exists)
	assertEqualOrchestration(t, expOrchestration, o)
}
