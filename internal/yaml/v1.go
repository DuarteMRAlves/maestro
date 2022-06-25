package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"

	"github.com/DuarteMRAlves/maestro/internal/arrays"
	"github.com/DuarteMRAlves/maestro/internal/spec"
	"gopkg.in/yaml.v2"
)

const (
	stageKind    = "stage"
	linkKind     = "link"
	pipelineKind = "pipeline"
)

var (
	ErrMissingKind = errors.New("kind not specified")
	ErrEmptySpec   = errors.New("empty spec")
)

type unknownKind struct {
	Kind string
}

func (err *unknownKind) Error() string {
	return fmt.Sprintf("unknown kind '%s'", err.Kind)
}

// ReadV1 reads a set of files in the Maestro V1 format and returns the
// discovered resources.
func ReadV1(files ...string) ([]*spec.Pipeline, error) {
	var pipelines []*spec.Pipeline
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("read v1: %w", err)
		}
		reader := bytes.NewReader(data)

		dec := yaml.NewDecoder(reader)
		dec.SetStrict(true)

		var resources []v1ReadResource
		for {
			var r v1ReadResource
			if err = dec.Decode(&r); err != nil {
				break
			}
			resources = append(resources, r)
		}
		if err != nil && err != io.EOF {
			switch concreteErr := err.(type) {
			case *yaml.TypeError:
				err = typeErrorToError(concreteErr)
				return nil, fmt.Errorf("read v1: %w", err)
			default:
				return nil, fmt.Errorf("read v1: %w", err)
			}
		}
		for _, r := range resources {
			switch r.Kind {
			case pipelineKind:
				o, err := resourceToPipeline(r)
				if err != nil {
					return nil, fmt.Errorf("read v1: %w", err)
				}
				pipelines = append(pipelines, o)
			}
		}
		for _, r := range resources {
			switch r.Kind {
			case stageKind:
				s, n, err := resourceToStage(r)
				if err != nil {
					return nil, fmt.Errorf("read v1: %w", err)
				}
				p := pipelineWithName(pipelines, n)
				if p == nil {
					return nil, fmt.Errorf("read v1: pipeline not found %s", n)
				}
				p.Stages = append(p.Stages, s)
			case linkKind:
				l, n, err := resourceToLink(r)
				if err != nil {
					return nil, fmt.Errorf("read v1: %w", err)
				}
				p := pipelineWithName(pipelines, n)
				if p == nil {
					return nil, fmt.Errorf("read v1: pipeline not found %s", n)
				}
				p.Links = append(p.Links, l)
			}
		}
	}
	return pipelines, nil
}

func pipelineWithName(pipelines []*spec.Pipeline, name string) *spec.Pipeline {
	filterFn := func(p *spec.Pipeline) bool {
		return p.Name == name
	}
	return arrays.FindFirst(filterFn, pipelines...)
}

func resourceToPipeline(r v1ReadResource) (*spec.Pipeline, error) {
	s, ok := r.Spec.(*v1PipelineSpec)
	if !ok {
		return nil, errors.New("pipeline spec cast error")
	}
	mode, err := stringToExecutionMode(s.Mode)
	if err != nil {
		return nil, err
	}
	return &spec.Pipeline{Name: s.Name, Mode: mode}, nil
}

func stringToExecutionMode(val string) (spec.ExecutionMode, error) {
	switch val {
	case "", "Offline":
		return spec.OfflineExecution, nil
	case "Online":
		return spec.OnlineExecution, nil
	default:
		err := fmt.Errorf("unknown execution mode: %s", val)
		return spec.OfflineExecution, err
	}
}

func resourceToStage(r v1ReadResource) (*spec.Stage, string, error) {
	stageSpec, ok := r.Spec.(*v1StageSpec)
	if !ok {
		return nil, "", errors.New("stage spec cast error")
	}

	s := &spec.Stage{
		Name: stageSpec.Name,
		MethodContext: spec.MethodContext{
			Address: stageSpec.Address,
			Service: stageSpec.Service,
			Method:  stageSpec.Method,
		},
	}
	return s, stageSpec.Pipeline, nil
}

func resourceToLink(r v1ReadResource) (*spec.Link, string, error) {
	linkSpec, ok := r.Spec.(*v1LinkSpec)
	if !ok {
		return nil, "", errors.New("link spec cast error")
	}

	l := &spec.Link{
		Name:             linkSpec.Name,
		SourceStage:      linkSpec.SourceStage,
		SourceField:      linkSpec.SourceField,
		TargetStage:      linkSpec.TargetStage,
		TargetField:      linkSpec.TargetField,
		NumEmptyMessages: linkSpec.NumEmptyMessages,
	}
	return l, linkSpec.Pipeline, nil
}

type v1ReadResource struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"-"`
}

func (r *v1ReadResource) String() string {
	return fmt.Sprintf("v1ReadResource{Kind:%v,Spec:%v}", r.Kind, r.Spec)
}

type yamlNode struct {
	unmarshal func(interface{}) error
}

func (n *yamlNode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	n.unmarshal = unmarshal
	return nil
}

// UnmarshalYAML changes the default unmarshalling behaviour for the Resource
// unmarshalling to account for the dynamic unmarshalling of the spec field.
func (r *v1ReadResource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	obj := &struct {
		Kind string `yaml:"kind"`
		// This field will not be decoded into a specific type but the
		// relevant information will be stored.
		Spec yamlNode `yaml:"spec"`
	}{}
	if err := unmarshal(obj); err != nil {
		return err
	}
	r.Kind = obj.Kind
	if r.Kind == "" {
		return ErrMissingKind
	}
	switch r.Kind {
	case stageKind:
		r.Spec = new(v1StageSpec)
	case linkKind:
		r.Spec = new(v1LinkSpec)
	case pipelineKind:
		r.Spec = new(v1PipelineSpec)
	default:
		return &unknownKind{Kind: r.Kind}
	}
	if obj.Spec.unmarshal == nil {
		return ErrEmptySpec
	}
	err := obj.Spec.unmarshal(r.Spec)
	if err != nil {
		return err
	}

	return valV1ReadResource(r)
}

func valV1ReadResource(r *v1ReadResource) error {
	switch r.Kind {
	case stageKind:
		spec, ok := r.Spec.(*v1StageSpec)
		if !ok {
			return errors.New("spec not v1StageSpec for stage kind")
		}
		return valV1StageSpec(spec)
	case linkKind:
		spec, ok := r.Spec.(*v1LinkSpec)
		if !ok {
			return errors.New("spec not v1LinkSpec for link kind")
		}
		return valV1LinkSpec(spec)
	case pipelineKind:
		spec, ok := r.Spec.(*v1PipelineSpec)
		if !ok {
			return errors.New("spec not v1PipelineSpec for pipeline kind")
		}
		return valV1PipelineSpec(spec)
	default:
		return &unknownKind{Kind: r.Kind}
	}
}

func valV1PipelineSpec(spec *v1PipelineSpec) error {
	if spec.Name == "" {
		return &missingRequiredField{Field: "name"}
	}
	return nil
}

func valV1StageSpec(spec *v1StageSpec) error {
	if spec.Name == "" {
		return &missingRequiredField{Field: "name"}
	}
	if spec.Address == "" {
		return &missingRequiredField{Field: "address"}
	}
	if spec.Pipeline == "" {
		return &missingRequiredField{Field: "pipeline"}
	}
	return nil
}

func valV1LinkSpec(spec *v1LinkSpec) error {
	if spec.Name == "" {
		return &missingRequiredField{Field: "name"}
	}
	if spec.SourceStage == "" {
		return &missingRequiredField{Field: "source_stage"}
	}
	if spec.TargetStage == "" {
		return &missingRequiredField{Field: "target_stage"}
	}
	return nil
}

// WriteV1 stores the resources set in a single file as a
func WriteV1(pipeline *spec.Pipeline, file string, perm fs.FileMode) error {
	var (
		buf bytes.Buffer
		err error
	)
	enc := yaml.NewEncoder(&buf)
	err = encodeResources(enc, pipelineToResource, pipeline)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	stageEncFunc := func(r *v1WriteResource, s *spec.Stage) {
		stageToResource(r, s, pipeline.Name)
	}
	err = encodeResources(enc, stageEncFunc, pipeline.Stages...)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	linkEncFunc := func(r *v1WriteResource, l *spec.Link) {
		linkToResource(r, l, pipeline.Name)
	}
	err = encodeResources(enc, linkEncFunc, pipeline.Links...)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	err = ioutil.WriteFile(file, buf.Bytes(), perm)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	return nil
}

func encodeResources[T any](
	enc *yaml.Encoder, encodeFn func(*v1WriteResource, T), resources ...T,
) error {
	for _, r := range resources {
		var w v1WriteResource
		encodeFn(&w, r)
		err := enc.Encode(w)
		if err != nil {
			return err
		}
	}
	return nil
}

type v1WriteResource struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"spec"`
}

func pipelineToResource(r *v1WriteResource, p *spec.Pipeline) {
	var pipelineSpec v1PipelineSpec
	pipelineSpec.Name = p.Name
	switch p.Mode {
	// No need to specify offline as it is the default.
	case spec.OfflineExecution:
		pipelineSpec.Mode = ""
	case spec.OnlineExecution:
		pipelineSpec.Mode = "Online"
	default:
		pipelineSpec.Mode = "Unknown"
	}

	r.Kind = pipelineKind
	r.Spec = pipelineSpec
}

func stageToResource(r *v1WriteResource, s *spec.Stage, pipelineName string) {
	var stageSpec v1StageSpec
	stageSpec.Name = s.Name
	stageSpec.Address = s.MethodContext.Address
	stageSpec.Service = s.MethodContext.Service
	stageSpec.Method = s.MethodContext.Method
	stageSpec.Pipeline = pipelineName

	r.Kind = stageKind
	r.Spec = stageSpec
}

func linkToResource(r *v1WriteResource, l *spec.Link, pipelineName string) {
	var linkSpec v1LinkSpec
	linkSpec.Name = l.Name
	linkSpec.SourceStage = l.SourceStage
	linkSpec.SourceField = l.SourceField
	linkSpec.TargetStage = l.TargetStage
	linkSpec.TargetField = l.TargetField
	linkSpec.NumEmptyMessages = l.NumEmptyMessages
	linkSpec.Pipeline = pipelineName

	r.Kind = linkKind
	r.Spec = linkSpec
}
