package server

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"google.golang.org/grpc"
	"log"
	"net"
)

const grpcNotConfigured = "grpc server not configured"

// Server is the main class that handles the requests
// It implements the InternalAPI interface and manages all requests
type Server struct {
	assetStore     asset.Store
	stageStore     stage.Store
	linkStore      link.Store
	blueprintStore blueprint.Store
	grpcServer     *grpc.Server
}

func (s *Server) ServeGrpc(lis net.Listener) error {
	if ok, err := assert.Status(s.grpcServer != nil, grpcNotConfigured); !ok {
		return err
	}
	return s.grpcServer.Serve(lis)
}

func (s *Server) CreateAsset(config *asset.Asset) error {
	log.Printf("Create Asset with config='%v'\n", config)
	return s.assetStore.Create(config)
}

func (s *Server) GetAsset(query *asset.Asset) []*asset.Asset {
	log.Printf("Get Asset with query='%v'\n", query)
	return s.assetStore.Get(query)
}

func (s *Server) CreateStage(config *stage.Stage) error {
	log.Printf("Create Stage with config='%v'\n", config)
	if config.Asset == "" {
		return errdefs.InvalidArgumentWithMsg("empty asset name")
	}
	if !s.assetStore.Contains(config.Asset) {
		return errdefs.NotFoundWithMsg("asset '%v' not found", config.Asset)
	}
	return s.stageStore.Create(config)
}

func (s *Server) GetStage(query *stage.Stage) []*stage.Stage {
	log.Printf("Get Stage with query=%v", query)
	return s.stageStore.Get(query)
}

func (s *Server) CreateLink(config *link.Link) error {
	log.Printf("Create Stage with config='%v'\n", config)
	if config.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if config.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if !s.stageStore.Contains(config.SourceStage) {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			config.SourceStage)
	}
	if !s.stageStore.Contains(config.TargetStage) {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			config.TargetStage)
	}
	return s.linkStore.Create(config)
}

func (s *Server) GetLink(query *link.Link) []*link.Link {
	log.Printf("Get Link with query=%v", query)
	return s.linkStore.Get(query)
}

func (s *Server) CreateBlueprint(config *blueprint.Blueprint) error {
	log.Printf("Create Blueprint with config=%v", config)
	return s.blueprintStore.Create(config)
}
