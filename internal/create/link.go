package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type LinkSaver interface {
	Save(internal.Link) error
}

type LinkLoader interface {
	Load(internal.LinkName) (internal.Link, error)
}

type LinkStorage interface {
	LinkSaver
	LinkLoader
}

var (
	emptyLinkName        = errors.New("empty link name")
	emptySourceStage     = errors.New("empty source stage")
	emptyTargetStage     = errors.New("empty target stage")
	equalSourceAndTarget = errors.New("equal source and target stages")
)

type linkAlreadyExists struct{ name string }

func (err *linkAlreadyExists) Error() string {
	return fmt.Sprintf("link '%s' already exists", err.name)
}

type stageNotInPipeline struct {
	Pipeline internal.PipelineName
	Stage    internal.StageName
}

func (err *stageNotInPipeline) Error() string {
	format := "stage '%s' not found in pipeline '%s'."
	return fmt.Sprintf(format, err.Stage, err.Pipeline)
}

type incompatibleLinks struct {
	A, B string
}

func (err *incompatibleLinks) Error() string {
	return fmt.Sprintf("incompatible links: %s, %s", err.A, err.B)
}

func Link(
	linkStorage LinkStorage,
	stageLoader StageLoader,
	pipelineStorage PipelineStorage,
) func(
	internal.LinkName,
	internal.LinkEndpoint,
	internal.LinkEndpoint,
	internal.PipelineName,
) error {
	return func(
		name internal.LinkName,
		source, target internal.LinkEndpoint,
		pipelineName internal.PipelineName,
	) error {
		if name.IsEmpty() {
			return emptyLinkName
		}

		sourceStage := source.Stage()
		if sourceStage.IsEmpty() {
			return emptySourceStage
		}

		targetStage := target.Stage()
		if targetStage.IsEmpty() {
			return emptyTargetStage
		}

		if pipelineName.IsEmpty() {
			return emptyPipelineName
		}

		_, err := linkStorage.Load(name)
		if err == nil {
			return &linkAlreadyExists{name: name.Unwrap()}
		}
		var nf interface{ NotFound() }
		if !errors.As(err, &nf) {
			return err
		}

		pipeline, err := pipelineStorage.Load(pipelineName)
		if err != nil {
			return fmt.Errorf("add link %s: %w", name, err)
		}

		_, err = stageLoader.Load(sourceStage)
		if err != nil {
			return err
		}

		_, err = stageLoader.Load(targetStage)
		if err != nil {
			return err
		}

		if !existsStage(sourceStage, pipeline.Stages()...) {
			return &stageNotInPipeline{Pipeline: pipelineName, Stage: sourceStage}
		}
		if !existsStage(targetStage, pipeline.Stages()...) {
			return &stageNotInPipeline{Pipeline: pipelineName, Stage: targetStage}
		}

		if sourceStage == targetStage {
			return equalSourceAndTarget
		}

		targetLinks, err := linksForTarget(linkStorage, targetStage, pipeline.Links()...)
		if err != nil {
			return fmt.Errorf("create link %s: %w", name, err)
		}

		for _, l := range targetLinks {
			// 1. Target receives entire message from this link but another exists.
			// 2. Target already receives entire message from existing link.
			// 3. Target receives same field from both links.
			if target.Field().IsEmpty() ||
				l.Target().Field().IsEmpty() ||
				target.Field().Unwrap() == l.Target().Field().Unwrap() {
				return &incompatibleLinks{A: name.Unwrap(), B: l.Name().Unwrap()}
			}
		}

		links := pipeline.Links()
		links = append(links, name)
		pipeline = internal.NewPipeline(pipeline.Name(), pipeline.Stages(), links)

		err = pipelineStorage.Save(pipeline)
		if err != nil {
			return fmt.Errorf("add link %s: %w", name, err)
		}

		l := internal.NewLink(name, source, target)

		return linkStorage.Save(l)
	}
}

func existsStage(name internal.StageName, ss ...internal.StageName) bool {
	for _, s := range ss {
		if s == name {
			return true
		}
	}
	return false
}

// Retrieves the links that have a specific target stage.
func linksForTarget(
	linkLoader LinkLoader, target internal.StageName, links ...internal.LinkName,
) ([]internal.Link, error) {
	var ret []internal.Link
	for _, n := range links {
		l, err := linkLoader.Load(n)
		if err != nil {
			return nil, err
		}
		if l.Target().Stage() == target {
			ret = append(ret, l)
		}
	}
	return ret, nil
}
