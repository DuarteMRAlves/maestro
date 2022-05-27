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

func TestCreatePipeline(t *testing.T) {
	tests := map[string]struct {
		name     internal.PipelineName
		mode     internal.ExecutionMode
		expected internal.Pipeline
	}{
		"required fields": {
			name: createPipelineName(t, "some-name"),
			expected: internal.NewPipeline(
				createPipelineName(t, "some-name"),
				internal.WithOfflineExec(),
			),
		},
		"all fields": {
			name: createPipelineName(t, "some-name"),
			mode: internal.OnlineExecution,
			expected: internal.NewPipeline(
				createPipelineName(t, "some-name"),
				internal.WithOnlineExec(),
			),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := mock.PipelineStorage{Pipelines: map[internal.PipelineName]internal.Pipeline{}}

			createFn := Pipeline(storage)

			err := createFn(tc.name, tc.mode)
			if err != nil {
				t.Fatalf("create error: %s", err)
			}

			if diff := cmp.Diff(1, len(storage.Pipelines)); diff != "" {
				t.Fatalf("number of pipelines mismatch:\n%s", diff)
			}

			p, exists := storage.Pipelines[tc.expected.Name()]
			if !exists {
				t.Fatalf("created pipeline does not exist in storage")
			}
			cmpPipeline(t, tc.expected, p, "created pipeline")
		})
	}
}

func TestCreatePipeline_Err(t *testing.T) {
	tests := map[string]struct {
		name     internal.PipelineName
		mode     internal.ExecutionMode
		valError func(*testing.T, error)
	}{
		"empty name": {
			name: createPipelineName(t, ""),
			mode: internal.OnlineExecution,
			valError: func(t *testing.T, err error) {
				if !errors.Is(err, emptyPipelineName) {
					format := "Wrong error: expected %s, got %s"
					t.Fatalf(format, emptyPipelineName, err)
				}
			},
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				storage := mock.PipelineStorage{
					Pipelines: map[internal.PipelineName]internal.Pipeline{},
				}

				createFn := Pipeline(storage)
				err := createFn(tc.name, tc.mode)
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				tc.valError(t, err)
				if diff := cmp.Diff(0, len(storage.Pipelines)); diff != "" {
					t.Fatalf("number of pipelines mismatch:\n%s", diff)
				}
			},
		)
	}
}

func TestCreatePipeline_AlreadyExists(t *testing.T) {
	name := "some-name"
	pipelineName := createPipelineName(t, name)
	mode := internal.OnlineExecution
	expected := internal.NewPipeline(pipelineName, internal.WithOnlineExec())
	storage := mock.PipelineStorage{Pipelines: map[internal.PipelineName]internal.Pipeline{}}

	createFn := Pipeline(storage)

	err := createFn(pipelineName, mode)
	if err != nil {
		t.Fatalf("first create error: %s", err)
	}
	if diff := cmp.Diff(1, len(storage.Pipelines)); diff != "" {
		t.Fatalf("number of pipelines mismatch:\n%s", diff)
	}

	p, exists := storage.Pipelines[expected.Name()]
	if !exists {
		t.Fatalf("created pipeline does not exist in storage")
	}
	cmpPipeline(t, expected, p, "first create pipeline")

	err = createFn(pipelineName, mode)
	if err == nil {
		t.Fatalf("expected create error but got none")
	}
	var alreadyExists *pipelineAlreadyExists
	if !errors.As(err, &alreadyExists) {
		format := "Wrong error type: expected *%s, got %s"
		t.Fatalf(format, reflect.TypeOf(alreadyExists), reflect.TypeOf(err))
	}
	if diff := cmp.Diff(name, alreadyExists.name); diff != "" {
		t.Fatalf("name mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(1, len(storage.Pipelines)); diff != "" {
		t.Fatalf("second create number of pipelines mismatch:\n%s", diff)
	}

	p, exists = storage.Pipelines[expected.Name()]
	if !exists {
		t.Fatalf("second created pipeline does not exist in storage")
	}
	cmpPipeline(t, expected, p, "second create pipeline")
}

func createPipelineName(t *testing.T, name string) internal.PipelineName {
	pipelineName, err := internal.NewPipelineName(name)
	if err != nil {
		t.Fatalf("create pipeline name %s: %s", name, err)
	}
	return pipelineName
}

func cmpPipeline(
	t *testing.T, x, y internal.Pipeline, msg string, args ...interface{},
) {
	cmpOpts := cmp.AllowUnexported(
		internal.Pipeline{},
		internal.PipelineName{},
		internal.ExecutionMode{},
		internal.StageName{},
		internal.LinkName{},
	)
	if diff := cmp.Diff(x, y, cmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
