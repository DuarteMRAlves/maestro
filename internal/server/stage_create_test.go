package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

const stageName = "stage-name"

func TestServer_CreateStage(t *testing.T) {
	var (
		lis                         net.Listener
		registerTest, registerExtra bool
	)

	lis = testutil.ListenAvailablePort(t)
	testAddr := lis.Addr().String()
	registerTest, registerExtra = true, false
	testServer := testutil.StartTestServer(t, lis, registerTest, registerExtra)
	defer testServer.GracefulStop()

	lis = testutil.ListenAvailablePort(t)
	extraAddr := lis.Addr().String()
	registerTest, registerExtra = false, true
	extraServer := testutil.StartTestServer(t, lis, registerTest, registerExtra)
	defer extraServer.GracefulStop()

	lis = testutil.ListenAvailablePort(t)
	bothAddr := lis.Addr().String()
	registerTest, registerExtra = true, true
	bothServer := testutil.StartTestServer(t, lis, registerTest, registerExtra)
	defer bothServer.GracefulStop()

	tests := []struct {
		name   string
		config *stage.Stage
	}{
		{
			name: "correct with nil asset, service and method",
			config: &stage.Stage{
				Name: stageName,
				// ExtraServer only has one server and method
				Address: extraAddr,
			},
		},
		{
			name: "correct with no service and specified method",
			config: &stage.Stage{
				Name:    stageName,
				Asset:   assetNameForNum(0),
				Service: "",
				Method:  "ClientStream",
				// testServer only has one service but four methods
				Address: testAddr,
			},
		},
		{
			name: "correct with service and no method",
			config: &stage.Stage{
				Name:    stageName,
				Asset:   assetNameForNum(0),
				Service: "pb.ExtraService",
				Method:  "",
				// both has two services and one method for ExtraService
				Address: bothAddr,
			},
		},
		{
			name: "correct with service and no method",
			config: &stage.Stage{
				Name:    stageName,
				Asset:   assetNameForNum(0),
				Service: "pb.TestService",
				Method:  "BidiStream",
				// both has two services and four methods for TestService
				Address: bothAddr,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				populateForStages(t, s)
				err = s.CreateStage(test.config)
				assert.NilError(t, err, "create stage error")
			})
	}
}

func TestServer_CreateStage_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForOrchestrations(t, s)

	err = s.CreateStage(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *stage.Stage
	}{
		{
			name: "empty name",
			config: &stage.Stage{
				Name:  "",
				Asset: assetNameForNum(0),
			},
		},
		{
			name: "invalid characters in name",
			config: &stage.Stage{
				Name:  "some@name",
				Asset: assetNameForNum(0),
			},
		},
		{
			name: "invalid character sequence",
			config: &stage.Stage{
				Name:  "other-/name",
				Asset: assetNameForNum(0),
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

				err = s.CreateStage(test.config)
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

func TestServer_CreateStage_AssetNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &stage.Stage{
		Name:  stageName,
		Asset: assetNameForNum(1),
	}

	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("asset '%v' not found", assetNameForNum(1))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_AlreadyExists(t *testing.T) {
	var err error

	lis := testutil.ListenAvailablePort(t)
	bothAddr := lis.Addr().String()
	bothServer := testutil.StartTestServer(t, lis, true, true)
	defer bothServer.GracefulStop()

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &stage.Stage{
		Name:    stageName,
		Asset:   assetNameForNum(0),
		Service: "pb.TestService",
		Method:  "BidiStream",
		// both has two services and four methods for TestService
		Address: bothAddr,
	}

	err = s.CreateStage(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not AlreadyExists")
	expectedMsg := fmt.Sprintf("stage '%v' already exists", stageName)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_Error(t *testing.T) {
	tests := []struct {
		name            string
		registerTest    bool
		registerExtra   bool
		config          *stage.Stage
		verifyErrTypeFn func(err error) bool
		expectedErrMsg  string
	}{
		{
			name:          "no services",
			registerTest:  false,
			registerExtra: false,
			config: &stage.Stage{
				Name:  stageName,
				Asset: assetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"find service without name for stage %v: expected 1 "+
					"available service but 0 found",
				stageName),
		},
		{
			name:          "too many services",
			registerTest:  true,
			registerExtra: true,
			config: &stage.Stage{
				Name:  stageName,
				Asset: assetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"find service without name for stage %v: expected 1 "+
					"available service but 2 found",
				stageName),
		},
		{
			name:          "no such service",
			registerTest:  true,
			registerExtra: true,
			config: &stage.Stage{
				Name:    stageName,
				Asset:   assetNameForNum(0),
				Service: "NoSuchService",
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg: fmt.Sprintf(
				"service with name NoSuchService not found for stage %v",
				stageName),
		},
		{
			name:         "too many methods",
			registerTest: true,
			config: &stage.Stage{
				Name:  stageName,
				Asset: assetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"find rpc without name for stage %v: expected 1 available "+
					"rpc but 4 found",
				stageName),
		},
		{
			name:          "no such method",
			registerTest:  true,
			registerExtra: false,
			config: &stage.Stage{
				Name:   stageName,
				Asset:  assetNameForNum(0),
				Method: "NoSuchMethod",
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg: fmt.Sprintf(
				"rpc with name NoSuchMethod not found for stage %v",
				stageName),
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var err error

				lis := testutil.ListenAvailablePort(t)
				bothAddr := lis.Addr().String()
				bothServer := testutil.StartTestServer(
					t,
					lis,
					test.registerTest,
					test.registerExtra)
				defer bothServer.GracefulStop()

				s, err := NewBuilder().
					WithGrpc().
					WithLogger(testutil.NewLogger(t)).
					Build()
				assert.NilError(t, err, "build server")
				populateForStages(t, s)

				test.config.Address = bothAddr

				err = s.CreateStage(test.config)
				assert.Assert(
					t,
					test.verifyErrTypeFn(err),
					"incorrect err type")
				assert.Error(t, err, test.expectedErrMsg)
			})
	}

}

func populateForStages(t *testing.T, s *Server) {
	assets := []*asset.Asset{
		assetForNum(0),
	}
	populateAssets(t, s, assets)
}
