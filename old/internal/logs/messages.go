package logs

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"go.uber.org/zap"
)

func FieldsForCreateOrchestrationRequest(
	req *api.CreateOrchestrationRequest,
) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	return fields
}

func FieldsForGetOrchestrationRequest(
	req *api.GetOrchestrationRequest,
) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	return fields
}

func FieldsForCreateStageRequest(req *api.CreateStageRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.Asset != "" {
		fields = append(fields, zap.String("asset", string(req.Asset)))
	}
	if req.Service != "" {
		fields = append(fields, zap.String("service", req.Service))
	}
	if req.Rpc != "" {
		fields = append(fields, zap.String("rpc", req.Rpc))
	}
	if req.Address != "" {
		fields = append(fields, zap.String("address", req.Address))
	}
	if req.Host != "" {
		fields = append(fields, zap.String("host", req.Host))
	}
	if req.Port != 0 {
		fields = append(fields, zap.Int32("port", req.Port))
	}
	return fields
}

func FieldsForGetStageRequest(req *api.GetStageRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.Phase != "" {
		fields = append(fields, zap.String("phase", string(req.Phase)))
	}
	if req.Asset != "" {
		fields = append(fields, zap.String("asset", string(req.Asset)))
	}
	if req.Service != "" {
		fields = append(fields, zap.String("service", req.Service))
	}
	if req.Rpc != "" {
		fields = append(fields, zap.String("rpc", req.Rpc))
	}
	if req.Orchestration != "" {
		fields = append(
			fields,
			zap.String("orchestration", string(req.Orchestration)),
		)
	}
	if req.Address != "" {
		fields = append(fields, zap.String("address", req.Address))
	}
	return fields
}

func FieldsForCreateLinkRequest(req *api.CreateLinkRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.SourceStage != "" {
		fields = append(
			fields,
			zap.String("source-stage", string(req.SourceStage)),
		)
	}
	if req.SourceField != "" {
		fields = append(fields, zap.String("source-field", req.SourceField))
	}
	if req.TargetStage != "" {
		fields = append(
			fields,
			zap.String("target-stage", string(req.TargetStage)),
		)
	}
	if req.TargetField != "" {
		fields = append(fields, zap.String("target-field", req.TargetField))
	}
	return fields
}

func FieldsForGetLinkRequest(req *api.GetLinkRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.SourceStage != "" {
		fields = append(
			fields,
			zap.String("source-stage", string(req.SourceStage)),
		)
	}
	if req.SourceField != "" {
		fields = append(fields, zap.String("source-field", req.SourceField))
	}
	if req.TargetStage != "" {
		fields = append(
			fields,
			zap.String("target-stage", string(req.TargetStage)),
		)
	}
	if req.TargetField != "" {
		fields = append(fields, zap.String("target-field", req.TargetField))
	}
	return fields
}

func FieldsForCreateAssetRequest(req *api.CreateAssetRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.Image != "" {
		fields = append(fields, zap.String("image", req.Image))
	}
	return fields
}

func FieldsForGetAssetRequest(req *api.GetAssetRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Name != "" {
		fields = append(fields, zap.String("name", string(req.Name)))
	}
	if req.Image != "" {
		fields = append(fields, zap.String("image", req.Image))
	}
	return fields
}

func FieldsForStartExecutionRequest(req *api.StartExecutionRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Orchestration != "" {
		fields = append(
			fields,
			zap.String("orchestration", string(req.Orchestration)),
		)
	}
	return fields
}

func FieldsForAttachExecutionRequest(req *api.AttachExecutionRequest) []zap.Field {
	if req == nil {
		return []zap.Field{zap.String("req", "nil")}
	}
	fields := make([]zap.Field, 0)
	if req.Orchestration != "" {
		fields = append(
			fields,
			zap.String("orchestration", string(req.Orchestration)),
		)
	}
	return fields
}
