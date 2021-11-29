package server

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"log"
)

func (s *Server) CreateAsset(config *asset.Asset) error {
	log.Printf("Create Asset with config='%v'\n", config)
	return s.assetStore.Create(config)
}

func (s *Server) GetAsset(query *asset.Asset) []*asset.Asset {
	log.Printf("Get Asset with query='%v'\n", query)
	return s.assetStore.Get(query)
}
