package logs

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"go.uber.org/zap"
)

func LogCreateOrchestrationRequest(
	logger *zap.Logger,
	req *api.CreateOrchestrationRequest,
) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		links := make([]string, 0, len(req.Links))
		for _, l := range req.Links {
			links = append(links, string(l))
		}
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.Strings("links", links),
		}
	}
	logger.Info("Create Orchestration.", logFields...)
}

func LogGetOrchestrationRequest(
	logger *zap.Logger,
	req *api.GetOrchestrationRequest,
) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{zap.String("name", string(req.Name))}
	}
	logger.Info("Get Orchestration.", logFields...)
}

func LogCreateStageRequest(logger *zap.Logger, req *api.CreateStageRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("asset", string(req.Asset)),
			zap.String("service", req.Service),
			zap.String("rpc", req.Rpc),
			zap.String("address", req.Address),
			zap.String("host", req.Host),
			zap.Int32("port", req.Port),
		}
	}
	logger.Info("Create Stage.", logFields...)
}

func LogGetStageRequest(logger *zap.Logger, req *api.GetStageRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("phase", string(req.Phase)),
			zap.String("asset", string(req.Asset)),
			zap.String("service", req.Service),
			zap.String("rpc", req.Rpc),
			zap.String("address", req.Address),
		}
	}
	logger.Info("Get Stage.", logFields...)
}

func LogCreateLinkRequest(logger *zap.Logger, req *api.CreateLinkRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("source-stage", string(req.SourceStage)),
			zap.String("source-field", req.SourceField),
			zap.String("target-stage", string(req.TargetStage)),
			zap.String("target-field", req.TargetField),
		}
	}
	logger.Info("Create Link.", logFields...)
}

func LogGetLinkRequest(logger *zap.Logger, req *api.GetLinkRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("source-stage", string(req.SourceStage)),
			zap.String("source-field", req.SourceField),
			zap.String("target-stage", string(req.TargetStage)),
			zap.String("target-field", req.TargetField),
		}
	}
	logger.Info("Get Link.", logFields...)
}

func LogCreateAssetRequest(logger *zap.Logger, req *api.CreateAssetRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("request", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("image", req.Image),
		}
	}
	logger.Info("Create Asset.", logFields...)
}

func LogGetAssetRequest(logger *zap.Logger, req *api.GetAssetRequest) {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("request", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("image", req.Image),
		}
	}
	logger.Info("Get Asset.", logFields...)
}
