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
		name         internal.StageName
		methodCtx    internal.MethodContext
		pipelineName internal.PipelineName
		expStage     internal.Stage
		loadPipeline internal.Pipeline
		expPipeline  internal.Pipeline
	}{
		"required fields": {
			name: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
			),
			pipelineName: createPipelineName(t, "pipeline"),
			expStage:     createStage(t, "some-name", true),
			loadPipeline: createPipeline(t, "pipeline", nil, nil),
			expPipeline: createPipeline(
				t, "pipeline", []string{"some-name"}, []string{},
			),
		},
		"all fields": {
			name: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.NewService("some-service"),
				internal.NewMethod("some-method"),
			),
			pipelineName: createPipelineName(t, "pipeline"),
			expStage:     createStage(t, "some-name", false),
			loadPipeline: createPipeline(t, "pipeline", nil, nil),
			expPipeline: createPipeline(
				t, "pipeline", []string{"some-name"}, []string{},
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

				pipelineStore := mock.PipelineStorage{
					Pipelines: map[internal.PipelineName]internal.Pipeline{
						tc.loadPipeline.Name(): tc.loadPipeline,
					},
				}

				createFn := Stage(stageStore, pipelineStore)

				err := createFn(tc.name, tc.methodCtx, tc.pipelineName)
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

				if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
					t.Fatalf("number of pipelines mismatch:\n%s", diff)
				}
				p, exists := pipelineStore.Pipelines[tc.expPipeline.Name()]
				if !exists {
					t.Fatalf("updated pipeline does not exist in storage")
				}
				cmpPipeline(t, tc.expPipeline, p, "updated pipeline")
			},
		)
	}
}

func TestCreateStage_Err(t *testing.T) {
	tests := map[string]struct {
		stageName    internal.StageName
		methodCtx    internal.MethodContext
		pipelineName internal.PipelineName
		isError      error
		loadPipeline internal.Pipeline
	}{
		"empty name": {
			stageName: createStageName(t, ""),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
			),
			pipelineName: createPipelineName(t, "pipeline"),
			isError:      emptyStageName,
			loadPipeline: createPipeline(t, "pipeline", nil, nil),
		},
		"empty address": {
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress(""),
				internal.Service{},
				internal.Method{},
			),
			pipelineName: createPipelineName(t, "pipeline"),
			isError:      emptyAddress,
			loadPipeline: createPipeline(t, "pipeline", nil, nil),
		},
		"empty pipeline": {
			stageName: createStageName(t, "some-name"),
			methodCtx: internal.NewMethodContext(
				internal.NewAddress("some-address"),
				internal.Service{},
				internal.Method{},
			),
			pipelineName: createPipelineName(t, ""),
			isError:      emptyPipelineName,
			loadPipeline: createPipeline(t, "pipeline", nil, nil),
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}

				pipelineStore := mock.PipelineStorage{
					Pipelines: map[internal.PipelineName]internal.Pipeline{
						tc.loadPipeline.Name(): tc.loadPipeline,
					},
				}

				createFn := Stage(stageStore, pipelineStore)
				err := createFn(tc.stageName, tc.methodCtx, tc.pipelineName)
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
	pipelineName := createPipelineName(t, "pipeline")

	expStage := createStage(t, "some-name", false)
	storedPipeline := createPipeline(t, "pipeline", nil, nil)
	expPipeline := createPipeline(
		t, "pipeline", []string{"some-name"}, []string{},
	)

	stageStore := mock.StageStorage{
		Stages: map[internal.StageName]internal.Stage{},
	}

	pipelineStore := mock.PipelineStorage{
		Pipelines: map[internal.PipelineName]internal.Pipeline{
			storedPipeline.Name(): storedPipeline,
		},
	}

	createFn := Stage(stageStore, pipelineStore)

	err := createFn(stageName, methodCtx, pipelineName)
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

	if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
		t.Fatalf("first number of pipelines mismatch:\n%s", diff)
	}
	p, exists := pipelineStore.Pipelines[expPipeline.Name()]
	if !exists {
		t.Fatalf("first updated pipeline does not exist in storage")
	}
	cmpPipeline(t, expPipeline, p, "first update pipeline")

	err = createFn(stageName, methodCtx, pipelineName)
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

	if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
		t.Fatalf("second number of pipelines mismatch:\n%s", diff)
	}
	p, exists = pipelineStore.Pipelines[expPipeline.Name()]
	if !exists {
		t.Fatalf("second updated pipeline does not exist in storage")
	}
	cmpPipeline(t, expPipeline, p, "second update pipeline")
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
