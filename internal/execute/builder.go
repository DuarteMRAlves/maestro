package execute

import (
	"errors"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/message"
)

const defaultChanSize = 10

type Builder func(pipeline *compiled.Pipeline) (Execution, error)

func NewBuilder(logger Logger) Builder {
	return func(pipeline *compiled.Pipeline) (Execution, error) {
		return buildOfflineExecution(pipeline, logger)
	}
}

func buildOfflineExecution(pipeline *compiled.Pipeline, logger Logger) (*offlineExecution, error) {
	// allChans stores all the channels, including the ones for aux stages.
	// linkChans stores the channels associates with the pipeline links.
	var allChans []chan offlineState

	chans := make(map[compiled.LinkName]chan offlineState)

	err := pipeline.VisitLinks(func(l *compiled.Link) error {
		ch := make(chan offlineState, defaultChanSize)
		allChans = append(allChans, ch)
		chans[l.Name()] = ch
		return nil
	})
	if err != nil {
		return nil, err
	}

	stages := make(map[compiled.StageName]Stage)
	err = pipeline.VisitStages(func(s *compiled.Stage) error {
		execStage, err := buildOfflineStage(s, chans, logger)
		if err != nil {
			return fmt.Errorf("build stage: %w", err)
		}
		stages[s.Name()] = execStage
		return nil
	})
	if err != nil {
		return nil, err
	}

	initOfflineChans(pipeline, chans)

	return newOfflineExecution(stages, logger), nil
}

func buildOfflineStage(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState, l Logger) (Stage, error) {
	switch s.Type() {
	case compiled.StageTypeUnary:
		s, err := buildOfflineUnary(s, chans, l)
		if err != nil {
			return nil, fmt.Errorf("build unary: %w", err)
		}
		return s, nil
	case compiled.StageTypeSource:
		s, err := buildOfflineSource(s, chans)
		if err != nil {
			return nil, fmt.Errorf("build source: %w", err)
		}
		return s, nil
	case compiled.StageTypeSink:
		s, err := buildOfflineSink(s, chans)
		if err != nil {
			return nil, fmt.Errorf("build sink: %w", err)
		}
		return s, nil
	case compiled.StageTypeMerge:
		s, err := buildOfflineMerge(s, chans)
		if err != nil {
			return nil, fmt.Errorf("build merge: %w", err)
		}
		return s, nil
	case compiled.StageTypeSplit:
		s, err := buildOfflineSplit(s, chans)
		if err != nil {
			return nil, fmt.Errorf("build split: %w", err)
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unknown stage type: %s", s.Type())
	}
}

func buildOfflineUnary(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState, l Logger) (Stage, error) {
	name := s.Name()
	inputs := s.CopyInputs()
	outputs := s.CopyOutputs()

	if len(inputs) != 1 {
		return nil, fmt.Errorf("inputs size mismatch: expected 1, actual %d", len(inputs))
	}
	if len(outputs) != 1 {
		return nil, fmt.Errorf("outputs size mismatch: expected 1, actual %d", len(outputs))
	}
	inChan, exists := chans[inputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown input link name: %s", inputs[0].Name())
	}
	outChan, exists := chans[outputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown output link name: %s", outputs[0].Name())
	}
	dialer := s.Dialer()
	if dialer == nil {
		return nil, errors.New("nil dialer")
	}
	return newOfflineUnary(name, inChan, outChan, dialer, l), nil
}

func buildOfflineSource(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState) (Stage, error) {
	input := s.InputDesc()
	if input == nil {
		return nil, errors.New("nil method input")
	}

	inputs := s.CopyInputs()
	outputs := s.CopyOutputs()
	if len(inputs) != 0 {
		return nil, fmt.Errorf("inputs size mismatch: expected 0, actual %d", len(inputs))
	}
	if len(outputs) != 1 {
		return nil, fmt.Errorf("outputs size mismatch: expected 1, actual %d", len(outputs))
	}

	outChan, exists := chans[outputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown output link name: %s", outputs[0].Name())
	}
	return newOfflineSource(message.BuildFunc(input.Build), outChan), nil
}

func buildOfflineSink(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState) (Stage, error) {
	inputs := s.CopyInputs()
	outputs := s.CopyOutputs()
	if len(inputs) != 1 {
		return nil, fmt.Errorf("inputs size mismatch: expected 1, actual %d", len(inputs))
	}
	if len(outputs) != 0 {
		return nil, fmt.Errorf("outputs size mismatch: expected 0, actual %d", len(outputs))
	}
	inChan, exists := chans[inputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown input link name: %s", inputs[0].Name())
	}
	return newOfflineSink(inChan), nil
}

func buildOfflineMerge(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState) (Stage, error) {
	inputs := s.CopyInputs()
	fields := make([]message.Field, 0, len(inputs))
	// channels where the stage will receive the several inputs.
	inChans := make([]<-chan offlineState, 0, len(inputs))
	for _, l := range inputs {
		fields = append(fields, l.Target().Field())
		inChan, exists := chans[l.Name()]
		if !exists {
			return nil, fmt.Errorf("unknown input link name: %s", l.Name())
		}
		inChans = append(inChans, inChan)
	}

	outputs := s.CopyOutputs()
	if len(outputs) != 1 {
		return nil, fmt.Errorf("outputs size mismatch: expected 1, actual %d", len(outputs))
	}
	outChan, exists := chans[outputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown output link name: %s", outputs[0].Name())
	}

	input := s.InputDesc()
	if input == nil {
		return nil, errors.New("nil method input")
	}
	return newOfflineMerge(fields, inChans, outChan, message.BuildFunc(input.Build)), nil
}

func buildOfflineSplit(s *compiled.Stage, chans map[compiled.LinkName]chan offlineState) (Stage, error) {
	inputs := s.CopyInputs()
	if len(inputs) != 1 {
		return nil, fmt.Errorf("inputs size mismatch: expected 1, actual %d", len(inputs))
	}
	inChan, exists := chans[inputs[0].Name()]
	if !exists {
		return nil, fmt.Errorf("unknown input link name: %s", inputs[0].Name())
	}

	outputs := s.CopyOutputs()
	fields := make([]message.Field, 0, len(outputs))
	// channels to split the received states.
	outChans := make([]chan<- offlineState, 0, len(outputs))
	for _, l := range outputs {
		fields = append(fields, l.Source().Field())
		outChan, exists := chans[l.Name()]
		if !exists {
			return nil, fmt.Errorf("unknown output link name: %s", l.Name())
		}
		outChans = append(outChans, outChan)
	}
	return newOfflineSplit(fields, inChan, outChans), nil
}

func initOfflineChans(
	pipeline *compiled.Pipeline, chans map[compiled.LinkName]chan offlineState,
) error {
	return pipeline.VisitLinks(func(l *compiled.Link) error {
		if l.NumEmptyMessages() == 0 {
			return nil
		}
		if l.NumEmptyMessages() > defaultChanSize {
			return fmt.Errorf("link %q: too many empty messages", l.Name())
		}
		ch, ok := chans[l.Name()]
		if !ok {
			return fmt.Errorf("link %q: chan not found", l.Name())
		}
		src, ok := pipeline.Stage(l.Source().Stage())
		if !ok {
			return fmt.Errorf("link %q: source %q not found", l.Name(), l.Source().Stage())
		}
		msgType := src.InputDesc()
		if !l.Source().Field().IsUnspecified() {
			var err error
			msgType, err = msgType.Subfield(l.Source().Field())
			if err != nil {
				return err
			}
		}
		for i := 0; i < int(l.NumEmptyMessages()); i++ {
			ch <- newOfflineState(msgType.Build())
		}
		return nil
	})
}
