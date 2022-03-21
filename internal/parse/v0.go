package parse

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// FromV0 parses configuration files for the orchestrator
// https://github.com/DuarteMRAlves/Pipeline-Orchestrator for compatibility
// purposes.
func FromV0(file string) (ResourceSet, error) {
	var (
		resources ResourceSet
		fileSpec  v0FileSpec
	)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
	}

	err = yaml.UnmarshalStrict(data, &fileSpec)
	if err != nil {
		return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
	}

	orchName, err := internal.NewOrchestrationName("v0-orchestration")
	if err != nil {
		return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
	}
	orch := Orchestration{Name: orchName}
	resources.Orchestrations = append(resources.Orchestrations, orch)

	for _, spec := range fileSpec.Stages {
		name, err := internal.NewStageName(spec.Name)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
		}
		addr := internal.NewAddress(fmt.Sprintf("%s:%d", spec.Host, spec.Port))
		serv := internal.NewService(spec.Service)
		meth := internal.NewMethod(spec.Method)

		methCtx := MethodContext{Address: addr, Service: serv, Method: meth}
		s := Stage{Name: name, Method: methCtx, Orchestration: orchName}
		resources.Stages = append(resources.Stages, s)
	}

	for _, spec := range fileSpec.Links {
		name := fmt.Sprintf(
			"v0-link-%s-to-%s", spec.Source.Stage, spec.Target.Stage,
		)
		linkName, err := internal.NewLinkName(name)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
		}

		srcStage, err := internal.NewStageName(spec.Source.Stage)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
		}
		srcField := internal.NewMessageField(spec.Source.Field)

		tgtStage, err := internal.NewStageName(spec.Target.Stage)
		if err != nil {
			return ResourceSet{}, fmt.Errorf("parse v0: %w", err)
		}
		tgtField := internal.NewMessageField(spec.Target.Field)

		l := Link{
			Name:          linkName,
			Source:        LinkEndpoint{Stage: srcStage, Field: srcField},
			Target:        LinkEndpoint{Stage: tgtStage, Field: tgtField},
			Orchestration: orchName,
		}
		resources.Links = append(resources.Links, l)
	}

	return resources, nil
}

type v0FileSpec struct {
	Stages []v0StageSpec `yaml:"stages" info:"required"`
	Links  []v0LinkSpec  `yaml:"links" info:"required"`
}

type v0StageSpec struct {
	// Name that should be associated with the stage.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Host where the server for the stage is running.
	// (required)
	Host string `yaml:"host" info:"required"`
	// Port where the server for the stage is running.
	// (required)
	Port int `yaml:"port" info:"required"`
	// Service specifies the name of the service to call.
	// (optional)
	Service string `yaml:"service"`
	// Method specifies the name of the method to call.
	// (optional)
	Method string `yaml:"method"`
}

type v0LinkSpec struct {
	Source v0LinkEndpoint `yaml:"source" info:"required"`
	Target v0LinkEndpoint `yaml:"target" info:"required"`
}

type v0LinkEndpoint struct {
	// Stage is the name of the stage to connect to.
	// (required)
	Stage string `yaml:"stage" info:"required"`
	// Field of the stage message to use. Don't specify for entire message.
	// (optional)
	Field string `yaml:"field"`
}
