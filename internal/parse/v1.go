package parse

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

const (
	assetKind         = "asset"
	stageKind         = "stage"
	linkKind          = "link"
	orchestrationKind = "orchestration"
)

var (
	MissingKind = errors.New("kind not specified")
	EmptySpec   = errors.New("empty spec")
)

type UnknownKind struct {
	Kind string
}

func (err *UnknownKind) Error() string {
	return fmt.Sprintf("unknown kind '%s'", err.Kind)
}

type MissingRequiredField struct {
	Field string
}

func (err *MissingRequiredField) Error() string {
	return fmt.Sprintf("missing required field '%s'", err.Field)
}

type UnknownFields struct {
	Fields []string
}

func (err *UnknownFields) Error() string {
	return fmt.Sprintf("unknown fields '%s'", strings.Join(err.Fields, ","))
}

// FromV1 parses a set of files in the Maestro V1 format and returns the
// discovered resources.
func FromV1(files ...string) (ResourceSet, error) {
	var resources ResourceSet
	for _, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
		}
		reader := bytes.NewReader(data)

		dec := yaml.NewDecoder(reader)
		dec.SetStrict(true)

		for {
			var r v1Resource
			if err = dec.Decode(&r); err != nil {
				break
			}

			switch r.Kind {
			case orchestrationKind:
				o, err := resourceToOrchestration(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
				}
				resources.Orchestrations = append(resources.Orchestrations, o)
			case stageKind:
				s, err := resourceToStage(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
				}
				resources.Stages = append(resources.Stages, s)
			case linkKind:
				l, err := resourceToLink(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
				}
				resources.Links = append(resources.Links, l)
			case assetKind:
				a, err := resourceToAsset(r)
				if err != nil {
					return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
				}
				resources.Assets = append(resources.Assets, a)
			}
		}
		if err != nil && err != io.EOF {
			switch err.(type) {
			case *yaml.TypeError:
				err = typeErrorToError(err.(*yaml.TypeError))
				return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
			default:
				return ResourceSet{}, fmt.Errorf("parse v1: %w", err)
			}
		}
	}
	return resources, nil
}

func resourceToOrchestration(r v1Resource) (Orchestration, error) {
	spec, ok := r.Spec.(*v1OrchestrationSpec)
	if !ok {
		return Orchestration{}, errors.New("orchestration spec cast error")
	}
	name, err := internal.NewOrchestrationName(spec.Name)
	if err != nil {
		return Orchestration{}, err
	}
	return Orchestration{Name: name}, nil
}

func resourceToStage(r v1Resource) (Stage, error) {
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

	orchName, err := internal.NewOrchestrationName(spec.Orchestration)
	if err != nil {
		return Stage{}, err
	}

	s := Stage{
		Name:          name,
		Method:        MethodContext{Address: addr, Service: serv, Method: meth},
		Orchestration: orchName,
	}
	return s, err
}

func resourceToLink(r v1Resource) (Link, error) {
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

	orchName, err := internal.NewOrchestrationName(spec.Orchestration)
	if err != nil {
		return Link{}, err
	}

	l := Link{
		Name:          name,
		Source:        LinkEndpoint{Stage: srcStage, Field: srcField},
		Target:        LinkEndpoint{Stage: tgtStage, Field: tgtField},
		Orchestration: orchName,
	}
	return l, nil
}

func resourceToAsset(r v1Resource) (Asset, error) {
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

func typeErrorToError(typeErr *yaml.TypeError) error {
	var unknownFields []string
	unknownRegex := regexp.MustCompile(
		`line \d+: field (?P<field>[\w\W_]+) not found in type [\w\W_.]+`,
	)
	for _, errMsg := range typeErr.Errors {
		unknownMatch := unknownRegex.FindStringSubmatch(errMsg)
		if len(unknownMatch) > 0 {
			unknownFields = append(
				unknownFields, unknownMatch[unknownRegex.SubexpIndex("field")],
			)
		}
	}
	if len(unknownFields) > 0 {
		sort.Strings(unknownFields)
		return &UnknownFields{Fields: unknownFields}
	}
	return typeErr
}

type v1Resource struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"-"`
}

func (r *v1Resource) String() string {
	return fmt.Sprintf("v1Resource{Kind:%v,Spec:%v}", r.Kind, r.Spec)
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
func (r *v1Resource) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
	case orchestrationKind:
		r.Spec = new(v1OrchestrationSpec)
	case assetKind:
		r.Spec = new(v1AssetSpec)
	default:
		return &UnknownKind{Kind: r.Kind}
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
					return &MissingRequiredField{Field: yamlName(objTypeField)}
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
