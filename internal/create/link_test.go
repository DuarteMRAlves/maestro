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

func TestCreateLink(t *testing.T) {
	tests := map[string]struct {
		name              internal.LinkName
		source            internal.LinkEndpoint
		target            internal.LinkEndpoint
		orchName          internal.OrchestrationName
		expLink           internal.Link
		loadOrchestration internal.Orchestration
		expOrch           internal.Orchestration
		storedStages      []internal.Stage
	}{
		"required fields": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.MessageField{},
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.MessageField{},
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expLink: createLink(
				t,
				"some-name",
				true,
			),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{},
			),
			expOrch: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{"some-name"},
			),
			storedStages: []internal.Stage{
				createStage(t, "source", true),
				createStage(t, "target", true),
			},
		},
		"all fields": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField("source-field"),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField("target-field"),
			),
			orchName: createOrchestrationName(t, "orchestration"),
			expLink:  createLink(t, "some-name", false),
			loadOrchestration: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{},
			),
			expOrch: createOrchestration(
				t,
				"orchestration",
				[]string{"source", "target"},
				[]string{"some-name"},
			),
			storedStages: []internal.Stage{
				createStage(t, "source", false),
				createStage(t, "target", false),
			},
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

				stageStore := mock.StageStorage{
					Stages: map[internal.StageName]internal.Stage{},
				}
				for _, s := range tc.storedStages {
					stageStore.Stages[s.Name()] = s
				}

				orchStore := mock.OrchestrationStorage{
					Orchs: map[internal.OrchestrationName]internal.Orchestration{
						tc.loadOrchestration.Name(): tc.loadOrchestration,
					},
				}

				createFn := Link(linkStore, stageStore, orchStore)
				err := createFn(tc.name, tc.source, tc.target, tc.orchName)
				if err != nil {
					t.Fatalf("create error: %s", err)
				}

				if diff := cmp.Diff(1, len(linkStore.Links)); diff != "" {
					t.Fatalf("number of links mismatch:\n%s", diff)
				}
				l, exists := linkStore.Links[tc.expLink.Name()]
				if !exists {
					t.Fatalf("created stage does not exist in storage")
				}
				cmpLink(t, tc.expLink, l, "created link")

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

func TestCreateLink_Err(t *testing.T) {
	tests := map[string]struct {
		name     internal.LinkName
		source   internal.LinkEndpoint
		target   internal.LinkEndpoint
		orchName internal.OrchestrationName
		expErr   error
	}{
		"empty name": {
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expErr:   emptyLinkName},
		"empty source stage": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, ""),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expErr:   emptySourceStage,
		},
		"empty target stage": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, ""),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expErr:   emptyTargetStage,
		},
		"empty orchestration": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField(""),
			),
			expErr: emptyOrchestrationName,
		},
		"equal source and target": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expErr:   equalSourceAndTarget,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

			stageStore := mock.StageStorage{
				Stages: map[internal.StageName]internal.Stage{
					createStageName(t, "source"): createStage(t, "source", true),
					createStageName(t, "target"): createStage(t, "target", true),
				},
			}

			orch := createOrchestration(t, "orch", []string{"source", "target"}, nil)
			orchStore := mock.OrchestrationStorage{
				Orchs: map[internal.OrchestrationName]internal.Orchestration{
					createOrchestrationName(t, "orch"): orch,
				},
			}

			createFn := Link(linkStore, stageStore, orchStore)
			err := createFn(tc.name, tc.source, tc.target, tc.orchName)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !errors.Is(err, tc.expErr) {
				t.Fatalf("Wrong error: expected %s, got %s", tc.expErr, err)
			}

			if diff := cmp.Diff(0, len(linkStore.Links)); diff != "" {
				t.Fatalf("number of links mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
				t.Fatalf("number of orchestrations mismatch:\n%s", diff)
			}
			o, exists := orchStore.Orchs[orch.Name()]
			if !exists {
				t.Fatalf("orchestration does not exist in storage")
			}
			cmpOrchestration(t, orch, o, "orchestration is not updated")
		})
	}
}

func TestCreateLink_AlreadyExists(t *testing.T) {
	linkName := createLinkName(t, "some-name")
	source := internal.NewLinkEndpoint(
		createStageName(t, "source"),
		internal.NewMessageField("source-field"),
	)
	target := internal.NewLinkEndpoint(
		createStageName(t, "target"),
		internal.NewMessageField("target-field"),
	)
	orchName := createOrchestrationName(t, "orchestration")
	expLink := createLink(t, "some-name", false)
	storedOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"source", "target"},
		[]string{},
	)
	expOrchestration := createOrchestration(
		t,
		"orchestration",
		[]string{"source", "target"},
		[]string{"some-name"},
	)
	storedStages := []internal.Stage{
		createStage(t, "source", false),
		createStage(t, "target", false),
	}

	linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

	stageStore := mock.StageStorage{
		Stages: map[internal.StageName]internal.Stage{},
	}
	for _, s := range storedStages {
		stageStore.Stages[s.Name()] = s
	}

	orchStore := mock.OrchestrationStorage{
		Orchs: map[internal.OrchestrationName]internal.Orchestration{
			storedOrchestration.Name(): storedOrchestration,
		},
	}

	createFn := Link(linkStore, stageStore, orchStore)
	err := createFn(linkName, source, target, orchName)
	if err != nil {
		t.Fatalf("first create error: %s", err)
	}
	if diff := cmp.Diff(1, len(linkStore.Links)); diff != "" {
		t.Fatalf("first create number of links mismatch:\n%s", diff)
	}
	l, exists := linkStore.Links[expLink.Name()]
	if !exists {
		t.Fatalf("first created link does not exist in storage")
	}
	cmpLink(t, expLink, l, "first create link")

	if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
		t.Fatalf("first number of orchestrations mismatch:\n%s", diff)
	}
	o, exists := orchStore.Orchs[expOrchestration.Name()]
	if !exists {
		t.Fatalf("first updated orchestration does not exist in storage")
	}
	cmpOrchestration(t, expOrchestration, o, "first update orchestration")

	err = createFn(linkName, source, target, orchName)
	if err == nil {
		t.Fatalf("expected create error but got none")
	}
	var alreadyExists *linkAlreadyExists
	if !errors.As(err, &alreadyExists) {
		format := "Wrong error type: expected *%s, got %s"
		t.Fatalf(format, reflect.TypeOf(alreadyExists), reflect.TypeOf(err))
	}
	if diff := cmp.Diff(linkName.Unwrap(), alreadyExists.name); diff != "" {
		t.Fatalf("name mismatch:\n%s", diff)
	}

	if diff := cmp.Diff(1, len(linkStore.Links)); diff != "" {
		t.Fatalf("second create number of links mismatch:\n%s", diff)
	}
	l, exists = linkStore.Links[expLink.Name()]
	if !exists {
		t.Fatalf("second created link does not exist in storage")
	}
	cmpLink(t, expLink, l, "second create link")

	if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
		t.Fatalf("second number of orchestrations mismatch:\n%s", diff)
	}
	o, exists = orchStore.Orchs[expOrchestration.Name()]
	if !exists {
		t.Fatalf("second updated orchestration does not exist in storage")
	}
	cmpOrchestration(t, expOrchestration, o, "second update orchestration")
}

func TestStageNotInOrchestration_Error(t *testing.T) {
	tests := map[string]struct {
		name     internal.LinkName
		source   internal.LinkEndpoint
		target   internal.LinkEndpoint
		orchName internal.OrchestrationName
		expStage internal.StageName
	}{
		"source": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "unknown"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expStage: createStageName(t, "unknown"),
		},
		"target": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "unknown"),
				internal.NewMessageField(""),
			),
			orchName: createOrchestrationName(t, "orch"),
			expStage: createStageName(t, "unknown"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}

			stageStore := mock.StageStorage{
				Stages: map[internal.StageName]internal.Stage{
					createStageName(t, "source"):  createStage(t, "source", true),
					createStageName(t, "target"):  createStage(t, "target", true),
					createStageName(t, "unknown"): createStage(t, "unknown", true),
				},
			}

			orch := createOrchestration(t, "orch", []string{"source", "target"}, nil)
			orchStore := mock.OrchestrationStorage{
				Orchs: map[internal.OrchestrationName]internal.Orchestration{
					createOrchestrationName(t, "orch"): orch,
				},
			}

			createFn := Link(linkStore, stageStore, orchStore)
			err := createFn(tc.name, tc.source, tc.target, tc.orchName)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			var concreteErr *stageNotInOrchestration
			if !errors.As(err, &concreteErr) {
				format := "Wrong error type: expected *stageNotInOrchestration, got %s"
				t.Fatalf(format, reflect.TypeOf(err))
			}
			expErr := &stageNotInOrchestration{Stage: tc.expStage, Orch: tc.orchName}
			cmpOpts := cmp.AllowUnexported(
				internal.StageName{}, internal.OrchestrationName{},
			)
			if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
				t.Fatalf("error mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(0, len(linkStore.Links)); diff != "" {
				t.Fatalf("number of links mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
				t.Fatalf("number of orchestrations mismatch:\n%s", diff)
			}
			o, exists := orchStore.Orchs[orch.Name()]
			if !exists {
				t.Fatalf("orchestration does not exist in storage")
			}
			cmpOrchestration(t, orch, o, "orchestration is not updated")
		})
	}
}

func TestCreateLink_IncompatibleLinks(t *testing.T) {
	tests := map[string]struct {
		first, second internal.Link
	}{
		"first entire, second entire": {
			first:  createLink(t, "first", true),
			second: createLink(t, "second", true),
		},
		"first entire, second field": {
			first:  createLink(t, "first", true),
			second: createLink(t, "second", false),
		},
		"first field, second entire": {
			first:  createLink(t, "first", false),
			second: createLink(t, "second", true),
		},
		"first field, second field": {
			first:  createLink(t, "first", false),
			second: createLink(t, "second", false),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initialOrch := createOrchestration(
				t, "orchestration", []string{"source", "target"}, nil,
			)
			expOrch := createOrchestration(
				t, "orchestration", []string{"source", "target"}, []string{"first"},
			)
			orchStore := mock.OrchestrationStorage{
				Orchs: map[internal.OrchestrationName]internal.Orchestration{
					initialOrch.Name(): initialOrch,
				},
			}

			sourceStage := createStage(t, "source", true)
			targetStage := createStage(t, "target", true)
			stageStore := mock.StageStorage{
				Stages: map[internal.StageName]internal.Stage{
					sourceStage.Name(): sourceStage,
					targetStage.Name(): targetStage,
				},
			}

			linkStore := mock.LinkStorage{Links: map[internal.LinkName]internal.Link{}}
			createFn := Link(linkStore, stageStore, orchStore)
			err := createFn(
				tc.first.Name(),
				tc.first.Source(),
				tc.first.Target(),
				initialOrch.Name(),
			)
			if err != nil {
				t.Fatalf("first create error: %s", err)
			}
			if diff := cmp.Diff(1, len(linkStore.Links)); diff != "" {
				t.Fatalf("first create number of links mismatch:\n%s", diff)
			}
			l, exists := linkStore.Links[tc.first.Name()]
			if !exists {
				t.Fatalf("first create stage does not exist in storage")
			}
			cmpLink(t, tc.first, l, "first create link")

			if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
				t.Fatalf("number of orchestrations mismatch:\n%s", diff)
			}
			o, exists := orchStore.Orchs[expOrch.Name()]
			if !exists {
				t.Fatalf("first create updated orchestration does not exist in storage")
			}
			cmpOrchestration(t, expOrch, o, "first create updated orchestration")

			err = createFn(
				tc.second.Name(),
				tc.second.Source(),
				tc.second.Target(),
				initialOrch.Name(),
			)
			if err == nil {
				t.Fatalf("second create expected create error but got nil")
			}
			var incompatibleErr *incompatibleLinks
			if !errors.As(err, &incompatibleErr) {
				format := "second create wrong error type: expected %s, got %s"
				t.Fatalf(format, reflect.TypeOf(incompatibleErr), reflect.TypeOf(err))
			}
			expError := &incompatibleLinks{
				A: tc.second.Name().Unwrap(), B: tc.first.Name().Unwrap(),
			}
			if diff := cmp.Diff(expError, incompatibleErr); diff != "" {
				t.Fatalf("second create error mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(1, len(linkStore.Links)); diff != "" {
				t.Fatalf("second create number of links mismatch:\n%s", diff)
			}
			l, exists = linkStore.Links[tc.first.Name()]
			if !exists {
				t.Fatalf("second created link does not exist in storage")
			}
			cmpLink(t, tc.first, l, "second create link")

			if diff := cmp.Diff(1, len(orchStore.Orchs)); diff != "" {
				t.Fatalf("second number of orchestrations mismatch:\n%s", diff)
			}
			o, exists = orchStore.Orchs[expOrch.Name()]
			if !exists {
				t.Fatalf("second create updated orchestration does not exist in storage")
			}
			cmpOrchestration(t, expOrch, o, "second create update orchestration")
		})
	}
}

func createLinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("create link name %s: %s", name, err)
	}
	return linkName
}

func createLink(
	t *testing.T,
	linkName string,
	requiredOnly bool,
) internal.Link {
	var (
		sourceField internal.MessageField
		targetField internal.MessageField
	)
	name := createLinkName(t, linkName)
	sourceStage := createStageName(t, "source")
	if !requiredOnly {
		sourceField = internal.NewMessageField("source-field")
	}
	sourceEndpoint := internal.NewLinkEndpoint(sourceStage, sourceField)

	targetStage := createStageName(t, "target")
	if !requiredOnly {
		targetField = internal.NewMessageField("target-field")
	}
	targetEndpoint := internal.NewLinkEndpoint(targetStage, targetField)

	return internal.NewLink(name, sourceEndpoint, targetEndpoint)
}

func cmpLink(t *testing.T, x, y internal.Link, msg string, args ...interface{}) {
	cmpOpts := cmp.AllowUnexported(
		internal.Link{},
		internal.LinkName{},
		internal.LinkEndpoint{},
		internal.StageName{},
		internal.MessageField{},
	)
	if diff := cmp.Diff(x, y, cmpOpts); diff != "" {
		prepend := fmt.Sprintf(msg, args...)
		t.Fatalf("%s:\n%s", prepend, diff)
	}
}
