package execute

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type StageLoader interface {
	Load(internal.StageName) (internal.Stage, error)
}

type LinkLoader interface {
	Load(internal.LinkName) (internal.Link, error)
}

type MethodLoader interface {
	Load(internal.MethodContext) (internal.UnaryMethod, error)
}

type Builder func(orchestration internal.Orchestration) (execution, error)

func NewBuilder(
	stageLoader StageLoader,
	linkLoader LinkLoader,
	methodLoader MethodLoader,
) Builder {
	return func(orchestration internal.Orchestration) (execution, error) {
		stageNames := orchestration.Stages()
		stageCtxs := make(map[internal.StageName]stageContext, len(stageNames))
		for _, n := range stageNames {
			s, err := stageLoader.Load(n)
			if err != nil {
				return execution{}, err
			}
			m, err := methodLoader.Load(s.MethodContext())
			if err != nil {
				return execution{}, err
			}
			stageCtxs[n] = stageContext{stage: s, method: m}
		}

		linkNames := orchestration.Links()
		links := make(map[internal.LinkName]linkContext, len(linkNames))
		for _, n := range linkNames {
			link, err := linkLoader.Load(n)
			if err != nil {
				return execution{}, err
			}
			linkCtx := linkContext{link: link, ch: make(chan state)}
			links[n] = linkCtx

			source, ok := stageCtxs[link.Source().Stage()]
			if !ok {
				err = fmt.Errorf("stage not found %s", link.Source().Stage())
				return execution{}, err
			}
			target, ok := stageCtxs[link.Target().Stage()]
			if !ok {
				err = fmt.Errorf("stage not found %s", link.Source().Stage())
				return execution{}, err
			}

			sourceMsg := source.method.Output()
			if link.Source().Field().Present() {
				sourceMsg, err = sourceMsg.GetField(link.Source().Field().Unwrap())
				if err != nil {
					return execution{}, err
				}
			}
			targetMsg := target.method.Input()
			if link.Target().Field().Present() {
				targetMsg, err = targetMsg.GetField(link.Target().Field().Unwrap())
				if err != nil {
					return execution{}, err
				}
			}
			if !sourceMsg.Compatible(targetMsg) {
				err = &internal.IncompatibleMessageDesc{
					A: sourceMsg,
					B: targetMsg,
				}
				return execution{}, err
			}

			target.inputs = append(target.inputs, linkCtx)
			source.outputs = append(target.outputs, linkCtx)
		}

		stages := newStageMap()
		for name, stageCtx := range stageCtxs {
			var (
				inChan, outChan chan state
				aux             Stage
				err             error
			)
			inChan, aux, err = stageCtx.buildInputResources()
			if err != nil {
				return execution{}, err
			}
			if aux != nil {
				stages.addInputStage(name, aux)
			}
			outChan, aux = stageCtx.buildOutputResources()
			if aux != nil {
				stages.addOutputStage(name, aux)
			}
			invokeFn := stageCtx.method.InvokeFn()
			rpcStage := newUnaryStage(inChan, outChan, invokeFn)
			stages.addRpcStage(name, rpcStage)
		}

		return newExecution(stages), nil
	}
}

type stageContext struct {
	stage   internal.Stage
	method  internal.UnaryMethod
	inputs  []linkContext
	outputs []linkContext
}

type linkContext struct {
	link internal.Link
	ch   chan state
}

func (ctx stageContext) buildInputResources() (chan state, Stage, error) {
	switch len(ctx.inputs) {
	case 0:
		ch := make(chan state)
		s := newSourceStage(1, ctx.method.Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		if ctx.inputs[0].link.Target().Field().Present() {
			output, s := ctx.buildMergeStage()
			return output, s, nil
		}
		return ctx.inputs[0].ch, nil, nil
	default:
		output, s := ctx.buildMergeStage()
		return output, s, nil
	}
}

func (ctx stageContext) buildMergeStage() (chan state, *mergeStage) {
	fields := make([]internal.MessageField, 0, len(ctx.inputs))
	// channels where the stage will receive the several inputs.
	inputs := make([]<-chan state, 0, len(ctx.inputs))
	// channel where the stage will send the constructed messages.
	outputChan := make(chan state)
	for _, l := range ctx.inputs {
		fields = append(fields, l.link.Target().Field().Unwrap())
		inputs = append(inputs, l.ch)
	}
	gen := ctx.method.Input().EmptyGen()
	return outputChan, newMergeStage(fields, inputs, outputChan, gen)
}

func (ctx stageContext) buildOutputResources() (chan state, Stage) {
	switch len(ctx.outputs) {
	case 0:
		ch := make(chan state)
		s := newSinkStage(ch)
		return ch, s
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		if ctx.outputs[0].link.Source().Field().Present() {
			return ctx.buildSplitStage()
		}
		return ctx.outputs[0].ch, nil
	default:
		return ctx.buildSplitStage()
	}
}

func (ctx stageContext) buildSplitStage() (chan state, Stage) {
	fields := make([]internal.OptionalMessageField, 0, len(ctx.outputs))
	// channel where the stage will send the produced states.
	inputChan := make(chan state)
	// channels to split the received states.
	outputs := make([]chan<- state, 0, len(ctx.outputs))
	for _, l := range ctx.outputs {
		fields = append(fields, l.link.Source().Field())
		outputs = append(outputs, l.ch)
	}
	return inputChan, newSplitStage(fields, inputChan, outputs)
}
