package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

func TestServer_CreateStage(t *testing.T) {
	var (
		lis                         net.Listener
		registerTest, registerExtra bool
	)
	const name = "stage-name"

	lis = testutil.ListenAvailablePort(t)
	tcpAddr, ok := lis.Addr().(*net.TCPAddr)
	assert.Assert(t, ok, "address type cast")
	testAddr := tcpAddr.String()
	testHost := tcpAddr.IP.String()
	testPort := tcpAddr.Port
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
		config *types.Stage
	}{
		{
			name: "nil asset, service and rpc",
			config: &types.Stage{
				Name: name,
				// ExtraServer only has one server and rpc
				Address: extraAddr,
			},
		},
		{
			name: "no service and specified rpc",
			config: &types.Stage{
				Name:    name,
				Asset:   testutil.AssetNameForNum(0),
				Service: "",
				Rpc:     "Unary",
				// testServer only has one service but four rpcs
				Address: testAddr,
			},
		},
		{
			name: "with service and no rpc",
			config: &types.Stage{
				Name:    name,
				Asset:   testutil.AssetNameForNum(0),
				Service: "pb.ExtraService",
				Rpc:     "",
				// both has two services and one rpc for ExtraService
				Address: bothAddr,
			},
		},
		{
			name: "with service and rpc",
			config: &types.Stage{
				Name:    name,
				Asset:   testutil.AssetNameForNum(0),
				Service: "pb.TestService",
				Rpc:     "Unary",
				// both has two services and four rpcs for TestService
				Address: bothAddr,
			},
		},
		{
			name: "from host and port",
			config: &types.Stage{
				Name:    name,
				Asset:   testutil.AssetNameForNum(0),
				Service: "pb.TestService",
				Rpc:     "Unary",
				// both has two services and four rpcs for TestService
				Host: testHost,
				Port: int32(testPort),
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
	populateForStages(t, s)

	err = s.CreateStage(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'cfg' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *types.Stage
	}{
		{
			name: "empty name",
			config: &types.Stage{
				Name:  "",
				Asset: testutil.AssetNameForNum(0),
			},
		},
		{
			name: "invalid characters in name",
			config: &types.Stage{
				Name:  "some@name",
				Asset: testutil.AssetNameForNum(0),
			},
		},
		{
			name: "invalid character sequence",
			config: &types.Stage{
				Name:  "other-/name",
				Asset: testutil.AssetNameForNum(0),
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
	const name = "stage-name"
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &types.Stage{
		Name:  name,
		Asset: testutil.AssetNameForNum(1),
	}

	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("asset '%v' not found", testutil.AssetNameForNum(1))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_AlreadyExists(t *testing.T) {
	var err error
	const name = "stage-name"

	lis := testutil.ListenAvailablePort(t)
	bothAddr := lis.Addr().String()
	bothServer := testutil.StartTestServer(t, lis, true, true)
	defer bothServer.GracefulStop()

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &types.Stage{
		Name:    name,
		Asset:   testutil.AssetNameForNum(0),
		Service: "pb.TestService",
		Rpc:     "Unary",
		// both has two services and four rpcs for TestService
		Address: bothAddr,
	}

	err = s.CreateStage(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not AlreadyExists")
	expectedMsg := fmt.Sprintf("stage '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_Error(t *testing.T) {
	const name = "stage-name"
	tests := []struct {
		name            string
		registerTest    bool
		registerExtra   bool
		config          *types.Stage
		verifyErrTypeFn func(err error) bool
		expectedErrMsg  string
	}{
		{
			name:          "no services",
			registerTest:  false,
			registerExtra: false,
			config: &types.Stage{
				Name:  name,
				Asset: testutil.AssetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"stage %s: find rpc: find service without name: expected 1 "+
					"available service but found 0",
				name),
		},
		{
			name:          "too many services",
			registerTest:  true,
			registerExtra: true,
			config: &types.Stage{
				Name:  name,
				Asset: testutil.AssetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"stage %s: find rpc: find service without name: expected 1 "+
					"available service but found 2",
				name),
		},
		{
			name:          "no such service",
			registerTest:  true,
			registerExtra: true,
			config: &types.Stage{
				Name:    name,
				Asset:   testutil.AssetNameForNum(0),
				Service: "NoSuchService",
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg: fmt.Sprintf(
				"stage %s: find rpc: find service with name NoSuchService:"+
					" not found",
				name),
		},
		{
			name:         "too many rpcs",
			registerTest: true,
			config: &types.Stage{
				Name:  name,
				Asset: testutil.AssetNameForNum(0),
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg: fmt.Sprintf(
				"stage %v: find rpc: find rpc without name: expected 1 "+
					"available rpc but found 4",
				name),
		},
		{
			name:          "no such rpc",
			registerTest:  true,
			registerExtra: false,
			config: &types.Stage{
				Name:  name,
				Asset: testutil.AssetNameForNum(0),
				Rpc:   "NoSuchRpc",
				// Address injected during the test to point to the server
			},
			verifyErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg: fmt.Sprintf(
				"stage %v: find rpc: find rpc with name NoSuchRpc: not found",
				name),
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
