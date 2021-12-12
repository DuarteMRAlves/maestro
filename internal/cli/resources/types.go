package resources

type AssetResource struct {
	Name  string `yaml:"name,required"`
	Image string `yaml:"image"`
}

type StageResource struct {
	Name    string `yaml:"name,required"`
	Asset   string `yaml:"asset"`
	Service string `yaml:"service"`
	Method  string `yaml:"method"`
}

type LinkResource struct {
	Name        string `yaml:"name,required"`
	SourceStage string `yaml:"source_stage,required"`
	SourceField string `yaml:"source_field"`
	TargetStage string `yaml:"target_stage,required"`
	TargetField string `yaml:"target_field"`
}
