package resources

type LinkResource struct {
	Name        string `yaml:"name,required"`
	SourceStage string `yaml:"source_stage,required"`
	SourceField string `yaml:"source_field"`
	TargetStage string `yaml:"target_stage,required"`
	TargetField string `yaml:"target_field"`
}
