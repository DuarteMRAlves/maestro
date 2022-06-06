package execute

import (
	"fmt"
	"time"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/compiled"
)

const defaultChanSize = 10

type Builder func(pipeline *compiled.Pipeline) (Execution, error)

func NewBuilder(logger Logger) Builder {
	return func(pipeline *compiled.Pipeline) (Execution, error) {
		switch pipeline.Mode() {
		case internal.OfflineExecution:
			return buildOfflineExecution(pipeline, logger)
		case internal.OnlineExecution:
			return buildOnlineExecution(pipeline, logger)
		default:
			return nil, fmt.Errorf("unknown execution format: %v", pipeline.Mode())
		}
	}
}

func buildOfflineExecution(pipeline *compiled.Pipeline, logger Logger) (*offlineExecution, error) {
	// allChans stores all the channels, including the ones for aux stages.
	// linkChans stores the channels associates with the pipeline links.
	var allChans []chan offlineState

	linkChans := make(map[internal.LinkName]chan offlineState)

	pipeline.VisitLinks(func(l *internal.Link) error {
		ch := make(chan offlineState, defaultChanSize)
		allChans = append(allChans, ch)
		linkChans[l.Name()] = ch
		return nil
	})

	stages := newStageMap()
	pipeline.VisitStages(func(s *compiled.Stage) error {
		name := s.Name()
		inChan, aux, err := buildInputResources(s, &allChans, linkChans)
		if err != nil {
			return err
		}
		if aux != nil {
			stages.addInputStage(name, aux)
		}
		outChan, aux := buildOfflineOutputResources(s, &allChans, linkChans)
		if aux != nil {
			stages.addOutputStage(name, aux)
		}
		address := s.Address()
		clientBuilder := s.Method().ClientBuilder()
		rpcStage := newOfflineUnary(
			name, inChan, outChan, address, clientBuilder, logger,
		)
		stages.addRpcStage(name, rpcStage)
		return nil
	})
	return newOfflineExecution(stages, logger), nil
}

func buildInputResources(
	s *compiled.Stage, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage, error) {
	switch len(s.Inputs()) {
	case 0:
		ch := make(chan offlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOfflineSource(s.Method().Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		l := s.Inputs()[0]
		if !l.Target().Field().IsEmpty() {
			output, s := buildOfflineMergeStage(s, allChans, linkChans)
			return output, s, nil
		}
		return linkChans[l.Name()], nil, nil
	default:
		output, s := buildOfflineMergeStage(s, allChans, linkChans)
		return output, s, nil
	}
}

func buildOfflineMergeStage(
	s *compiled.Stage, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(s.Inputs()))
	// channels where the stage will receive the several inputs.
	inputs := make([]<-chan offlineState, 0, len(s.Inputs()))
	// channel where the stage will send the constructed messages.
	outputChan := make(chan offlineState, defaultChanSize)
	*allChans = append(*allChans, outputChan)
	for _, l := range s.Inputs() {
		fields = append(fields, l.Target().Field())
		inputs = append(inputs, linkChans[l.Name()])
	}
	gen := s.Method().Input().EmptyGen()
	return outputChan, newOfflineMerge(fields, inputs, outputChan, gen)
}

func buildOfflineOutputResources(
	s *compiled.Stage, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	switch len(s.Outputs()) {
	case 0:
		ch := make(chan offlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		return ch, newOfflineSink(ch)
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		l := s.Outputs()[0]
		if !l.Source().Field().IsEmpty() {
			return buildOfflineSplitStage(s, allChans, linkChans)
		}
		return linkChans[l.Name()], nil
	default:
		return buildOfflineSplitStage(s, allChans, linkChans)
	}
}

func buildOfflineSplitStage(
	s *compiled.Stage, allChans *[]chan offlineState, linkChans map[internal.LinkName]chan offlineState,
) (chan offlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(s.Outputs()))
	// channel where the stage will send the produced states.
	inputChan := make(chan offlineState, defaultChanSize)
	*allChans = append(*allChans, inputChan)
	// channels to split the received states.
	outputs := make([]chan<- offlineState, 0, len(s.Outputs()))
	for _, l := range s.Outputs() {
		fields = append(fields, l.Source().Field())
		outputs = append(outputs, linkChans[l.Name()])
	}
	return inputChan, newOfflineSplit(fields, inputChan, outputs)
}

func buildOnlineExecution(pipeline *compiled.Pipeline, logger Logger) (*onlineExecution, error) {
	// allChans stores all the channels, including the ones for aux stages.
	// linkChans stores the channels associates with the pipeline links.
	var allChans []chan onlineState

	linkChans := make(map[internal.LinkName]chan onlineState)

	pipeline.VisitLinks(func(l *internal.Link) error {
		ch := make(chan onlineState, defaultChanSize)
		allChans = append(allChans, ch)
		linkChans[l.Name()] = ch
		return nil
	})

	stages := newStageMap()
	pipeline.VisitStages(func(s *compiled.Stage) error {
		name := s.Name()
		inChan, aux, err := buildOnlineInputResources(s, &allChans, linkChans)
		if err != nil {
			return err
		}
		if aux != nil {
			stages.addInputStage(name, aux)
		}
		outChan, aux := buildOnlineOutputResources(s, &allChans, linkChans)
		if aux != nil {
			stages.addOutputStage(name, aux)
		}
		address := s.Address()
		clientBuilder := s.Method().ClientBuilder()
		rpcStage := newOnlineUnary(
			name, inChan, outChan, address, clientBuilder, logger,
		)
		stages.addRpcStage(name, rpcStage)
		return nil
	})
	drainFunc := newChanDrainer(5*time.Millisecond, allChans...)
	return newOnlineExecution(stages, drainFunc, logger), nil
}

func buildOnlineInputResources(
	s *compiled.Stage, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage, error) {
	switch len(s.Inputs()) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		s := newOnlineSource(1, s.Method().Input().EmptyGen(), ch)
		return ch, s, nil
	case 1:
		l := s.Inputs()[0]
		if !l.Target().Field().IsEmpty() {
			output, s := buildOnlineMergeStage(s, allChans, linkChans)
			return output, s, nil
		}
		return linkChans[l.Name()], nil, nil
	default:
		output, s := buildOnlineMergeStage(s, allChans, linkChans)
		return output, s, nil
	}
}

func buildOnlineMergeStage(
	s *compiled.Stage, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(s.Inputs()))
	// channels where the stage will receive the several inputs.
	inputs := make([]<-chan onlineState, 0, len(s.Inputs()))
	// channel where the stage will send the constructed messages.
	outputChan := make(chan onlineState, defaultChanSize)
	*allChans = append(*allChans, outputChan)
	for _, l := range s.Inputs() {
		fields = append(fields, l.Target().Field())
		inputs = append(inputs, linkChans[l.Name()])
	}
	gen := s.Method().Input().EmptyGen()
	return outputChan, newOnlineMerge(fields, inputs, outputChan, gen)
}

func buildOnlineOutputResources(
	s *compiled.Stage, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	switch len(s.Outputs()) {
	case 0:
		ch := make(chan onlineState, defaultChanSize)
		*allChans = append(*allChans, ch)
		return ch, newOnlineSink(ch)
	case 1:
		// We have only one link, but we want a sub message. We can use the
		// split stage with just one output that retrieves the desired message
		// part.
		l := s.Outputs()[0]
		if !l.Source().Field().IsEmpty() {
			return buildOnlineSplitStage(s, allChans, linkChans)
		}
		return linkChans[l.Name()], nil
	default:
		return buildOnlineSplitStage(s, allChans, linkChans)
	}
}

func buildOnlineSplitStage(
	s *compiled.Stage, allChans *[]chan onlineState, linkChans map[internal.LinkName]chan onlineState,
) (chan onlineState, Stage) {
	fields := make([]internal.MessageField, 0, len(s.Outputs()))
	// channel where the stage will send the produced states.
	inputChan := make(chan onlineState, defaultChanSize)
	*allChans = append(*allChans, inputChan)
	// channels to split the received states.
	outputs := make([]chan<- onlineState, 0, len(s.Outputs()))
	for _, l := range s.Outputs() {
		fields = append(fields, l.Source().Field())
		outputs = append(outputs, linkChans[l.Name()])
	}
	return inputChan, newOnlineSplit(fields, inputChan, outputs)
}
