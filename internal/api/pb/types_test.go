package pb

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalAsset(t *testing.T) {
	const (
		name  api.AssetName = "Asset Name"
		image               = "user/image:version"
	)
	tests := []struct {
		in api.Asset
	}{
		{api.Asset{Name: name, Image: image}},
		{
			api.Asset{
				Name:  name,
				Image: image,
			},
		},
		{api.Asset{Name: name}},
		{api.Asset{Image: image}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res, err := MarshalAsset(&in)
				assert.NilError(t, err, "marshal error")
				assert.Equal(t, string(in.Name), res.Name, "Asset Name")
				assert.Equal(t, in.Image, res.Image, "Asset Image")
			},
		)
	}
}

func TestUnmarshalAsset(t *testing.T) {
	tests := []struct {
		in *pb.Asset
	}{
		{&pb.Asset{Name: "Asset Name"}},
		{&pb.Asset{Name: "Asset Name"}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res, err := UnmarshalAsset(in)
				assert.Equal(t, nil, err, "Error")
				assert.Equal(t, in.Name, string(res.Name), "Asset Name")
			},
		)
	}
}

func TestMarshalAssetNil(t *testing.T) {
	res, err := MarshalAsset(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'a' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalAssetNil(t *testing.T) {
	res, err := UnmarshalAsset(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestMarshalOrchestration(t *testing.T) {
	const (
		name  api.OrchestrationName = "OrchestrationName"
		phase                       = api.OrchestrationRunning
		link1 api.LinkName          = "Link Name 1"
		link2 api.LinkName          = "Link Name 2"
		link3 api.LinkName          = "Link Name 3"
	)
	tests := []struct {
		name string
		o    *api.Orchestration
	}{
		{
			name: "all fields with non default values",
			o: &api.Orchestration{
				Name:  name,
				Phase: phase,
				Links: []api.LinkName{link1, link2, link3},
			},
		},
		{
			name: "all field with default values",
			o: &api.Orchestration{
				Name:  "",
				Phase: "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				o := test.o
				res, err := MarshalOrchestration(o)
				assert.NilError(t, err, "marshal error")
				assert.Equal(t, string(o.Name), res.Name, "orchestration name")
				assert.Equal(
					t,
					string(o.Phase),
					res.Phase,
					"orchestration phase",
				)
				assert.Equal(
					t,
					len(o.Links),
					len(res.Links),
					"orchestration links len",
				)
				for i, l := range o.Links {
					assert.Equal(t, string(l), res.Links[i], "link %d equal", i)
				}
			},
		)
	}
}

func TestUnmarshalOrchestrationCorrect(t *testing.T) {
	const (
		name  = "OrchestrationName"
		phase = string(api.OrchestrationPending)
		link1 = "Link Name 1"
		link2 = "Link Name 2"
		link3 = "Link Name 3"
	)
	tests := []struct {
		name string
		o    *pb.Orchestration
	}{
		{
			name: "all fields with non default values",
			o: &pb.Orchestration{
				Name:  name,
				Phase: phase,
				Links: []string{link1, link2, link3},
			},
		},
		{
			name: "all field with default values",
			o: &pb.Orchestration{
				Name:  "",
				Phase: "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				o := test.o
				res, err := UnmarshalOrchestration(o)
				assert.Equal(t, nil, err, "unmarshal error")
				assert.Equal(t, o.Name, string(res.Name), "orchestration name")
				assert.Equal(
					t,
					o.Phase,
					string(res.Phase),
					"orchestration phase",
				)
				assert.Equal(
					t,
					len(o.Links),
					len(res.Links),
					"orchestration links len",
				)
				for i, l := range o.Links {
					assert.Equal(t, l, string(res.Links[i]), "link %d equal", i)
				}
			},
		)
	}
}

func TestMarshalOrchestrationError(t *testing.T) {
	res, err := MarshalOrchestration(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'o' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalOrchestrationNil(t *testing.T) {
	res, err := UnmarshalOrchestration(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestMarshalStage(t *testing.T) {
	const (
		stageName    api.StageName = "Stage Name"
		stagePhase                 = api.StagePending
		stageAsset                 = "Stage Asset"
		stageService               = "stageService"
		stageRpc                   = "stageRpc"
		stageAddress               = "stageAddress"
	)
	tests := []*api.Stage{
		{
			Name:    stageName,
			Phase:   stagePhase,
			Asset:   stageAsset,
			Service: stageService,
			Rpc:     stageRpc,
			Address: stageAddress,
		},
		{
			Name:    "",
			Phase:   "",
			Asset:   "",
			Service: "",
			Rpc:     "",
			Address: "",
		},
	}

	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName, func(t *testing.T) {
				res, err := MarshalStage(s)
				assert.NilError(t, err, "marshal error")
				assertStage(t, s, res)
			},
		)
	}
}

func TestUnmarshalStageCorrect(t *testing.T) {
	const (
		stageName    = "Stage Name"
		stagePhase   = "Running"
		stageAsset   = "Stage Asset"
		stageService = "stageService"
		stageRpc     = "stageRpc"
		stageAddress = "stageAddress"
	)
	tests := []*pb.Stage{
		{
			Name:    stageName,
			Phase:   stagePhase,
			Asset:   stageAsset,
			Service: stageService,
			Rpc:     stageRpc,
			Address: stageAddress,
		},
		{
			Name:    "",
			Phase:   "",
			Asset:   "",
			Service: "",
			Rpc:     "",
			Address: "",
		},
	}
	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName,
			func(t *testing.T) {
				res, err := UnmarshalStage(s)
				assert.Equal(t, nil, err, "unmarshal error")
				assertPbStage(t, s, res)
			},
		)
	}
}

func TestMarshalStageNil(t *testing.T) {
	res, err := MarshalStage(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'s' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalStageNil(t *testing.T) {
	res, err := UnmarshalStage(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertStage(t *testing.T, expected *api.Stage, actual *pb.Stage) {
	assert.Equal(t, string(expected.Name), actual.Name, "stage assetName")
	assert.Equal(t, string(expected.Phase), actual.Phase, "stage phase")
	assert.Equal(t, string(expected.Asset), actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Rpc, actual.Rpc, "stage rpc")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *api.Stage) {
	assert.Equal(t, expected.Name, string(actual.Name), "stage assetName")
	assert.Equal(t, expected.Phase, string(actual.Phase), "stage phase")
	assert.Equal(t, expected.Asset, string(actual.Asset), "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Rpc, actual.Rpc, "stage rpc")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
}

func TestMarshalLink(t *testing.T) {
	const (
		name        api.LinkName  = "name"
		sourceStage api.StageName = "sourceStage"
		sourceField               = "sourceField"
		targetStage api.StageName = "targetStage"
		targetField               = "targetField"
	)
	tests := []*api.Link{
		{
			Name:        name,
			SourceStage: sourceStage,
			SourceField: sourceField,
			TargetStage: targetStage,
			TargetField: targetField,
		},
		{
			Name:        "",
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	}

	for _, l := range tests {
		testName := fmt.Sprintf("link=%v", l)

		t.Run(
			testName, func(t *testing.T) {
				res, err := MarshalLink(l)
				assert.NilError(t, err, "marshal error")
				assertLink(t, l, res)
			},
		)
	}
}

func TestUnmarshalLinkCorrect(t *testing.T) {
	const (
		name        = "name"
		sourceStage = "sourceStage"
		sourceField = "sourceField"
		targetStage = "targetStage"
		targetField = "targetField"
	)
	tests := []*pb.Link{
		{
			Name:        name,
			SourceStage: sourceStage,
			SourceField: sourceField,
			TargetStage: targetStage,
			TargetField: targetField,
		},
		{
			Name:        "",
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
			TargetField: "",
		},
	}
	for _, l := range tests {
		testName := fmt.Sprintf("link=%v", l)

		t.Run(
			testName,
			func(t *testing.T) {
				res, err := UnmarshalLink(l)
				assert.NilError(t, err, "unmarshal error")
				assertPbLink(t, l, res)
			},
		)
	}
}

func TestMarshalLinkNil(t *testing.T) {
	res, err := MarshalLink(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'l' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalLinkNil(t *testing.T) {
	res, err := UnmarshalLink(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertLink(t *testing.T, expected *api.Link, actual *pb.Link) {
	assert.Equal(t, string(expected.Name), actual.Name, "name")
	assert.Equal(
		t,
		string(expected.SourceStage),
		actual.SourceStage,
		"source id",
	)
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field",
	)
	assert.Equal(
		t,
		string(expected.TargetStage),
		actual.TargetStage,
		"target id",
	)
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field",
	)
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *api.Link) {
	assert.Equal(t, expected.Name, string(actual.Name), "name")
	assert.Equal(
		t,
		expected.SourceStage,
		string(actual.SourceStage),
		"source id",
	)
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field",
	)
	assert.Equal(
		t,
		expected.TargetStage,
		string(actual.TargetStage),
		"target id",
	)
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field",
	)
}
