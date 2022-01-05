package resources

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

const (
	assetKind         = "asset"
	stageKind         = "stage"
	linkKind          = "link"
	orchestrationKind = "orchestration"
)

type Resource struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"-"`
}

func (r *Resource) IsValidKind() bool {
	return r.IsAssetKind() ||
		r.IsStageKind() ||
		r.IsLinkKind() ||
		r.IsOrchestrationKind()
}

func (r *Resource) IsAssetKind() bool {
	return r.Kind == assetKind
}

func (r *Resource) IsStageKind() bool {
	return r.Kind == stageKind
}

func (r *Resource) IsLinkKind() bool {
	return r.Kind == linkKind
}

func (r *Resource) IsOrchestrationKind() bool {
	return r.Kind == orchestrationKind
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource{Kind:%v,Spec:%v}", r.Kind, r.Spec)
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
func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
		return errdefs.InvalidArgumentWithMsg("kind not specified")
	}
	switch r.Kind {
	case assetKind:
		r.Spec = new(AssetSpec)
	case stageKind:
		r.Spec = new(types.Stage)
	case linkKind:
		r.Spec = new(LinkSpec)
	case orchestrationKind:
		r.Spec = new(OrchestrationSpec)
	default:
		return errdefs.InvalidArgumentWithMsg("unknown kind: '%v'", r.Kind)
	}
	if obj.Spec.unmarshal == nil {
		return errdefs.InvalidArgumentWithMsg("empty spec")
	}
	err := obj.Spec.unmarshal(r.Spec)
	if err != nil {
		return err
	}

	return validateInfo(r.Spec)
}
