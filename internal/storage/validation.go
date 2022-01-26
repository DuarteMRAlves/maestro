package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
)

// validateCreateOrchestrationConfig verifies if all the conditions to create an
// orchestration are met. It returns an error if one condition is not met and
// nil otherwise.
func validateCreateOrchestrationConfig(
	txn *badger.Txn,
	req *api.CreateOrchestrationRequest,
) error {
	if ok, err := util.ArgNotNil(req, "req"); !ok {
		return err
	}
	if !IsValidOrchestrationName(req.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", req.Name)
	}
	prev, _ := txn.Get(orchestrationKey(req.Name))
	if prev != nil {
		return errdefs.AlreadyExistsWithMsg(
			"orchestration '%s' already exists",
			req.Name,
		)
	}
	return nil
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateStageConfig(
	txn *badger.Txn,
	req *api.CreateStageRequest,
) error {
	if ok, err := util.ArgNotNil(req, "req"); !ok {
		return err
	}
	if !IsValidStageName(req.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			req.Name,
		)
	}
	prev, _ := txn.Get(stageKey(req.Name))
	if prev != nil {
		return errdefs.AlreadyExistsWithMsg(
			"stage '%v' already exists",
			req.Name,
		)
	}
	// Asset is not required but if specified should exist.
	if req.Asset != "" && !m.ContainsAsset(txn, req.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			req.Asset,
		)
	}
	if req.Address != "" && req.Host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and host for stage",
		)
	}
	if req.Address != "" && req.Port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and port for stage",
		)
	}
	return nil
}

// validateCreateLinkConfig verifies if all conditions to create a link are met.
// It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateLinkConfig(
	txn *badger.Txn,
	req *api.CreateLinkRequest,
) error {
	if ok, err := util.ArgNotNil(req, "req"); !ok {
		return err
	}
	if !IsValidLinkName(req.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", req.Name)
	}
	prev, _ := txn.Get(linkKey(req.Name))
	if prev != nil {
		return errdefs.AlreadyExistsWithMsg(
			"link '%v' already exists",
			req.Name,
		)
	}
	if req.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if req.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if req.SourceStage == req.TargetStage {
		return errdefs.InvalidArgumentWithMsg(
			"source and target stages are equal",
		)
	}
	source, ok := m.GetStageByName(txn, req.SourceStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			req.SourceStage,
		)
	}
	target, ok := m.GetStageByName(txn, req.TargetStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			req.TargetStage,
		)
	}

	if source.Phase != api.StagePending {
		return errdefs.FailedPreconditionWithMsg(
			"source stage is not in Pending phase for link %s",
			req.Name,
		)
	}
	if target.Phase != api.StagePending {
		return errdefs.FailedPreconditionWithMsg(
			"target stage is not in Pending phase for link %s",
			req.Name,
		)
	}
	return nil
}

// validateCreateAssetRequest verifies if all conditions to create an asset are
// met. It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateAssetRequest(
	txn *badger.Txn,
	req *api.CreateAssetRequest,
) error {
	if ok, err := util.ArgNotNil(req, "req"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !IsValidAssetName(req.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			req.Name,
		)
	}
	if m.ContainsAsset(txn, req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			req.Name,
		)
	}
	return nil
}
