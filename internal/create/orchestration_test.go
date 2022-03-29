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

func TestCreateOrchestration(t *testing.T) {
	name := "some-name"
	orchName := createOrchestrationName(t, name)
	expected := createOrchestration(t, name, nil, nil)
	storage := mock.OrchestrationStorage{Orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Orchestration(storage)

	err := createFn(orchName)
	if err != nil {
		t.Fatalf("create error: %s", err)
	}

	if diff := cmp.Diff(1, len(storage.Orchs)); diff != "" {
		t.Fatalf("number of orchestrations mismatch:\n%s", diff)
	}

	o, exists := storage.Orchs[expected.Name()]
	if !exists {
		t.Fatalf("created orchestration does not exist in storage")
	}
	cmpOrchestration(t, expected, o, "created orchestration")
}

func TestCreateOrchestration_Err(t *testing.T) {
	tests := map[string]struct {
		name    internal.OrchestrationName
		isError error
	}{
		"empty name": {
			name:    createOrchestrationName(t, ""),
			isError: emptyOrchestrationName,
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				storage := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{},
				}

				createFn := Orchestration(storage)
				err := createFn(tc.name)
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !errors.Is(err, tc.isError) {
					format := "Wrong error: expected %s, got %s"
					t.Fatalf(format, tc.isError, err)
				}
				if diff := cmp.Diff(0, len(storage.Orchs)); diff != "" {
					t.Fatalf("number of orchestrations mismatch:\n%s", diff)
				}
			},
		)
	}
}

func TestCreateOrchestration_AlreadyExists(t *testing.T) {
	name := "some-name"
	orchName := createOrchestrationName(t, name)
	expected := createOrchestration(t, name, nil, nil)
	storage := mock.OrchestrationStorage{Orchs: map[internal.OrchestrationName]internal.Orchestration{}}

	createFn := Orchestration(storage)

	err := createFn(orchName)
	if err != nil {
		t.Fatalf("first create error: %s", err)
	}
	if diff := cmp.Diff(1, len(storage.Orchs)); diff != "" {
		t.Fatalf("number of orchestrations mismatch:\n%s", diff)
	}

	o, exists := storage.Orchs[expected.Name()]
	if !exists {
		t.Fatalf("created orchestration does not exist in storage")
	}
	cmpOrchestration(t, expected, o, "first create orchestration")

	err = createFn(orchName)
	if err == nil {
		t.Fatalf("expected create error but got none")
	}
	var alreadyExists *orchestrationAlreadyExists
	if !errors.As(err, &alreadyExists) {
		format := "Wrong error type: expected *%s, got %s"
		t.Fatalf(format, reflect.TypeOf(alreadyExists), reflect.TypeOf(err))
	}
	if diff := cmp.Diff(name, alreadyExists.name); diff != "" {
		t.Fatalf("name mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(1, len(storage.Orchs)); diff != "" {
		t.Fatalf("second create number of orchestrations mismatch:\n%s", diff)
	}

	o, exists = storage.Orchs[expected.Name()]
	if !exists {
		t.Fatalf("second created orchestration does not exist in storage")
	}
	cmpOrchestration(t, expected, o, "second create orchestration")
}

func createOrchestrationName(
	t *testing.T,
	orchName string,
) internal.OrchestrationName {
	name, err := internal.NewOrchestrationName(orchName)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", orchName, err)
	}
	return name
}

func createOrchestration(
	t *testing.T,
	orchName string,
	stages, links []string,
) internal.Orchestration {
	var (
		stageNames []internal.StageName
		linkNames  []internal.LinkName
	)
	name := createOrchestrationName(t, orchName)
	for _, s := range stages {
		stageNames = append(stageNames, createStageName(t, s))
	}
	for _, l := range links {
		linkNames = append(linkNames, createLinkName(t, l))
	}
	return internal.NewOrchestration(name, stageNames, linkNames)
}

func cmpOrchestration(
	t *testing.T, x, y internal.Orchestration, msg string, args ...interface{},
) {
	orchCmpOpts := cmp.AllowUnexported(
		internal.Orchestration{},
		internal.OrchestrationName{},
		internal.StageName{},
		internal.LinkName{},
	)
	if diff := cmp.Diff(x, y, orchCmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
