package compiled

import (
	"errors"
	"fmt"
)

// MethodLoader resolves a method from its identifier.
type MethodLoader interface {
	Load(MethodID) (MethodDesc, error)
}

// MethodLoaderFunc is an adapter to use functions as MethodLoader objects.
type MethodLoaderFunc func(MethodID) (MethodDesc, error)

func (fn MethodLoaderFunc) Load(mid MethodID) (MethodDesc, error) {
	return fn(mid)
}

// Context specifies a compilation context to create the pipeline.
type Context struct {
	methodLoader MethodLoader
}

func NewContext(methodLoader MethodLoader) Context {
	c := Context{}
	c.methodLoader = methodLoader
	return c
}

var (
	errEmptyPipelineName    = errors.New("empty pipeline name")
	errEmptyStageName       = errors.New("empty stage name")
	errEmptyLinkName        = errors.New("empty link name")
	errEmptyAddress         = errors.New("empty address")
	errEmptySourceName      = errors.New("empty source name")
	errEmptyTargetName      = errors.New("empty target name")
	errEqualSourceAndTarget = errors.New("equal source and target stages")
)

type stageAlreadyExists struct{ name string }

func (err *stageAlreadyExists) Error() string {
	return fmt.Sprintf("stage '%s' already exists", err.name)
}

type stageNotFound struct{ name string }

func (err *stageNotFound) Error() string {
	return fmt.Sprintf("stage '%s' not found", err.name)
}

type linkSetsFullMessage struct {
	name string
}

func (err *linkSetsFullMessage) Error() string {
	format := "link '%s' sets full message (incompatible with any other links)"
	return fmt.Sprintf(format, err.name)
}

type linksSetSameField struct{ A, B, field string }

func (err *linksSetSameField) Error() string {
	return fmt.Sprintf("links '%s' and '%s' set same field '%s'", err.A, err.B, err.field)
}

type incompatibleMessageDesc struct{ A, B MessageDesc }

func (err *incompatibleMessageDesc) Error() string {
	return fmt.Sprintf("incompatible message descriptors: %s, %s", err.A, err.B)
}

// New compiles a Pipeline from its specification.
func New(ctx Context, cfg *PipelineConfig) (*Pipeline, error) {
	name, err := compileName(cfg.Name)
	if err != nil {
		return nil, err
	}

	// condensed graph contains rpc stages with multiple inputs and outputs.
	condensedGraph := make(stageGraph, len(cfg.Stages))
	for _, stageCfg := range cfg.Stages {
		stageName := stageCfg.Name
		stage, err := compileStage(ctx, stageCfg)
		if err != nil {
			return nil, fmt.Errorf("compile stage '%s': %w", stageName, err)
		}
		err = validateStage(condensedGraph, stage)
		if err != nil {
			return nil, fmt.Errorf("validate stage '%s': %w", stageName, err)
		}
		condensedGraph[stage.name] = stage
	}

	for _, linkCfg := range cfg.Links {
		linkName := linkCfg.Name
		link, err := compileLink(linkCfg)
		if err != nil {
			return nil, fmt.Errorf("compile link '%s': %w", linkName, err)
		}
		err = validateLink(condensedGraph, link)
		if err != nil {
			return nil, fmt.Errorf("validate link '%s': %w", linkName, err)
		}

		source := condensedGraph[link.Source().Stage()]
		target := condensedGraph[link.Target().Stage()]

		target.inputs = append(target.inputs, link)
		source.outputs = append(source.outputs, link)
	}

	augmentedGraph := augmentedGraphFromCondensed(condensedGraph)

	p := &Pipeline{
		name:   name,
		mode:   cfg.Mode,
		stages: augmentedGraph,
	}
	return p, nil
}

func compileName(name string) (PipelineName, error) {
	pipelineName, err := NewPipelineName(name)
	if err != nil {
		return pipelineName, err
	}
	if pipelineName.IsEmpty() {
		return pipelineName, errEmptyPipelineName
	}
	return pipelineName, nil
}

func compileStage(ctx Context, cfg *StageConfig) (*Stage, error) {
	name, err := compileStageName(cfg.Name)
	if err != nil {
		return nil, err
	}
	method, err := ctx.methodLoader.Load(cfg.MethodID)
	if err != nil {
		return nil, fmt.Errorf("load method id %s: %w", cfg.MethodID, err)
	}
	stage := &Stage{
		name:    name,
		sType:   StageTypeUnary,
		mid:     cfg.MethodID,
		method:  method,
		inputs:  []*Link{},
		outputs: []*Link{},
	}
	return stage, nil
}

func compileStageName(name string) (StageName, error) {
	stageName, err := NewStageName(name)
	if err != nil {
		return stageName, err
	}
	if stageName.IsEmpty() {
		return stageName, errEmptyStageName
	}
	return stageName, nil
}

func compileLink(cfg *LinkConfig) (*Link, error) {
	name, err := compileLinkName(cfg.Name)
	if err != nil {
		return nil, err
	}
	source, err := compileEndpoint(cfg.SourceStage, cfg.SourceField)
	if err != nil {
		return nil, err
	}
	if source.Stage().IsEmpty() {
		return nil, errEmptySourceName
	}
	target, err := compileEndpoint(cfg.TargetStage, cfg.TargetField)
	if err != nil {
		return nil, err
	}
	if target.Stage().IsEmpty() {
		return nil, errEmptyTargetName
	}
	l := NewLink(name, source, target)
	return &l, nil
}

func compileLinkName(name string) (LinkName, error) {
	linkName, err := NewLinkName(name)
	if err != nil {
		return linkName, err
	}
	if linkName.IsEmpty() {
		return linkName, errEmptyLinkName
	}
	return linkName, nil
}

func compileEndpoint(stage string, field string) (*LinkEndpoint, error) {
	var endpt LinkEndpoint
	stageName, err := NewStageName(stage)
	if err != nil {
		return nil, err
	}
	fieldName := MessageField(field)
	endpt = NewLinkEndpoint(StageName(stageName), fieldName)
	return &endpt, nil
}

func validateStage(stages stageGraph, stage *Stage) error {
	_, exists := stages[stage.name]
	if exists {
		return &stageAlreadyExists{name: stage.name.Unwrap()}
	}
	return nil
}

func validateLink(stages stageGraph, link *Link) error {
	var err error

	sourceStage := link.Source().Stage()
	source, exists := stages[sourceStage]
	if !exists {
		return &stageNotFound{name: sourceStage.Unwrap()}
	}

	targetStage := link.Target().Stage()
	target, exists := stages[targetStage]
	if !exists {
		return &stageNotFound{name: targetStage.Unwrap()}
	}

	if sourceStage == targetStage {
		return errEqualSourceAndTarget
	}

	sourceMsg := source.method.Output()
	if !link.Source().Field().IsUnspecified() {
		sourceMsg, err = sourceMsg.GetField(link.Source().Field())
		if err != nil {
			return err
		}
	}

	targetMsg := target.method.Input()
	if !link.Target().Field().IsUnspecified() {
		targetMsg, err = targetMsg.GetField(link.Target().Field())
		if err != nil {
			return err
		}
	}
	if !sourceMsg.Compatible(targetMsg) {
		return &incompatibleMessageDesc{A: sourceMsg, B: targetMsg}
	}

	for _, prev := range target.Inputs() {
		targetFieldLink := link.Target().Field()
		targetFieldPrev := prev.Target().Field()
		// Target receives entire message from this link but another exists.
		if targetFieldLink.IsUnspecified() {
			return &linkSetsFullMessage{name: link.Name().Unwrap()}
		}
		// 2. Target already receives entire message from existing link.
		if targetFieldPrev.IsUnspecified() {
			return &linkSetsFullMessage{name: prev.Name().Unwrap()}
		}
		// 3. Target receives same field from both links.
		if targetFieldLink == targetFieldPrev {
			return &linksSetSameField{
				A:     link.Name().Unwrap(),
				B:     prev.Name().Unwrap(),
				field: string(targetFieldPrev),
			}
		}
	}
	return nil
}

func augmentedGraphFromCondensed(condensedGraph stageGraph) stageGraph {
	augmentedGraph := make(stageGraph)

	for _, s := range condensedGraph {
		if auxInput := compileAuxInputIfNecessary(s); auxInput != nil {
			augmentedGraph[auxInput.name] = auxInput
		}
		if auxOutput := compileAuxOutputIfNecessary(s); auxOutput != nil {
			augmentedGraph[auxOutput.name] = auxOutput
		}
		augmentedGraph[s.name] = s
	}
	return augmentedGraph
}

func compileAuxInputIfNecessary(s *Stage) *Stage {
	switch len(s.inputs) {
	// Stage s has no inputs and is a source of the pipeline. We create
	// a source stage and add it as an input to s.
	case 0:
		return compileSourceInput(s)
	case 1:
		l := s.inputs[0]
		// We only have one link but we have to set a field and so
		// we use a single link merge stage.
		if !l.Target().Field().IsUnspecified() {
			return compileMergeInput(s)
		}
		return nil
	default:
		return compileMergeInput(s)
	}
}

func compileSourceInput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-source", s.name.val)}
	l := &Link{
		name:   LinkName{fmt.Sprintf("%s:aux-source-link", s.name.val)},
		source: &LinkEndpoint{stage: name},
		target: &LinkEndpoint{stage: s.name},
	}
	source := &Stage{
		name:  name,
		sType: StageTypeSource,
		// give access to method information for later usage
		mid:     s.mid,
		method:  s.method,
		inputs:  []*Link{},
		outputs: []*Link{l},
	}
	s.inputs = []*Link{l}
	return source
}

func compileMergeInput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-merge", s.name.val)}
	l := &Link{
		name:   LinkName{fmt.Sprintf("%s:aux-merge-link", s.name.val)},
		source: &LinkEndpoint{stage: name},
		target: &LinkEndpoint{stage: s.name},
	}
	merge := &Stage{
		name:  name,
		sType: StageTypeMerge,
		// give access to method information for later usage
		mid:     s.mid,
		method:  s.method,
		inputs:  s.inputs,
		outputs: []*Link{l},
	}
	s.inputs = []*Link{l}
	for _, i := range merge.inputs {
		i.target.stage = name
	}
	return merge
}

func compileAuxOutputIfNecessary(s *Stage) *Stage {
	switch len(s.outputs) {
	// Stage s has no inputs and is a source of the pipeline. We create
	// a source stage and add it as an input to s.
	case 0:
		return compileSinkOutput(s)
	case 1:
		l := s.inputs[0]
		if !l.Target().Field().IsUnspecified() {
			return compileSplitOutput(s)
		}
		return nil
	default:
		return compileSplitOutput(s)
	}
}

func compileSinkOutput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-sink", s.name.val)}
	l := &Link{
		name:   LinkName{fmt.Sprintf("%s:aux-sink-link", s.name.val)},
		source: &LinkEndpoint{stage: s.name},
		target: &LinkEndpoint{stage: name},
	}
	sink := &Stage{
		name:  name,
		sType: StageTypeSink,
		// give access to method information for later usage
		mid:     s.mid,
		method:  s.method,
		inputs:  []*Link{l},
		outputs: []*Link{},
	}
	s.outputs = []*Link{l}
	return sink
}

func compileSplitOutput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-split", s.name.val)}
	l := &Link{
		name:   LinkName{fmt.Sprintf("%s:aux-split-link", s.name.val)},
		source: &LinkEndpoint{stage: s.name},
		target: &LinkEndpoint{stage: name},
	}
	split := &Stage{
		name:  name,
		sType: StageTypeSplit,
		// give access to method information for later usage
		mid:     s.mid,
		method:  s.method,
		inputs:  []*Link{l},
		outputs: s.outputs,
	}
	s.outputs = []*Link{l}
	for _, i := range split.outputs {
		i.source.stage = name
	}
	return split
}
