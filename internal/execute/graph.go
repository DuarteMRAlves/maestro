package execute

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal"
)

type graph map[internal.StageName]*stageInfo

// stageInfo stores information about stages that
// is unrelated to links.
type stageInfo struct {
	stage   internal.Stage
	method  internal.UnaryMethod
	inputs  []internal.Link
	outputs []internal.Link
}

type incompatibleMessageDesc struct {
	A, B internal.MessageDesc
}

func (err *incompatibleMessageDesc) Error() string {
	return fmt.Sprintf("incompatible message descriptors: %s, %s", err.A, err.B)
}

type buildGraphFunc func([]internal.StageName, []internal.LinkName) (graph, error)

func newBuildGraphFunc(
	stageLoader StageLoader, linkLoader LinkLoader, methodLoader MethodLoader,
) buildGraphFunc {
	return func(stages []internal.StageName, links []internal.LinkName) (graph, error) {
		execGraph := make(graph, len(stages))
		for _, stageName := range stages {
			stage, err := stageLoader.Load(stageName)
			if err != nil {
				return nil, err
			}
			method, err := methodLoader.Load(stage.MethodContext())
			if err != nil {
				return nil, err
			}
			execGraph[stageName] = &stageInfo{stage: stage, method: method}
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

			sourceMsg := source.method.Output()
			if !link.Source().Field().IsEmpty() {
				sourceMsg, err = sourceMsg.GetField(link.Source().Field())
				if err != nil {
					return nil, err
				}
			}
			targetMsg := target.method.Input()
			if !link.Target().Field().IsEmpty() {
				targetMsg, err = targetMsg.GetField(link.Target().Field())
				if err != nil {
					return nil, err
				}
			}
			if !sourceMsg.Compatible(targetMsg) {
				return nil, &incompatibleMessageDesc{A: sourceMsg, B: targetMsg}
			}
			target.inputs = append(target.inputs, link)
			source.outputs = append(source.outputs, link)
		}
		return execGraph, nil
	}
}
