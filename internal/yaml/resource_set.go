package yaml

import "github.com/DuarteMRAlves/maestro/internal"

type ResourceSet struct {
	Pipelines []Pipeline
	Stages    []Stage
	Links     []Link
	Assets    []Asset
}

type Pipeline struct {
	Name internal.PipelineName
	Mode internal.ExecutionMode
}

type Stage struct {
	Name     internal.StageName
	Method   MethodContext
	Pipeline internal.PipelineName
}

type MethodContext struct {
	Address internal.Address
	Service internal.Service
	Method  internal.Method
}

type Link struct {
	Name     internal.LinkName
	Source   LinkEndpoint
	Target   LinkEndpoint
	Pipeline internal.PipelineName
}

type LinkEndpoint struct {
	Stage internal.StageName
	Field internal.MessageField
}

type Asset struct {
	Name  internal.AssetName
	Image internal.Image
}
