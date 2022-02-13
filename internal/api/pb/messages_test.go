package pb

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"gotest.tools/v3/assert"
	"testing"
)

func TestCreateOrchestrationRequest(t *testing.T) {
	var unmarshalled api.CreateOrchestrationRequest

	orig := &pb.CreateOrchestrationRequest{
		Name: "Orchestration",
	}

	UnmarshalCreateOrchestrationRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
}

func TestGetOrchestrationRequest(t *testing.T) {
	var unmarshalled api.GetOrchestrationRequest

	orig := &pb.GetOrchestrationRequest{
		Name:  "Orchestration",
		Phase: string(api.OrchestrationRunning),
	}

	UnmarshalGetOrchestrationRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.Phase, string(unmarshalled.Phase))
}

func TestCreateStageRequest(t *testing.T) {
	var unmarshalled api.CreateStageRequest

	orig := &pb.CreateStageRequest{
		Name:          "Stage",
		Service:       "Service",
		Rpc:           "Rpc",
		Address:       "Address",
		Host:          "Host",
		Port:          1234,
		Orchestration: "Orchestration",
		Asset:         "Asset",
	}

	UnmarshalCreateStageRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.Service, unmarshalled.Service)
	assert.Equal(t, orig.Rpc, unmarshalled.Rpc)
	assert.Equal(t, orig.Address, unmarshalled.Address)
	assert.Equal(t, orig.Host, unmarshalled.Host)
	assert.Equal(t, orig.Port, unmarshalled.Port)
	assert.Equal(t, orig.Orchestration, string(unmarshalled.Orchestration))
	assert.Equal(t, orig.Asset, string(unmarshalled.Asset))
}

func TestGetStageRequest(t *testing.T) {
	var unmarshalled api.GetStageRequest

	orig := &pb.GetStageRequest{
		Name:          "Stage",
		Phase:         string(api.StageFailed),
		Service:       "Service",
		Rpc:           "Rpc",
		Address:       "Address",
		Orchestration: "Orchestration",
		Asset:         "Asset",
	}

	UnmarshalGetStageRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.Phase, string(unmarshalled.Phase))
	assert.Equal(t, orig.Service, unmarshalled.Service)
	assert.Equal(t, orig.Rpc, unmarshalled.Rpc)
	assert.Equal(t, orig.Address, unmarshalled.Address)
	assert.Equal(t, orig.Orchestration, string(unmarshalled.Orchestration))
	assert.Equal(t, orig.Asset, string(unmarshalled.Asset))
}

func TestCreateLinkRequest(t *testing.T) {
	var unmarshalled api.CreateLinkRequest

	orig := &pb.CreateLinkRequest{
		Name:          "Stage",
		SourceStage:   "SourceStage",
		SourceField:   "SourceField",
		TargetStage:   "TargetStage",
		TargetField:   "TargetField",
		Orchestration: "Orchestration",
	}

	UnmarshalCreateLinkRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.SourceStage, string(unmarshalled.SourceStage))
	assert.Equal(t, orig.SourceField, unmarshalled.SourceField)
	assert.Equal(t, orig.TargetStage, string(unmarshalled.TargetStage))
	assert.Equal(t, orig.TargetField, unmarshalled.TargetField)
	assert.Equal(t, orig.Orchestration, string(unmarshalled.Orchestration))
}

func TestGetLinkRequest(t *testing.T) {
	var unmarshalled api.GetLinkRequest

	orig := &pb.GetLinkRequest{
		Name:          "Stage",
		SourceStage:   "SourceStage",
		SourceField:   "SourceField",
		TargetStage:   "TargetStage",
		TargetField:   "TargetField",
		Orchestration: "Orchestration",
	}

	UnmarshalGetLinkRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.SourceStage, string(unmarshalled.SourceStage))
	assert.Equal(t, orig.SourceField, unmarshalled.SourceField)
	assert.Equal(t, orig.TargetStage, string(unmarshalled.TargetStage))
	assert.Equal(t, orig.TargetField, unmarshalled.TargetField)
	assert.Equal(t, orig.Orchestration, string(unmarshalled.Orchestration))
}

func TestCreateAssetRequest(t *testing.T) {
	var unmarshalled api.CreateAssetRequest

	orig := &pb.CreateAssetRequest{
		Name:  "Asset",
		Image: "Image",
	}

	UnmarshalCreateAssetRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.Image, unmarshalled.Image)
}

func TestGetAssetRequest(t *testing.T) {
	var unmarshalled api.GetAssetRequest

	orig := &pb.GetAssetRequest{
		Name:  "Asset",
		Image: "Image",
	}

	UnmarshalGetAssetRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Name, string(unmarshalled.Name))
	assert.Equal(t, orig.Image, unmarshalled.Image)
}

func TestStartExecutionRequest(t *testing.T) {
	var unmarshalled api.StartExecutionRequest

	orig := &pb.StartExecutionRequest{Orchestration: "Orchestration"}

	UnmarshalStartExecutionRequest(&unmarshalled, orig)
	assert.Equal(t, orig.Orchestration, string(unmarshalled.Orchestration))
}
