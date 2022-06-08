package compiled

import (
	"errors"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/spec"
)

// MethodLoader resolves a method from its context.
type MethodLoader interface {
	Load(*MethodContext) (UnaryMethod, error)
}

// MethodLoaderFunc is an adapter to use functions as MethodLoader objects.
type MethodLoaderFunc func(*MethodContext) (UnaryMethod, error)

func (fn MethodLoaderFunc) Load(methodContext *MethodContext) (UnaryMethod, error) {
	return fn(methodContext)
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

type unsupportedExecutionMode struct{ mode spec.ExecutionMode }

func (err *unsupportedExecutionMode) Error() string {
	return fmt.Sprintf("unsupported execution mode: %s", err.mode)
}

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
func New(ctx Context, pipelineSpec *spec.Pipeline) (*Pipeline, error) {
	name, err := compileName(pipelineSpec.Name)
	if err != nil {
		return nil, err
	}

	mode, err := compileExecutionMode(pipelineSpec.Mode)
	if err != nil {
		return nil, err
	}

	stages := make(stageGraph, len(pipelineSpec.Stages))
	for _, stageSpec := range pipelineSpec.Stages {
		stageName := stageSpec.Name
		stage, err := compileStage(ctx, stageSpec)
		if err != nil {
			return nil, fmt.Errorf("compile stage '%s': %w", stageName, err)
		}
		err = validateStage(stages, stage)
		if err != nil {
			return nil, fmt.Errorf("validate stage '%s': %w", stageName, err)
		}
		stages[stage.name] = stage
	}

	for _, linkSpec := range pipelineSpec.Links {
		linkName := linkSpec.Name
		link, err := compileLink(linkSpec)
		if err != nil {
			return nil, fmt.Errorf("compile link '%s': %w", linkName, err)
		}
		err = validateLink(stages, link)
		if err != nil {
			return nil, fmt.Errorf("validate link '%s': %w", linkName, err)
		}

		source := stages[link.Source().Stage()]
		target := stages[link.Target().Stage()]

		target.inputs = append(target.inputs, link)
		source.outputs = append(source.outputs, link)
	}

	p := &Pipeline{
		name:   name,
		mode:   mode,
		stages: stages,
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

func compileExecutionMode(mode spec.ExecutionMode) (ExecutionMode, error) {
	var executionMode ExecutionMode
	switch mode {
	case spec.OfflineExecution:
		executionMode = OfflineExecution
	case spec.OnlineExecution:
		executionMode = OnlineExecution
	default:
		return executionMode, &unsupportedExecutionMode{mode: mode}
	}
	return executionMode, nil
}

func compileStage(ctx Context, stageSpec *spec.Stage) (*Stage, error) {
	name, err := compileStageName(stageSpec.Name)
	if err != nil {
		return nil, err
	}
	address, err := compileAddress(stageSpec.MethodContext.Address)
	if err != nil {
		return nil, err
	}
	methodCtx, err := compileMethodContext(&stageSpec.MethodContext)
	if err != nil {
		return nil, err
	}
	unaryMethod, err := ctx.methodLoader.Load(methodCtx)
	if err != nil {
		return nil, err
	}
	stage := &Stage{
		name:    name,
		address: address,
		method:  unaryMethod,
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

func compileMethodContext(
	methodContextSpec *spec.MethodContext,
) (*MethodContext, error) {
	address, err := compileAddress(methodContextSpec.Address)
	if err != nil {
		return nil, err
	}
	service := NewService(methodContextSpec.Service)
	method := NewMethod(methodContextSpec.Method)
	methodCtx := NewMethodContext(address, service, method)
	return &methodCtx, nil
}

func compileAddress(address string) (Address, error) {
	addr := NewAddress(address)
	if addr.IsEmpty() {
		return addr, errEmptyAddress
	}
	return addr, nil
}

func compileLink(linkSpec *spec.Link) (*Link, error) {
	name, err := compileLinkName(linkSpec.Name)
	if err != nil {
		return nil, err
	}
	source, err := compileEndpoint(linkSpec.SourceStage, linkSpec.SourceField)
	if err != nil {
		return nil, err
	}
	if source.Stage().IsEmpty() {
		return nil, errEmptySourceName
	}
	target, err := compileEndpoint(linkSpec.TargetStage, linkSpec.TargetField)
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
	fieldName := NewMessageField(field)
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

	sourceMsg := source.Method().Output()
	if !link.Source().Field().IsEmpty() {
		sourceMsg, err = sourceMsg.GetField(link.Source().Field())
		if err != nil {
			return err
		}
	}

	targetMsg := target.Method().Input()
	if !link.Target().Field().IsEmpty() {
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
		if targetFieldLink.IsEmpty() {
			return &linkSetsFullMessage{name: link.Name().Unwrap()}
		}
		// 2. Target already receives entire message from existing link.
		if targetFieldPrev.IsEmpty() {
			return &linkSetsFullMessage{name: prev.Name().Unwrap()}
		}
		// 3. Target receives same field from both links.
		if targetFieldLink.Unwrap() == targetFieldPrev.Unwrap() {
			return &linksSetSameField{
				A:     link.Name().Unwrap(),
				B:     prev.Name().Unwrap(),
				field: targetFieldPrev.Unwrap(),
			}
		}
	}
	return nil
}
