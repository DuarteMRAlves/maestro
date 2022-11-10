package compiled

import (
	"context"
	"errors"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
)

// Context specifies a compilation context to create the pipeline.
type Context struct {
	resolver method.Resolver
}

func NewContext(methodLoader method.Resolver) Context {
	c := Context{}
	c.resolver = methodLoader
	return c
}

var (
	errEmptyPipelineName    = errors.New("empty pipeline name")
	errEmptyStageName       = errors.New("empty stage name")
	errEmptyLinkName        = errors.New("empty link name")
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

type incompatibleMessageDesc struct{ A, B message.Type }

func (err *incompatibleMessageDesc) Error() string {
	return fmt.Sprintf("incompatible message descriptors: %s, %s", err.A, err.B)
}

// New compiles a Pipeline from its specification.
func New(ctx Context, cfg *api.Pipeline) (*Pipeline, error) {
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

func compileStage(ctx Context, cfg *api.Stage) (*Stage, error) {
	name, err := compileStageName(cfg.Name)
	if err != nil {
		return nil, err
	}
	address := compileStageAddr(cfg.Address, cfg.Service, cfg.Method)
	method, err := ctx.resolver.Resolve(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("load method %q: %w", cfg.Address, err)
	}
	stage := &Stage{
		name:    name,
		sType:   StageTypeUnary,
		address: address,
		desc:    method,
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

func compileStageAddr(address, service, method string) string {
	a := address
	s := "*"
	m := "*"
	if service != "" {
		s = service
	}
	if method != "" {
		m = method
	}
	return fmt.Sprintf("%s/%s/%s", a, s, m)
}

func compileLink(cfg *api.Link) (*Link, error) {
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
	size := defaultLinkSize
	if cfg.Size > 0 {
		size = cfg.Size
	}
	l := NewLink(name, source, target, size, cfg.NumEmptyMessages)
	return l, nil
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
	fieldName := message.Field(field)
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

	sourceMsg := source.desc.Output()
	if !link.Source().Field().IsUnspecified() {
		sourceMsg, err = sourceMsg.Subfield(link.Source().Field())
		if err != nil {
			return err
		}
	}

	targetMsg := target.desc.Input()
	if !link.Target().Field().IsUnspecified() {
		targetMsg, err = targetMsg.Subfield(link.Target().Field())
		if err != nil {
			return err
		}
	}
	if !sourceMsg.Compatible(targetMsg) {
		return &incompatibleMessageDesc{A: sourceMsg, B: targetMsg}
	}

	return compatibleWithPreviousLinks(link, target)
}

func compatibleWithPreviousLinks(link *Link, target *Stage) error {
	var err error
	target.RangeInputs(func(prev *Link) bool {
		targetFieldLink := link.Target().Field()
		targetFieldPrev := prev.Target().Field()

		// 1. Target receives entire message from this link but another exists.
		if targetFieldLink.IsUnspecified() {
			err = &linkSetsFullMessage{name: link.Name().Unwrap()}
			return false
		}

		// 2. Target already receives entire message from existing link.
		if targetFieldPrev.IsUnspecified() {
			err = &linkSetsFullMessage{name: prev.Name().Unwrap()}
			return false
		}

		// 3. Target receives same field from both links.
		if targetFieldLink == targetFieldPrev {
			err = &linksSetSameField{
				A:     link.Name().Unwrap(),
				B:     prev.Name().Unwrap(),
				field: string(targetFieldPrev),
			}
			return false
		}
		return true
	})
	return err
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
	l := NewLink(
		LinkName{fmt.Sprintf("%s:aux-source-link", s.name.val)},
		&LinkEndpoint{stage: name},
		&LinkEndpoint{stage: s.name},
		defaultLinkSize,
		0,
	)
	source := &Stage{
		name:  name,
		sType: StageTypeSource,
		// give access to method information for later usage
		address: s.address,
		desc:    s.desc,
		inputs:  []*Link{},
		outputs: []*Link{l},
	}
	s.inputs = []*Link{l}
	return source
}

func compileMergeInput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-merge", s.name.val)}
	l := NewLink(
		LinkName{fmt.Sprintf("%s:aux-merge-link", s.name.val)},
		&LinkEndpoint{stage: name},
		&LinkEndpoint{stage: s.name},
		defaultLinkSize,
		0,
	)
	merge := &Stage{
		name:  name,
		sType: StageTypeMerge,
		// give access to method information for later usage
		address: s.address,
		desc:    s.desc,
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
	l := NewLink(
		LinkName{fmt.Sprintf("%s:aux-sink-link", s.name.val)},
		&LinkEndpoint{stage: s.name},
		&LinkEndpoint{stage: name},
		defaultLinkSize,
		0,
	)
	sink := &Stage{
		name:  name,
		sType: StageTypeSink,
		// give access to method information for later usage
		address: s.address,
		desc:    s.desc,
		inputs:  []*Link{l},
		outputs: []*Link{},
	}
	s.outputs = []*Link{l}
	return sink
}

func compileSplitOutput(s *Stage) *Stage {
	name := StageName{fmt.Sprintf("%s:aux-split", s.name.val)}
	l := NewLink(
		LinkName{fmt.Sprintf("%s:aux-split-link", s.name.val)},
		&LinkEndpoint{stage: s.name},
		&LinkEndpoint{stage: name},
		defaultLinkSize,
		0,
	)
	split := &Stage{
		name:  name,
		sType: StageTypeSplit,
		// give access to method information for later usage
		address: s.address,
		desc:    s.desc,
		inputs:  []*Link{l},
		outputs: s.outputs,
	}
	s.outputs = []*Link{l}
	for _, i := range split.outputs {
		i.source.stage = name
	}
	return split
}

func hasLoops(g stageGraph) bool {
	visited := make(map[StageName]bool, len(g))
	for name := range g {
		visited[name] = false
	}
	for n, curr := range g {
		if !visited[n] {
			inStack := make(map[StageName]bool)
			if visitStage(g, curr, visited, inStack) {
				return true
			}
		}
	}
	return false
}

func visitStage(g stageGraph, curr *Stage, visited map[StageName]bool, inStack map[StageName]bool) bool {
	inStack[curr.name] = true
	visited[curr.name] = true
	for _, o := range curr.outputs {
		neigh := g[o.target.stage]
		if inStack[neigh.name] {
			return true
		}
		if !visited[neigh.name] && visitStage(g, neigh, visited, inStack) {
			return true
		}
	}
	inStack[curr.name] = false
	return false
}
