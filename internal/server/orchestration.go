package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"go.uber.org/zap"
)

func (s *Server) CreateOrchestration(config *apitypes.Orchestration) error {
	s.logger.Info(
		"Create Orchestration.",
		logOrchestration(config, "config")...)
	return s.orchestrationManager.CreateOrchestration(config)
}

func (s *Server) GetOrchestration(
	query *apitypes.Orchestration,
) []*apitypes.Orchestration {
	s.logger.Info("Get Orchestration.", logOrchestration(query, "query")...)
	return s.orchestrationManager.GetMatchingOrchestration(query)
}

func logOrchestration(
	o *apitypes.Orchestration,
	field string,
) []zap.Field {
	if o == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, string(l))
	}
	return []zap.Field{
		zap.String("name", string(o.Name)),
		zap.Strings("links", links),
	}
}
