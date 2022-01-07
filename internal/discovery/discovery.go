package discovery

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
)

// Config specifies the search config to be applied
type Config struct {
	// Service specifies the name of the service to be searched. If not defined,
	// only one service should exist, which is the one used.
	Service string
	// Rpc specifies the name of the method to be searched. If not defined,
	// only one method should exist inside the grpc service.
	Rpc string
}

func FindRpc(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	cfg *Config,
) (reflection.RPC, error) {
	reflectionClient := reflection.NewClient(ctx, conn)
	availableServices, err := reflectionClient.ListServices()
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	serviceName, err := findService(availableServices, cfg)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	service, err := reflectionClient.ResolveService(serviceName)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	rpc, err := findRpc(service.RPCs(), cfg)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "find rpc")
	}
	return rpc, nil
}

// findService finds the service that should be used to call the rpc.
// It tries to find the specified service among the available services. If the
// service is not specified, then only one available service must exist that
// will be used. An error is returned if none of the above conditions is
// verified.
func findService(available []string, cfg *Config) (string, error) {
	search := cfg.Service
	if search == "" {
		if len(available) == 1 {
			return available[0], nil
		}
		return "", errdefs.InvalidArgumentWithMsg(
			"find service without name: expected 1 available service but found %v",
			len(available))
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return "", errdefs.NotFoundWithMsg(
			"find service with name %v: not found",
			search)
	}
}

// findRpc find the rpc to be called. If a rpc name was specified in the config,
// then it verifies it exists in the available rpcs and returns it. Otherwise,
// it verifies only a single rpc is available and returns it.
func findRpc(
	available []reflection.RPC,
	cfg *Config,
) (reflection.RPC, error) {
	search := cfg.Rpc
	if search == "" {
		if len(available) == 1 {
			return available[0], nil
		}
		return nil, errdefs.InvalidArgumentWithMsg(
			"find rpc without name: expected 1 available rpc but found %v",
			len(available))
	} else {
		for _, rpc := range available {
			if search == rpc.Name() {
				return rpc, nil
			}
		}
		return nil, errdefs.NotFoundWithMsg(
			"find rpc with name %v: not found",
			search)
	}
}
