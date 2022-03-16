package execute

type execution struct {
	stages *stageMap
}

func newExecution(stages *stageMap) execution {
	return execution{stages: stages}
}
