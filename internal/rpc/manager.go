package rpc

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc"
	"sync"
)

// Manager stores the RPC objects for a given stage.
type Manager interface {
	// FindRpc searches for a given RPC using rpc. It then caches the rpc
	// to be later used. The RPC is associated with the received stage name.
	FindRpc(context.Context, api.StageName, *FindQuery) error
	// GetRpc retrieves an already loaded RPC associated with the given stage.
	// Returns the (RPC, true) if it exists and (nil, false) otherwise.
	GetRpc(api.StageName) (RPC, bool)
}

// FindQuery specifies the search config to be applied
type FindQuery struct {
	// Conn specifies the connection to be used when performing the search.
	Conn grpc.ClientConnInterface
	// Service specifies the name of the service to be searched. If not defined,
	// only one service should exist, which is the one used.
	Service string
	// Rpc specifies the name of the method to be searched. If not defined,
	// only one method should exist inside the grpc service.
	Rpc string
}

type manager struct {
	rpcs sync.Map
}

func NewManager() Manager {
	return &manager{rpcs: sync.Map{}}
}

func (m *manager) FindRpc(
	ctx context.Context,
	stage api.StageName,
	cfg *FindQuery,
) error {
	reflectionClient := NewReflectionClient(ctx, cfg.Conn)
	availableServices, err := reflectionClient.ListServices()
	if err != nil {
		return errdefs.PrependMsg(err, "find rpc")
	}
	serviceName, err := findService(availableServices, cfg)
	if err != nil {
		return errdefs.PrependMsg(err, "find rpc")
	}
	service, err := reflectionClient.ResolveService(serviceName)
	if err != nil {
		return errdefs.PrependMsg(err, "find rpc")
	}
	rpc, err := findRpc(service.RPCs(), cfg)
	if err != nil {
		return errdefs.PrependMsg(err, "find rpc")
	}
	m.rpcs.Store(stage, rpc)
	return nil
}

// findService finds the service that should be used to call the rpc.
// It tries to find the specified service among the available services. If the
// service is not specified, then only one available service must exist that
// will be used. An error is returned if none of the above conditions is
// verified.
func findService(available []string, cfg *FindQuery) (string, error) {
	search := cfg.Service
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
func findRpc(
	available []RPC,
	cfg *FindQuery,
) (RPC, error) {
	search := cfg.Rpc
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

func (m *manager) GetRpc(stage api.StageName) (RPC, bool) {
	rpc, ok := m.rpcs.Load(stage)
	if !ok {
		return nil, false
	}
	return rpc.(RPC), true
}
