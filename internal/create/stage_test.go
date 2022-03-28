package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestCreateStage(t *testing.T) {
	tests := map[string]struct {
		name              internal.StageName
		methodCtx         internal.MethodContext
		orchName          internal.OrchestrationName
		expStage          internal.Stage
		loadOrchestration internal.Orchestration
		expOrch           internal.Orchestration
	}{
		"required fields": {
			name: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expStage: createStage(t, "some-name", true),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
			expOrch: createOrchestration(
				t,
				"orchestration",
				[]string{"some-name"},
				[]string{},
			),
		},
		"all fields": {
			name: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewService("some-service"),
				internal.NewMethod("some-method"),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expStage: createStage(t, "some-name", false),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				nil,
				nil,
			),
			expOrch: createOrchestration(
				t,
				"orchestration",
				[]string{"some-name"},
				[]string{},
			),
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						tc.loadOrchestration.Name(): tc.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)

				err := createFn(tc.name, tc.methodCtx, tc.orchName)
				if err != nil {
					t.Fatalf("create error: %s", err)
				}

				if diff := cmp.Diff(1, len(stageStore.Stages)); diff != "" {
					t.Fatalf("number of stages mismatch:\n%s", diff)
				}
				s, exists := stageStore.Stages[tc.expStage.Name()]
				if !exists {
					t.Fatalf("created stage does not exist in storage")
				}
				cmpStage(t, tc.expStage, s, "created stage")

				if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
					t.Fatalf("number of orchestrations mismatch:\n%s", diff)
				}
				o, exists := orchStore.Orchs[tc.expOrch.Name()]
				if !exists {
					t.Fatalf("updated orchestration does not exist in storage")
				}
				cmpOrchestration(t, tc.expOrch, o, "updated orchestration")
			},
		)
	}
}

func TestCreateStage_Err(t *testing.T) {
	tests := map[string]struct {
		stageName         internal.StageName
		methodCtx         internal.MethodContext
		orchName          internal.OrchestrationName
		isError           error
		loadOrchestration internal.Orchestration
	}{
		"empty name": {
			stageName: createStageName(t, ""),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
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
		"empty address": {
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress(""),
				internal.Service{},
				internal.Method{},
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
		"empty orchestration": {
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
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
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						tc.loadOrchestration.Name(): tc.loadOrchestration,
					},
				}

				createFn := Stage(stageStore, orchStore)
				err := createFn(tc.stageName, tc.methodCtx, tc.orchName)
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !errors.Is(err, tc.isError) {
					format := "Wrong error: expected %s, got %s"
					t.Fatalf(format, tc.isError, err)
				}
				if diff := cmp.Diff(0, len(stageStore.Stages)); diff != "" {
					t.Fatalf("number of stages mismatch:\n%s", diff)
				}
			},
		)
	}
}

func TestCreateStage_AlreadyExists(t *testing.T) {
	stageName := createStageName(t, "some-name")
	methodCtx := internal.NewMethodContext(
		internal.NewAddress("some-address"),
		internal.NewService("some-service"),
		internal.NewMethod("some-method"),
	)
	orchName := createOrchestrationName(t, "orchestration")

	expStage := createStage(t, "some-name", false)
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
	if err != nil {
		t.Fatalf("first create error: %s", err)
	}
	if diff := cmp.Diff(1, len(stageStore.Stages)); diff != "" {
		t.Fatalf("first create number of stages mismatch:\n%s", diff)
	}
	s, exists := stageStore.Stages[expStage.Name()]
	if !exists {
		t.Fatalf("first created stage does not exist in storage")
	}
	cmpStage(t, expStage, s, "first stage create")

	if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
		t.Fatalf("first number of orchestrations mismatch:\n%s", diff)
	}
	o, exists := orchStore.Orchs[expOrchestration.Name()]
	if !exists {
		t.Fatalf("first updated orchestration does not exist in storage")
	}
	cmpOrchestration(t, expOrchestration, o, "first update orchestration")

	err = createFn(stageName, methodCtx, orchName)
	if err == nil {
		t.Fatalf("expected create error but got none")
	}
	var alreadyExists *stageAlreadyExists
	if !errors.As(err, &alreadyExists) {
		format := "Wrong error type: expected *%s, got %s"
		t.Fatalf(format, reflect.TypeOf(alreadyExists), reflect.TypeOf(err))
	}
	if diff := cmp.Diff(stageName.Unwrap(), alreadyExists.name); diff != "" {
		t.Fatalf("name mismatch:\n%s", diff)
	}

	if diff := cmp.Diff(1, len(stageStore.Stages)); diff != "" {
		t.Fatalf("second create number of stages mismatch:\n%s", diff)
	}
	s, exists = stageStore.Stages[expStage.Name()]
	if !exists {
		t.Fatalf("second created stage does not exist in storage")
	}
	cmpStage(t, expStage, s, "second stage create")

	if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
		t.Fatalf("second number of orchestrations mismatch:\n%s", diff)
	}
	o, exists = orchStore.Orchs[expOrchestration.Name()]
	if !exists {
		t.Fatalf("second updated orchestration does not exist in storage")
	}
	cmpOrchestration(t, expOrchestration, o, "second update orchestration")
}

func createStageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	return stageName
}

func createStage(
	t *testing.T,
	stageName string,
	requiredOnly bool,
) internal.Stage {
	var (
		service internal.Service
		method  internal.Method
	)
	name := createStageName(t, stageName)
	address := internal.NewAddress("some-address")
	if !requiredOnly {
		service = internal.NewService("some-service")
		method = internal.NewMethod("some-method")
	}
	ctx := internal.NewMethodContext(address, service, method)
	return internal.NewStage(name, ctx)
}

func cmpStage(t *testing.T, x, y internal.Stage, msg string, args ...interface{}) {
	stageCmpOpts := cmp.AllowUnexported(
		internal.Stage{},
		internal.StageName{},
		internal.MethodContext{},
		internal.Address{},
		internal.Service{},
		internal.Method{},
	)
	if diff := cmp.Diff(x, y, stageCmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
