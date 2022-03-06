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
			loadOrchestration: createEmptyOrchestration(t, "orchestration"),
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
			loadOrchestration: createEmptyOrchestration(t, "orchestration"),
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

				existsOrchestrationCount := 0
				loadOrchestrationCount := 0
				saveOrchestrationCount := 0

				existsStage := existsStageFn(
					test.expStage.Name(),
					&existsStageCount,
				)
				saveStage := saveStageFn(t, test.expStage, &saveStageCount)

				existsOrchestration := existsOrchestrationFn(
					test.loadOrchestration.Name(),
					&existsOrchestrationCount,
					0,
				)
				loadOrchestration := loadOrchestrationFn(
					t,
					test.loadOrchestration,
					&loadOrchestrationCount,
				)
				saveOrchestration := saveOrchestrationFn(
					t,
					test.expOrchestration,
					&saveOrchestrationCount,
				)

				createFn := CreateStage(
					existsStage,
					saveStage,
					existsOrchestration,
					loadOrchestration,
					saveOrchestration,
				)

				res := createFn(test.req)
				assert.Assert(t, !res.Err.Present())
				assert.Equal(t, existsStageCount, 1)
				assert.Equal(t, saveStageCount, 1)
				assert.Equal(t, existsOrchestrationCount, 1)
				assert.Equal(t, loadOrchestrationCount, 1)
				assert.Equal(t, saveOrchestrationCount, 1)
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
	storedOrchestration := createEmptyOrchestration(t, "orchestration")
	expOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"some-name"},
		[]string{},
	)
	existsStageCount := 0
	saveStageCount := 0

	existsOrchestrationCount := 0
	loadOrchestrationCount := 0
	saveOrchestrationCount := 0

	existsStage := existsStageFn(expStage.Name(), &existsStageCount)
	saveStage := saveStageFn(t, expStage, &saveStageCount)

	existsOrchestration := existsOrchestrationFn(
		storedOrchestration.Name(),
		&existsOrchestrationCount,
		0,
	)
	loadOrchestration := loadOrchestrationFn(
		t,
		storedOrchestration,
		&loadOrchestrationCount,
	)
	saveOrchestration := saveOrchestrationFn(
		t,
		expOrchestration,
		&saveOrchestrationCount,
	)

	createFn := CreateStage(
		existsStage,
		saveStage,
		existsOrchestration,
		loadOrchestration,
		saveOrchestration,
	)

	res := createFn(req)
	assert.Assert(t, !res.Err.Present())
	assert.Equal(t, existsStageCount, 1)
	assert.Equal(t, saveStageCount, 1)
	assert.Equal(t, existsOrchestrationCount, 1)
	assert.Equal(t, loadOrchestrationCount, 1)
	assert.Equal(t, saveOrchestrationCount, 1)

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
	// Should not run due to order of verifications
	assert.Equal(t, existsOrchestrationCount, 1)
	// Should not run due to order of verifications
	assert.Equal(t, loadOrchestrationCount, 1)
	// Orchestration Should not have been updated
	assert.Equal(t, saveOrchestrationCount, 1)
}
