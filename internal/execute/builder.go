package execute

import (
	"fmt"
	"time"

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
		graphBuildFunc := newBuildGraphFunc(stageLoader, linkLoader, methodLoader)

		execGraph, err := graphBuildFunc(pipeline.Stages(), pipeline.Links())
		if err != nil {
			return nil, fmt.Errorf("build execution graph: %w", err)
		}

		switch pipeline.Mode() {
		case internal.OfflineExecution:
			return buildOfflineExecution(execGraph, logger)
		case internal.OnlineExecution:
			return buildOnlineExecution(execGraph, logger)
		default:
			return nil, fmt.Errorf("unknown execution format: %v", pipeline.Mode())
		}
	}
}

func buildOnlineExecution(execGraph graph, logger Logger) (*Execution, error) {
	// allChans stores all the channels, including the ones for aux stages.
	// linkChans stores the channels associates with the pipeline links.
	var allChans []chan onlineState

	linkChans := make(map[internal.LinkName]chan onlineState)

	for _, l := range execGraph.links() {
		ch := make(chan onlineState, defaultChanSize)
		allChans = append(allChans, ch)
		linkChans[l.Name()] = ch
	}

	stages := newStageMap()
	for name, info := range execGraph {
		var (
			inChan, outChan chan onlineState
			aux             Stage
			err             error
		)
		inChan, aux, err = buildOnlineInputResources(info, &allChans, linkChans)
		if err != nil {
			return nil, err
		}
		if aux != nil {
			stages.addInputStage(name, aux)
		}
		outChan, aux = buildOnlineOutputResources(info, &allChans, linkChans)
		if aux != nil {
			stages.addOutputStage(name, aux)
		}
		address := info.stage.MethodContext().Address()
		clientBuilder := info.method.ClientBuilder()
		rpcStage := newOnlineUnary(
			name, inChan, outChan, address, clientBuilder, logger,
		)
		stages.addRpcStage(name, rpcStage)
	}
	drainFunc := newChanDrainer(5*time.Millisecond, allChans...)
	return newExecution(stages, drainFunc, logger), nil
}

func buildOnlineInputResources(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage, error) {
	switch len(info.inputs) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOnlineSource(1, info.method.Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		l := info.inputs[0]
		if !l.Target().Field().IsEmpty() {
			output, s := buildOnlineMergeStage(info, allChans, linkChans)
			return output, s, nil
		}
		return linkChans[l.Name()], nil, nil
	default:
		output, s := buildOnlineMergeStage(info, allChans, linkChans)
		return output, s, nil
	}
}

func buildOnlineMergeStage(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
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
	return outputChan, newOnlineMerge(fields, inputs, outputChan, gen)
}

func buildOnlineOutputResources(
	info *stageInfo, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	switch len(info.outputs) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		return ch, newOnlineSink(ch)
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		l := info.outputs[0]
		if !l.Source().Field().IsEmpty() {
			return buildOnlineSplitStage(info, allChans, linkChans)
		}
		return linkChans[l.Name()], nil
	default:
		return buildOnlineSplitStage(info, allChans, linkChans)
	}
}

func buildOnlineSplitStage(
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
	return inputChan, newOnlineSplit(fields, inputChan, outputs)
}

func buildOfflineExecution(execGraph graph, logger Logger) (*Execution, error) {
	// allChans stores all the channels, including the ones for aux stages.
	// linkChans stores the channels associates with the pipeline links.
	var allChans []chan offlineState

	linkChans := make(map[internal.LinkName]chan offlineState)

	for _, l := range execGraph.links() {
		ch := make(chan offlineState, defaultChanSize)
		allChans = append(allChans, ch)
		linkChans[l.Name()] = ch
	}

	stages := newStageMap()
	for name, info := range execGraph {
		var (
			inChan, outChan chan offlineState
			aux             Stage
			err             error
		)
		inChan, aux, err = buildInputResources(info, &allChans, linkChans)
		if err != nil {
			return nil, err
		}
		if aux != nil {
			stages.addInputStage(name, aux)
		}
		outChan, aux = buildOfflineOutputResources(info, &allChans, linkChans)
		if aux != nil {
			stages.addOutputStage(name, aux)
		}
		address := info.stage.MethodContext().Address()
		clientBuilder := info.method.ClientBuilder()
		rpcStage := newOfflineUnary(
			name, inChan, outChan, address, clientBuilder, logger,
		)
		stages.addRpcStage(name, rpcStage)
	}
	drainFunc := newChanDrainer(5*time.Millisecond, allChans...)
	return newExecution(stages, drainFunc, logger), nil
}

func buildInputResources(
	info *stageInfo, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage, error) {
	switch len(info.inputs) {
	case 0:
		ch := make(chan offlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOfflineSource(info.method.Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		l := info.inputs[0]
		if !l.Target().Field().IsEmpty() {
			output, s := buildOfflineMergeStage(info, allChans, linkChans)
			return output, s, nil
		}
		return linkChans[l.Name()], nil, nil
	default:
		output, s := buildOfflineMergeStage(info, allChans, linkChans)
		return output, s, nil
	}
}

func buildOfflineMergeStage(
	info *stageInfo, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(info.inputs))
	// channels where the stage will receive the several inputs.
	inputs := make([]<-chan offlineState, 0, len(info.inputs))
	// channel where the stage will send the constructed messages.
	outputChan := make(chan offlineState, defaultChanSize)
	*allChans = append(*allChans, outputChan)
	for _, l := range info.inputs {
		fields = append(fields, l.Target().Field())
		inputs = append(inputs, linkChans[l.Name()])
	}
	gen := info.method.Input().EmptyGen()
	return outputChan, newOfflineMerge(fields, inputs, outputChan, gen)
}

func buildOfflineOutputResources(
	info *stageInfo, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	switch len(info.outputs) {
	case 0:
		ch := make(chan offlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		return ch, newOfflineSink(ch)
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		l := info.outputs[0]
		if !l.Source().Field().IsEmpty() {
			return buildOfflineSplitStage(info, allChans, linkChans)
		}
		return linkChans[l.Name()], nil
	default:
		return buildOfflineSplitStage(info, allChans, linkChans)
	}
}

func buildOfflineSplitStage(
	info *stageInfo, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(info.outputs))
	// channel where the stage will send the produced states.
	inputChan := make(chan offlineState, defaultChanSize)
	*allChans = append(*allChans, inputChan)
	// channels to split the received states.
	outputs := make([]chan<- offlineState, 0, len(info.outputs))
	for _, l := range info.outputs {
		fields = append(fields, l.Source().Field())
		outputs = append(outputs, linkChans[l.Name()])
	}
	return inputChan, newOfflineSplit(fields, inputChan, outputs)
}
