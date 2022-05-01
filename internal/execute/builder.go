package execute

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal"
)

const defaultChanSize = 10

type StageLoader interface {
	Load(internal.StageName) (internal.Stage, error)
}

type LinkLoader interface {
	Load(internal.LinkName) (internal.Link, error)
}

type MethodLoader interface {
	Load(internal.MethodContext) (internal.UnaryMethod, error)
}

type Builder func(pipeline internal.Pipeline) (*Execution, error)

func NewBuilder(
	stageLoader StageLoader,
	linkLoader LinkLoader,
	methodLoader MethodLoader,
	logger Logger,
) Builder {
	return func(pipeline internal.Pipeline) (*Execution, error) {
		buildGraph := newBuildGraphFunc(stageLoader, linkLoader, methodLoader)

		execGraph, err := buildGraph(pipeline.Stages(), pipeline.Links())
		if err != nil {
			return nil, fmt.Errorf("build execution graph: %w", err)
		}

		// allChans stores all the channels, including the ones for aux stages.
		var allChans []chan onlineState
		// linkChans stores the channels associates with the pipeline links.
		linkChans := make(map[internal.LinkName]chan onlineState)

		for _, l := range pipeline.Links() {
			ch := make(chan onlineState, defaultChanSize)
			allChans = append(allChans, ch)
			linkChans[l] = ch
		}

		stages := newStageMap()
		for name, info := range execGraph {
			var (
				inChan, outChan chan onlineState
				aux             Stage
				err             error
			)
			fmt.Printf("On stage %s\n", name)
			inChan, aux, err = buildInputResources(info, &allChans, linkChans)
			if err != nil {
				return nil, err
			}
			if aux != nil {
				stages.addInputStage(name, aux)
			}
			outChan, aux = buildOutputResources(info, &allChans, linkChans)
			if aux != nil {
				stages.addOutputStage(name, aux)
			}
			address := info.stage.MethodContext().Address()
			clientBuilder := info.method.ClientBuilder()
			rpcStage := newOnlineUnaryStage(
				name, inChan, outChan, address, clientBuilder, logger,
			)
			stages.addRpcStage(name, rpcStage)
		}

		return newExecution(stages, allChans, logger), nil
	}
}

func buildInputResources(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage, error) {
	switch len(info.inputs) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOnlineSourceStage(1, info.method.Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		l := info.inputs[0]
		if !l.Target().Field().IsEmpty() {
			output, s := buildMergeStage(info, allChans, linkChans)
			return output, s, nil
		}
		return linkChans[l.Name()], nil, nil
	default:
		output, s := buildMergeStage(info, allChans, linkChans)
		return output, s, nil
	}
}

func buildMergeStage(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, *onlineMergeStage) {
	fields := make([]internal.MessageField, 0, len(info.inputs))
	// channels where the stage will receive the several inputs.
	inputs := make([]<-chan onlineState, 0, len(info.inputs))
	// channel where the stage will send the constructed messages.
	outputChan := make(chan onlineState, defaultChanSize)
	*allChans = append(*allChans, outputChan)
	for _, l := range info.inputs {
		fields = append(fields, l.Target().Field())
		inputs = append(inputs, linkChans[l.Name()])
	}
	gen := info.method.Input().EmptyGen()
	return outputChan, newOnlineMergeStage(fields, inputs, outputChan, gen)
}

func buildOutputResources(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	switch len(info.outputs) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOnlineSinkStage(ch)
		return ch, s
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		l := info.outputs[0]
		if !l.Source().Field().IsEmpty() {
			return buildSplitStage(info, allChans, linkChans)
		}
		return linkChans[l.Name()], nil
	default:
		return buildSplitStage(info, allChans, linkChans)
	}
}

func buildSplitStage(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(info.outputs))
	// channel where the stage will send the produced states.
	inputChan := make(chan onlineState, defaultChanSize)
	*allChans = append(*allChans, inputChan)
	// channels to split the received states.
	outputs := make([]chan<- onlineState, 0, len(info.outputs))
	for _, l := range info.outputs {
		fields = append(fields, l.Source().Field())
		outputs = append(outputs, linkChans[l.Name()])
	}
	return inputChan, newOnlineSplitStage(fields, inputChan, outputs)
}
