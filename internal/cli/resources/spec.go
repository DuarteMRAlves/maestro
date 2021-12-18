package resources

type AssetSpec struct {
	Name  string `yaml:"name" info:"required"`
	Image string `yaml:"image"`
}

type StageSpec struct {
	Name    string `yaml:"name" info:"required"`
	Asset   string `yaml:"asset"`
	Service string `yaml:"service"`
	Method  string `yaml:"method"`
	Address string `yaml:"address"`
}

type LinkSpec struct {
	Name        string `yaml:"name" info:"required"`
	SourceStage string `yaml:"source_stage" info:"required"`
	SourceField string `yaml:"source_field"`
	TargetStage string `yaml:"target_stage" info:"required"`
	TargetField string `yaml:"target_field"`
}

type OrchestrationSpec struct {
	Name  string   `yaml:"name" info:"required"`
	Links []string `yaml:"links" info:"required"`
}
