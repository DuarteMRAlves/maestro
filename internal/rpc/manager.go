package rpc

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc"
)

// Manager stores the RPC objects for a given stage.
type Manager interface {
	// GetRpc searches for a given RPC using the grpc reflection api.
	GetRpc(context.Context, grpc.ClientConnInterface, *api.Stage) (RPC, error)
}

type manager struct{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) GetRpc(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	stage *api.Stage,
) (RPC, error) {
	var err error

	client := NewReflectionClient(ctx, conn)
	availableServices, err := client.ListServices()
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	srvName, err := findService(availableServices, stage)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	srv, err := client.ResolveService(srvName)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	stageRpc, err := findRpc(srv.RPCs(), stage)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	return stageRpc, nil
}

// findService finds the service that should be used to call the rpc.
// It tries to find the specified service among the available services. If the
// service is not specified, then only one available service must exist that
// will be used. An error is returned if none of the above conditions is
// verified.
func findService(available []string, stage *api.Stage) (string, error) {
	search := stage.Service
	if search == "" {
		if len(available) == 1 {
			return available[0], nil
		}
		return "", errdefs.InvalidArgumentWithMsg(
			"find service without name: expected 1 available service but found %v",
			len(available),
		)
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return "", errdefs.NotFoundWithMsg(
			"find service with name %v: not found",
			search,
		)
	}
}

// findRpc find the rpc to be called. If a rpc name was specified in the config,
// then it verifies it exists in the available rpcs and returns it. Otherwise,
// it verifies only a single rpc is available and returns it.
func findRpc(available []RPC, stage *api.Stage) (RPC, error) {
	search := stage.Rpc
	if search == "" {
		if len(available) == 1 {
			return available[0], nil
		}
		return nil, errdefs.InvalidArgumentWithMsg(
			"find rpc without name: expected 1 available rpc but found %v",
			len(available),
		)
	} else {
		for _, rpc := range available {
			if search == rpc.Name() {
				return rpc, nil
			}
		}
		return nil, errdefs.NotFoundWithMsg(
			"find rpc with name %v: not found",
			search,
		)
	}
}
