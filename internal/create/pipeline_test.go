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
	name := "some-name"
	pipelineName := createPipelineName(t, name)
	expected := createPipeline(t, name, nil, nil)
	storage := mock.PipelineStorage{Pipelines: map[internal.PipelineName]internal.Pipeline{}}

	createFn := Pipeline(storage)

	err := createFn(pipelineName)
	if err != nil {
		t.Fatalf("create error: %s", err)
	}

	if diff := cmp.Diff(1, len(storage.Pipelines)); diff != "" {
		t.Fatalf("number of pipelines mismatch:\n%s", diff)
	}

	p, exists := storage.Pipelines[expected.Name()]
	if !exists {
		t.Fatalf("created pipeline does not exist in storage")
	}
	cmpPipeline(t, expected, p, "created pipeline")
}

func TestCreatePipeline_Err(t *testing.T) {
	tests := map[string]struct {
		name    internal.PipelineName
		isError error
	}{
		"empty name": {
			name:    createPipelineName(t, ""),
			isError: emptyPipelineName,
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
				err := createFn(tc.name)
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !errors.Is(err, tc.isError) {
					format := "Wrong error: expected %s, got %s"
					t.Fatalf(format, tc.isError, err)
				}
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
	expected := createPipeline(t, name, nil, nil)
	storage := mock.PipelineStorage{Pipelines: map[internal.PipelineName]internal.Pipeline{}}

	createFn := Pipeline(storage)

	err := createFn(pipelineName)
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

	err = createFn(pipelineName)
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

func createPipeline(
	t *testing.T, name string, stages, links []string,
) internal.Pipeline {
	var (
		stageNames []internal.StageName
		linkNames  []internal.LinkName
	)
	pipelineName := createPipelineName(t, name)
	for _, s := range stages {
		stageNames = append(stageNames, createStageName(t, s))
	}
	for _, l := range links {
		linkNames = append(linkNames, createLinkName(t, l))
	}
	return internal.NewPipeline(pipelineName, stageNames, linkNames)
}

func cmpPipeline(
	t *testing.T, x, y internal.Pipeline, msg string, args ...interface{},
) {
	cmpOpts := cmp.AllowUnexported(
		internal.Pipeline{},
		internal.PipelineName{},
		internal.StageName{},
		internal.LinkName{},
	)
	if diff := cmp.Diff(x, y, cmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
