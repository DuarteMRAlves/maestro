package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"gopkg.in/yaml.v2"
	"io"
	"io/fs"
	"io/ioutil"
	"reflect"
	"strings"
)

const (
	assetKind    = "asset"
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
			case assetKind:
				a, err := resourceToAsset(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("read v1: %w", err)
				}
				resources.Assets = append(resources.Assets, a)
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
	return Pipeline{Name: name}, nil
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

func resourceToAsset(r v1ReadResource) (Asset, error) {
	spec, ok := r.Spec.(*v1AssetSpec)
	if !ok {
		return Asset{}, errors.New("asset spec cast error")
	}
	name, err := internal.NewAssetName(spec.Name)
	if err != nil {
		return Asset{}, err
	}
	image := internal.NewImage(spec.Image)
	return Asset{Name: name, Image: image}, nil
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
	case assetKind:
		r.Spec = new(v1AssetSpec)
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

	return validateInfo(r.Spec)
}

// validateInfo verifies if all the restrictions specified by the info tags are
// complied
func validateInfo(v interface{}) error {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
	case reflect.Struct:
		// Do nothing, we keep the same value to later analyze as we already
		// have the struct.
	default:
		return fmt.Errorf(
			"invalid type: expected Ptr or Struct but got %v", value.Kind(),
		)
	}

	objType := value.Type()
	for i := 0; i < objType.NumField(); i++ {
		objTypeField := objType.Field(i)
		// Ignore unexported fields
		if !objTypeField.IsExported() {
			continue
		}
		fieldValue := value.Field(i)
		if err := validateField(objTypeField, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func validateField(
	objTypeField reflect.StructField, fieldValue reflect.Value,
) error {
	tag, hasTag := objTypeField.Tag.Lookup("info")
	if hasTag {

		tagOpts := strings.Split(tag, ",")
		for _, opt := range tagOpts {
			switch opt {
			case "required":
				if fieldValue.IsZero() {
					return &missingRequiredField{Field: yamlName(objTypeField)}
				}
			}
		}
	}
	return nil
}

func yamlName(f reflect.StructField) string {
	tag, hasTag := f.Tag.Lookup("yaml")
	if hasTag {
		tagOpts := strings.Split(tag, ",")
		if tagOpts[0] != "" {
			return tagOpts[0]
		}
	}
	return f.Name
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
	err = encodeResources(enc, assetToResource, resources.Assets...)
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

func assetToResource(r *v1WriteResource, a Asset) {
	var spec v1AssetSpec
	spec.Name = a.Name.Unwrap()
	spec.Image = a.Image.Unwrap()

	r.Kind = assetKind
	r.Spec = spec
}
