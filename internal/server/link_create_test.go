package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/internal/testutil/mock"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"testing"
)

const linkName = "link-name"

func TestServer_CreateLink(t *testing.T) {

	tests := []struct {
		name   string
		config *link.Link
	}{
		{
			name: "correct with nil fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "correct with empty fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: "stage-1",
				SourceField: "",
				TargetStage: "stage-2",
				TargetField: "",
			},
		},
		{
			name: "correct with fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: "stage-1",
				SourceField: "field4",
				TargetStage: "stage-2",
				TargetField: "fieldName4",
			},
		},
		{
			name: "incompatible outer but compatible inner",
			config: &link.Link{
				Name:        linkName,
				SourceStage: "stage-1",
				SourceField: "field4",
				TargetStage: "stage-incompatible-outer",
				TargetField: "field4",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

				populateForLinks(t, s)
				err = s.CreateLink(test.config)
				assert.NilError(t, err, "create link error")
			})
	}
}

func TestServer_CreateLink_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	err = s.CreateLink(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *link.Link
	}{
		{
			name: "empty name",
			config: &link.Link{
				Name:        "",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid characters in name",
			config: &link.Link{
				Name:        "some'character",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid character sequence",
			config: &link.Link{
				Name:        "//invalid-name",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				populateForOrchestrations(t, s)

				err = s.CreateLink(test.config)
				assert.Assert(
					t,
					errdefs.IsInvalidArgument(err),
					"error is not InvalidArgument")
				expectedMsg := fmt.Sprintf(
					"invalid name '%v'",
					test.config.Name)
				assert.Error(t, err, expectedMsg)
			})
	}
}

func TestServer_CreateLink_SourceEmpty(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	assert.Error(t, err, "empty source stage name")
}

func TestServer_CreateLink_TargetEmpty(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-2",
		TargetStage: "",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	assert.Error(t, err, "empty target stage name")
}

func TestServer_CreateLink_EqualSourceAndTarget(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		TargetStage: "stage-1",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not NotFound")
	assert.Error(t, err, "source and target stages are equal")
}

func TestServer_CreateLink_SourceNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-3",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := "source stage 'stage-3' not found"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_TargetNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		TargetStage: "stage-3",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := "target stage 'stage-3' not found"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_AlreadyExists(t *testing.T) {
	var err error

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("link '%v' already exists", linkName)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_UnknownSourceField(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		SourceField: "unknown-field",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"field with name unknown-field not found for message "+
			"pb.TestMessage1 for source stage in link %v",
		linkName)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_UnknownTargetField(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		TargetStage: "stage-2",
		TargetField: "unknown-field",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"field with name unknown-field not found for message "+
			"pb.TestMessageDiffNames for target stage in link %v",
		linkName)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_IncompatibleMessages(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "stage-1",
		TargetStage: "stage-incompatible-outer",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not IsInvalidArgument")
	expectedMsg := fmt.Sprintf(
		"incompatible message types between source output pb.TestMessage1 "+
			"and target input pb.TestWrongOuterFieldType in link %v",
		linkName)
	assert.Error(t, err, expectedMsg)
}

func populateForLinks(t *testing.T, s *Server) {
	testMsg1Desc, err := desc.LoadMessageDescriptor("pb.TestMessage1")
	assert.NilError(t, err, "load desc test message 1")

	testMsg2Desc, err := desc.LoadMessageDescriptor("pb.TestMessageDiffNames")
	assert.NilError(t, err, "load desc test message 2")

	testIncompatibleDesc, err := desc.LoadMessageDescriptor(
		"pb.TestWrongOuterFieldType")
	assert.NilError(t, err, "load desc test message 2")

	message1, err := reflection.NewMessage(testMsg1Desc)
	assert.NilError(t, err, "test message 1")

	message2, err := reflection.NewMessage(testMsg2Desc)
	assert.NilError(t, err, "test message 2")

	messageIncompatible, err := reflection.NewMessage(testIncompatibleDesc)
	assert.NilError(t, err, "test message 2")

	stage1 := &stage.Stage{
		Name:    "stage-1",
		Asset:   "asset-1",
		Service: "service-1",
		Method:  "method-1",
		Address: "address-1",
		Rpc: &mock.RPC{
			Name_: "method-1",
			FQN:   "service-1/method-1",
			In:    message1,
			Out:   message1,
		},
	}

	stage2 := &stage.Stage{
		Name:    "stage-2",
		Asset:   "asset-2",
		Service: "service-2",
		Method:  "method-2",
		Address: "address-2",
		Rpc: &mock.RPC{
			Name_: "method-2",
			FQN:   "service-2/method-2",
			In:    message2,
			Out:   message2,
		},
	}

	stageIncompatible := &stage.Stage{
		Name:    "stage-incompatible-outer",
		Asset:   "asset-incompatible",
		Service: "service-incompatible",
		Method:  "method-incompatible",
		Address: "address-incompatible",
		Rpc: &mock.RPC{
			Name_: "method-2",
			FQN:   "service-2/method-incompatible",
			In:    messageIncompatible,
			Out:   messageIncompatible,
		},
	}

	stages := []*stage.Stage{stage1, stage2, stageIncompatible}

	populateStages(t, s, stages)
}
