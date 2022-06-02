package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"

	"github.com/DuarteMRAlves/maestro/internal"
	"gopkg.in/yaml.v2"
)

const (
	stageKind    = "stage"
	linkKind     = "link"
	pipelineKind = "pipeline"
)

var (
	MissingKind = errors.New("kind not specified")
	EmptySpec   = errors.New("empty spec")
)

type unknownKind struct {
	Kind string
}

func (err *unknownKind) Error() string {
	return fmt.Sprintf("unknown kind '%s'", err.Kind)
}

// ReadV1 reads a set of files in the Maestro V1 format and returns the
// discovered resources.
func ReadV1(files ...string) (ResourceSet, error) {
	var resources ResourceSet
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("read v1: %w", err)
		}
		reader := bytes.NewReader(data)

		dec := yaml.NewDecoder(reader)
		dec.SetStrict(true)

		for {
			var r v1ReadResource
			if err = dec.Decode(&r); err != nil {
				break
			}

			switch r.Kind {
			case pipelineKind:
				o, err := resourceToPipeline(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("read v1: %w", err)
				}
				resources.Pipelines = append(resources.Pipelines, o)
			case stageKind:
				s, err := resourceToStage(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("read v1: %w", err)
				}
				resources.Stages = append(resources.Stages, s)
			case linkKind:
				l, err := resourceToLink(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("read v1: %w", err)
				}
				resources.Links = append(resources.Links, l)
			}
		}
		if err != nil && err != io.EOF {
			switch err.(type) {
			case *yaml.TypeError:
				err = typeErrorToError(err.(*yaml.TypeError))
				return ResourceSet{}, fmt.Errorf("read v1: %w", err)
			default:
				return ResourceSet{}, fmt.Errorf("read v1: %w", err)
			}
		}
	}
	return resources, nil
}

func resourceToPipeline(r v1ReadResource) (Pipeline, error) {
	spec, ok := r.Spec.(*v1PipelineSpec)
	if !ok {
		return Pipeline{}, errors.New("pipeline spec cast error")
	}
	name, err := internal.NewPipelineName(spec.Name)
	if err != nil {
		return Pipeline{}, err
	}
	mode, err := stringToExecutionMode(spec.Mode)
	if err != nil {
		return Pipeline{}, err
	}
	return Pipeline{Name: name, Mode: mode}, nil
}

func stringToExecutionMode(val string) (internal.ExecutionMode, error) {
	switch val {
	case "", "Offline":
		return internal.OfflineExecution, nil
	case "Online":
		return internal.OnlineExecution, nil
	default:
		err := fmt.Errorf("unknown execution mode: %s", val)
		return internal.ExecutionMode{}, err
	}
}

func resourceToStage(r v1ReadResource) (Stage, error) {
	spec, ok := r.Spec.(*v1StageSpec)
	if !ok {
		return Stage{}, errors.New("stage spec cast error")
	}

	name, err := internal.NewStageName(spec.Name)
	if err != nil {
		return Stage{}, err
	}
	addr := internal.NewAddress(spec.Address)
	serv := internal.NewService(spec.Service)
	meth := internal.NewMethod(spec.Method)

	pipelineName, err := internal.NewPipelineName(spec.Pipeline)
	if err != nil {
		return Stage{}, err
	}

	s := Stage{
		Name:     name,
		Method:   MethodContext{Address: addr, Service: serv, Method: meth},
		Pipeline: pipelineName,
	}
	return s, err
}

func resourceToLink(r v1ReadResource) (Link, error) {
	spec, ok := r.Spec.(*v1LinkSpec)
	if !ok {
		return Link{}, errors.New("link spec cast error")
	}

	name, err := internal.NewLinkName(spec.Name)
	if err != nil {
		return Link{}, err
	}

	srcStage, err := internal.NewStageName(spec.SourceStage)
	if err != nil {
		return Link{}, err
	}
	srcField := internal.NewMessageField(spec.SourceField)

	tgtStage, err := internal.NewStageName(spec.TargetStage)
	if err != nil {
		return Link{}, err
	}
	tgtField := internal.NewMessageField(spec.TargetField)

	pipelineName, err := internal.NewPipelineName(spec.Pipeline)
	if err != nil {
		return Link{}, err
	}

	l := Link{
		Name:     name,
		Source:   LinkEndpoint{Stage: srcStage, Field: srcField},
		Target:   LinkEndpoint{Stage: tgtStage, Field: tgtField},
		Pipeline: pipelineName,
	}
	return l, nil
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
		return MissingKind
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
		return EmptySpec
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
func WriteV1(resources ResourceSet, file string, perm fs.FileMode) error {
	var (
		buf bytes.Buffer
		err error
	)
	enc := yaml.NewEncoder(&buf)
	err = encodeResources(enc, pipelineToResource, resources.Pipelines...)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	err = encodeResources(enc, stageToResource, resources.Stages...)
	if err != nil {
		return fmt.Errorf("write v1: %w", err)
	}
	err = encodeResources(enc, linkToResource, resources.Links...)
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

func pipelineToResource(r *v1WriteResource, o Pipeline) {
	var spec v1PipelineSpec
	spec.Name = o.Name.Unwrap()
	switch o.Mode {
	// No need to specify offline as it is the default.
	case internal.OfflineExecution:
		spec.Mode = ""
	case internal.OnlineExecution:
		spec.Mode = "Online"
	default:
		spec.Mode = "Unknown"
	}

	r.Kind = pipelineKind
	r.Spec = spec
}

func stageToResource(r *v1WriteResource, s Stage) {
	var spec v1StageSpec
	spec.Name = s.Name.Unwrap()
	spec.Address = s.Method.Address.Unwrap()
	spec.Service = s.Method.Service.Unwrap()
	spec.Method = s.Method.Method.Unwrap()
	spec.Pipeline = s.Pipeline.Unwrap()

	r.Kind = stageKind
	r.Spec = spec
}

func linkToResource(r *v1WriteResource, l Link) {
	var spec v1LinkSpec
	spec.Name = l.Name.Unwrap()
	spec.SourceStage = l.Source.Stage.Unwrap()
	spec.SourceField = l.Source.Field.Unwrap()
	spec.TargetStage = l.Target.Stage.Unwrap()
	spec.TargetField = l.Target.Field.Unwrap()
	spec.Pipeline = l.Pipeline.Unwrap()

	r.Kind = linkKind
	r.Spec = spec
}
