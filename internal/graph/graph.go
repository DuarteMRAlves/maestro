package graph

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal"
)

// Graph specifies the architecture of the pipeline to be
// executed.
type Graph map[internal.StageName]*StageInfo

func (g Graph) Links() []internal.Link {
	var links []internal.Link
	for _, s := range g {
		links = append(links, s.Inputs...)
	}
	return links
}

// StageInfo stores information about stages that
// is unrelated to links.
type StageInfo struct {
	Stage   internal.Stage
	Method  internal.UnaryMethod
	Inputs  []internal.Link
	Outputs []internal.Link
}

type BuildFunc func([]internal.StageName, []internal.LinkName) (Graph, error)

type incompatibleMessageDesc struct {
	A, B internal.MessageDesc
}

func (err *incompatibleMessageDesc) Error() string {
	return fmt.Sprintf("incompatible message descriptors: %s, %s", err.A, err.B)
}

type StageLoader interface {
	Load(internal.StageName) (internal.Stage, error)
}

type LinkLoader interface {
	Load(internal.LinkName) (internal.Link, error)
}

type MethodLoader interface {
	Load(internal.MethodContext) (internal.UnaryMethod, error)
}

func NewBuildFunc(
	stageLoader StageLoader, linkLoader LinkLoader, methodLoader MethodLoader,
) BuildFunc {
	return func(stages []internal.StageName, links []internal.LinkName) (Graph, error) {
		execGraph := make(Graph, len(stages))
		for _, stageName := range stages {
			stage, err := stageLoader.Load(stageName)
			if err != nil {
				return nil, err
			}
			method, err := methodLoader.Load(stage.MethodContext())
			if err != nil {
				return nil, fmt.Errorf("load method %v: %w", stage.MethodContext(), err)
			}
			execGraph[stageName] = &StageInfo{Stage: stage, Method: method}
		}
		for _, linkName := range links {
			link, err := linkLoader.Load(linkName)
			if err != nil {
				return nil, err
			}
			source, ok := execGraph[link.Source().Stage()]
			if !ok {
				err = fmt.Errorf("stage not found %s", link.Source().Stage())
				return nil, err
			}
			target, ok := execGraph[link.Target().Stage()]
			if !ok {
				err = fmt.Errorf("stage not found %s", link.Source().Stage())
				return nil, err
			}

			sourceMsg := source.Method.Output()
			if !link.Source().Field().IsEmpty() {
				sourceMsg, err = sourceMsg.GetField(link.Source().Field())
				if err != nil {
					return nil, err
				}
			}
			targetMsg := target.Method.Input()
			if !link.Target().Field().IsEmpty() {
				targetMsg, err = targetMsg.GetField(link.Target().Field())
				if err != nil {
					return nil, err
				}
			}
			if !sourceMsg.Compatible(targetMsg) {
				return nil, &incompatibleMessageDesc{A: sourceMsg, B: targetMsg}
			}
			target.Inputs = append(target.Inputs, link)
			source.Outputs = append(source.Outputs, link)
		}
		return execGraph, nil
	}
}
