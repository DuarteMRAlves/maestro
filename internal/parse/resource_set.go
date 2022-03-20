package parse

import "github.com/DuarteMRAlves/maestro/internal"

type ResourceSet struct {
	Orchestrations []Orchestration
	Stages         []Stage
	Links          []Link
	Assets         []Asset
}

type Orchestration struct {
	Name internal.OrchestrationName
}

type Stage struct {
	Name          internal.StageName
	Method        MethodContext
	Orchestration internal.OrchestrationName
}

type MethodContext struct {
	Address internal.Address
	Service internal.Service
	Method  internal.Method
}

type Link struct {
	Name          internal.LinkName
	Source        LinkEndpoint
	Target        LinkEndpoint
	Orchestration internal.OrchestrationName
}

type LinkEndpoint struct {
	Stage internal.StageName
	Field internal.MessageField
}

type Asset struct {
	Name  internal.AssetName
	Image internal.Image
}
