package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/dgraph-io/badger/v3"
)

// validateCreateOrchestrationConfig verifies if all the conditions to create an
// orchestration are met. It returns an error if one condition is not met and
// nil otherwise.
func validateCreateOrchestrationConfig(
	txn *badger.Txn,
	req *api.CreateOrchestrationRequest,
) error {
	if ok, err := validate.ArgNotNil(req, "req"); !ok {
		return err
	}
	if !naming.IsValidOrchestrationName(req.Name) {
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
	if ok, err := validate.ArgNotNil(req, "req"); !ok {
		return err
	}
	if !naming.IsValidStageName(req.Name) {
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
	cfg *api.Link,
) error {
	if ok, err := validate.ArgNotNil(cfg, "cfg"); !ok {
		return err
	}
	if !naming.IsValidLinkName(cfg.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", cfg.Name)
	}
	prev, _ := txn.Get(linkKey(cfg.Name))
	if prev != nil {
		return errdefs.AlreadyExistsWithMsg(
			"link '%v' already exists",
			cfg.Name,
		)
	}
	if cfg.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if cfg.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if cfg.SourceStage == cfg.TargetStage {
		return errdefs.InvalidArgumentWithMsg(
			"source and target stages are equal",
		)
	}
	source, ok := m.GetStageByName(txn, cfg.SourceStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			cfg.SourceStage,
		)
	}
	target, ok := m.GetStageByName(txn, cfg.TargetStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			cfg.TargetStage,
		)
	}

	if !source.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"source stage is not in Pending phase for link %s",
			cfg.Name,
		)
	}
	if !target.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"target stage is not in Pending phase for link %s",
			cfg.Name,
		)
	}
	return nil
}

// validateCreateAssetRequest verifies if all conditions to create an asset are
// met. It returns an error if a condition is not met and nil otherwise.
func validateCreateAssetRequest(req *api.CreateAssetRequest) error {
	if ok, err := validate.ArgNotNil(req, "req"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !naming.IsValidAssetName(req.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			req.Name,
		)
	}
	return nil
}
