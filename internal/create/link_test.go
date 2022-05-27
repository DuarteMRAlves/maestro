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
		name         internal.LinkName
		source       internal.LinkEndpoint
		target       internal.LinkEndpoint
		pipelineName internal.PipelineName
		expLink      internal.Link
		loadPipeline internal.Pipeline
		expPipeline  internal.Pipeline
		storedStages []internal.Stage
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
			pipelineName: createPipelineName(t, "pipeline"),
			expLink: createLink(
				t,
				"some-name",
				true,
			),
			loadPipeline: internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
			),
			expPipeline: internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
				internal.WithLinks(createLinkName(t, "some-name")),
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
			pipelineName: createPipelineName(t, "pipeline"),
			expLink:      createLink(t, "some-name", false),
			loadPipeline: internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
			),
			expPipeline: internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
				internal.WithLinks(createLinkName(t, "some-name")),
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

				pipelineStore := mock.PipelineStorage{
					Pipelines: map[internal.PipelineName]internal.Pipeline{
						tc.loadPipeline.Name(): tc.loadPipeline,
					},
				}

				createFn := Link(linkStore, stageStore, pipelineStore)
				err := createFn(tc.name, tc.source, tc.target, tc.pipelineName)
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

func TestCreateLink_Err(t *testing.T) {
	tests := map[string]struct {
		name         internal.LinkName
		source       internal.LinkEndpoint
		target       internal.LinkEndpoint
		pipelineName internal.PipelineName
		expErr       error
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
			pipelineName: createPipelineName(t, "pipeline"),
			expErr:       emptyLinkName},
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
			pipelineName: createPipelineName(t, "pipeline"),
			expErr:       emptySourceStage,
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
			pipelineName: createPipelineName(t, "pipeline"),
			expErr:       emptyTargetStage,
		},
		"empty pipeline": {
			name: createLinkName(t, "some-name"),
			source: internal.NewLinkEndpoint(
				createStageName(t, "source"),
				internal.NewMessageField(""),
			),
			target: internal.NewLinkEndpoint(
				createStageName(t, "target"),
				internal.NewMessageField(""),
			),
			expErr: emptyPipelineName,
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
			pipelineName: createPipelineName(t, "pipeline"),
			expErr:       equalSourceAndTarget,
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

			pipeline := internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
			)
			pipelineStore := mock.PipelineStorage{
				Pipelines: map[internal.PipelineName]internal.Pipeline{
					createPipelineName(t, "pipeline"): pipeline,
				},
			}

			createFn := Link(linkStore, stageStore, pipelineStore)
			err := createFn(tc.name, tc.source, tc.target, tc.pipelineName)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !errors.Is(err, tc.expErr) {
				t.Fatalf("Wrong error: expected %s, got %s", tc.expErr, err)
			}

			if diff := cmp.Diff(0, len(linkStore.Links)); diff != "" {
				t.Fatalf("number of links mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
				t.Fatalf("number of pipelines mismatch:\n%s", diff)
			}
			p, exists := pipelineStore.Pipelines[pipeline.Name()]
			if !exists {
				t.Fatalf("pipeline does not exist in storage")
			}
			cmpPipeline(t, pipeline, p, "pipeline is not updated")
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
	pipelineName := createPipelineName(t, "pipeline")
	expLink := createLink(t, "some-name", false)
	storedPipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(
			createStageName(t, "source"), createStageName(t, "target"),
		),
	)
	expPipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(
			createStageName(t, "source"), createStageName(t, "target"),
		),
		internal.WithLinks(linkName),
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

	pipelineStore := mock.PipelineStorage{
		Pipelines: map[internal.PipelineName]internal.Pipeline{
			storedPipeline.Name(): storedPipeline,
		},
	}

	createFn := Link(linkStore, stageStore, pipelineStore)
	err := createFn(linkName, source, target, pipelineName)
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

	if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
		t.Fatalf("first number of pipelines mismatch:\n%s", diff)
	}
	p, exists := pipelineStore.Pipelines[expPipeline.Name()]
	if !exists {
		t.Fatalf("first updated pipeline does not exist in storage")
	}
	cmpPipeline(t, expPipeline, p, "first update pipeline")

	err = createFn(linkName, source, target, pipelineName)
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

	if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
		t.Fatalf("second number of pipelines mismatch:\n%s", diff)
	}
	p, exists = pipelineStore.Pipelines[expPipeline.Name()]
	if !exists {
		t.Fatalf("second updated pipeline does not exist in storage")
	}
	cmpPipeline(t, expPipeline, p, "second update pipeline")
}

func TestStageNotInPipeline_Error(t *testing.T) {
	tests := map[string]struct {
		name         internal.LinkName
		source       internal.LinkEndpoint
		target       internal.LinkEndpoint
		pipelineName internal.PipelineName
		expStage     internal.StageName
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
			pipelineName: createPipelineName(t, "pipeline"),
			expStage:     createStageName(t, "unknown"),
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
			pipelineName: createPipelineName(t, "pipeline"),
			expStage:     createStageName(t, "unknown"),
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

			pipeline := internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
			)
			pipelineStore := mock.PipelineStorage{
				Pipelines: map[internal.PipelineName]internal.Pipeline{
					createPipelineName(t, "pipeline"): pipeline,
				},
			}

			createFn := Link(linkStore, stageStore, pipelineStore)
			err := createFn(tc.name, tc.source, tc.target, tc.pipelineName)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			var concreteErr *stageNotInPipeline
			if !errors.As(err, &concreteErr) {
				format := "Wrong error type: expected *stageNotInPipeline, got %s"
				t.Fatalf(format, reflect.TypeOf(err))
			}
			expErr := &stageNotInPipeline{Stage: tc.expStage, Pipeline: tc.pipelineName}
			cmpOpts := cmp.AllowUnexported(
				internal.StageName{}, internal.PipelineName{},
			)
			if diff := cmp.Diff(expErr, concreteErr, cmpOpts); diff != "" {
				t.Fatalf("error mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(0, len(linkStore.Links)); diff != "" {
				t.Fatalf("number of links mismatch:\n%s", diff)
			}

			if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
				t.Fatalf("number of pipelines mismatch:\n%s", diff)
			}
			p, exists := pipelineStore.Pipelines[pipeline.Name()]
			if !exists {
				t.Fatalf("pipeline does not exist in storage")
			}
			cmpPipeline(t, pipeline, p, "pipeline is not updated")
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
			initialPipeline := internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
			)
			expPipeline := internal.NewPipeline(
				createPipelineName(t, "pipeline"),
				internal.WithStages(
					createStageName(t, "source"), createStageName(t, "target"),
				),
				internal.WithLinks(createLinkName(t, "first")),
			)
			pipelineStore := mock.PipelineStorage{
				Pipelines: map[internal.PipelineName]internal.Pipeline{
					initialPipeline.Name(): initialPipeline,
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
			createFn := Link(linkStore, stageStore, pipelineStore)
			err := createFn(
				tc.first.Name(),
				tc.first.Source(),
				tc.first.Target(),
				initialPipeline.Name(),
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

			if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
				t.Fatalf("number of pipelines mismatch:\n%s", diff)
			}
			p, exists := pipelineStore.Pipelines[expPipeline.Name()]
			if !exists {
				t.Fatalf("first create updated pipeline does not exist in storage")
			}
			cmpPipeline(t, expPipeline, p, "first create updated pipeline")

			err = createFn(
				tc.second.Name(),
				tc.second.Source(),
				tc.second.Target(),
				initialPipeline.Name(),
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

			if diff := cmp.Diff(1, len(pipelineStore.Pipelines)); diff != "" {
				t.Fatalf("second number of pipelines mismatch:\n%s", diff)
			}
			p, exists = pipelineStore.Pipelines[expPipeline.Name()]
			if !exists {
				t.Fatalf("second create updated pipeline does not exist in storage")
			}
			cmpPipeline(t, expPipeline, p, "second create update pipeline")
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
