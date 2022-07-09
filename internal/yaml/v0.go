package yaml

import (
	"fmt"
	"io/ioutil"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"gopkg.in/yaml.v2"
)

// ReadV0 reads configuration files for the orchestrator
// https://github.com/DuarteMRAlves/Pipeline-Orchestrator for compatibility
// purposes.
func ReadV0(file string) (*api.Pipeline, error) {
	var fileSpec v0FileSpec

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read v0: %w", err)
	}

	err = yaml.UnmarshalStrict(data, &fileSpec)
	if err != nil {
		if typeErr, ok := err.(*yaml.TypeError); ok {
			err = typeErrorToError(typeErr)
		}
		return nil, fmt.Errorf("read v0: %w", err)
	}
	err = valV0FileSpec(fileSpec)
	if err != nil {
		return nil, fmt.Errorf("read v0: %w", err)
	}

	// v0 pipeline executes always online.
	pipeline := &api.Pipeline{Name: "v0-pipeline", Mode: api.OnlineExecution}

	for _, stageSpec := range fileSpec.Stages {
		s := &api.Stage{
			Name:    stageSpec.Name,
			Address: fmt.Sprintf("%s:%d", stageSpec.Host, stageSpec.Port),
			Service: stageSpec.Service,
			Method:  stageSpec.Method,
		}
		pipeline.Stages = append(pipeline.Stages, s)
	}

	for _, linkSpec := range fileSpec.Links {
		name := fmt.Sprintf(
			"v0-link-%s-to-%s", linkSpec.Source.Stage, linkSpec.Target.Stage,
		)

		l := &api.Link{
			Name:        name,
			SourceStage: linkSpec.Source.Stage,
			SourceField: linkSpec.Source.Field,
			TargetStage: linkSpec.Target.Stage,
			TargetField: linkSpec.Target.Field,
		}
		pipeline.Links = append(pipeline.Links, l)
	}

	return pipeline, nil
}

func valV0FileSpec(spec v0FileSpec) error {
	if spec.Stages == nil {
		return &missingRequiredField{Field: "stages"}
	}
	if spec.Links == nil {
		return &missingRequiredField{Field: "links"}
	}
	for _, s := range spec.Stages {
		if err := valV0StageSpec(s); err != nil {
			return err
		}
	}
	for _, l := range spec.Links {
		if err := valV0LinkSpec(l); err != nil {
			return err
		}
	}
	return nil
}

func valV0StageSpec(spec v0StageSpec) error {
	if spec.Name == "" {
		return &missingRequiredField{Field: "name"}
	}
	if spec.Host == "" {
		return &missingRequiredField{Field: "host"}
	}
	if spec.Port == 0 {
		return &missingRequiredField{Field: "port"}
	}
	return nil
}

func valV0LinkSpec(spec v0LinkSpec) error {
	if err := valV0LinkEndpointSpec(spec.Source); err != nil {
		return err
	}
	if err := valV0LinkEndpointSpec(spec.Target); err != nil {
		return err
	}
	return nil
}

func valV0LinkEndpointSpec(spec v0LinkEndpoint) error {
	if spec.Stage == "" {
		return &missingRequiredField{Field: "stage"}
	}
	return nil
}

type v0FileSpec struct {
	Stages []v0StageSpec `yaml:"stages"`
	Links  []v0LinkSpec  `yaml:"links"`
}

type v0StageSpec struct {
	// Name that should be associated with the stage.
	// (required, unique)
	Name string `yaml:"name"`
	// Host where the server for the stage is running.
	// (required)
	Host string `yaml:"host"`
	// Port where the server for the stage is running.
	// (required)
	Port int `yaml:"port"`
	// Service specifies the name of the service to call.
	// (optional)
	Service string `yaml:"service"`
	// Method specifies the name of the method to call.
	// (optional)
	Method string `yaml:"method"`
}

type v0LinkSpec struct {
	Source v0LinkEndpoint `yaml:"source"`
	Target v0LinkEndpoint `yaml:"target"`
}

type v0LinkEndpoint struct {
	// Stage is the name of the stage to connect to.
	// (required)
	Stage string `yaml:"stage"`
	// Field of the stage message to use. Don't specify for entire message.
	// (optional)
	Field string `yaml:"field"`
}
