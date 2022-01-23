package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// validateCreateOrchestrationConfig verifies if all the conditions to create an
// orchestration are met. It returns an error if one condition is not met and
// nil otherwise.
func validateCreateOrchestrationConfig(cfg *apitypes.Orchestration) error {
	if ok, err := validate.ArgNotNil(cfg, "cfg"); !ok {
		return err
	}
	if !naming.IsValidOrchestrationName(cfg.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", cfg.Name)
	}
	return nil
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateStageConfig(cfg *apitypes.Stage) error {
	if ok, err := validate.ArgNotNil(cfg, "cfg"); !ok {
		return err
	}
	if !naming.IsValidStageName(cfg.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			cfg.Name)
	}
	if cfg.Phase != "" {
		return errdefs.InvalidArgumentWithMsg("phase should not be specified")
	}
	// Asset is not required but if specified should exist.
	if cfg.Asset != "" && !m.assets.Contains(cfg.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			cfg.Asset)
	}
	if cfg.Address != "" && cfg.Host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and host for stage")
	}
	if cfg.Address != "" && cfg.Port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and port for stage")
	}
	return nil
}

// validateCreateLinkConfig verifies if all conditions to create a link are met.
// It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateLinkConfig(config *apitypes.Link) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidLinkName(config.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", config.Name)
	}
	if m.ContainsLink(config.Name) {
		return errdefs.AlreadyExistsWithMsg("link '%v' already exists", config.Name)
	}
	if config.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if config.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if config.SourceStage == config.TargetStage {
		return errdefs.InvalidArgumentWithMsg(
			"source and target stages are equal")
	}
	source, ok := m.GetStageByName(config.SourceStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			config.SourceStage)
	}
	target, ok := m.GetStageByName(config.TargetStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			config.TargetStage)
	}

	if !source.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"source stage is not in Pending phase for link %s",
			config.Name)
	}
	if !target.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"target stage is not in Pending phase for link %s",
			config.Name)
	}
	return nil
}
