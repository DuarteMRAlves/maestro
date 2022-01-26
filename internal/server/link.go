package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(req *api.CreateLinkRequest) error {
	logs.LogCreateLinkRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateLink(txn, req)
		},
	)
}

func (s *Server) GetLink(req *api.GetLinkRequest) ([]*api.Link, error) {
	var (
		links []*api.Link
		err   error
	)
	logs.LogGetLinkRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			links, err = s.storageManager.GetMatchingLinks(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}
