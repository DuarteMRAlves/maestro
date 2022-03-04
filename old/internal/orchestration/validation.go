package orchestration

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
)

// validateCreateOrchestrationConfig verifies if all the conditions to create an
// orchestration are met. It returns an error if one condition is not met and
// nil otherwise.
func validateCreateOrchestrationConfig(
	txn *badger.Txn,
	req *api.CreateOrchestrationRequest,
) error {
	if req == nil {
		return errdefs.InvalidArgumentWithMsg("'req' is nil")
	}
	if !IsValidOrchestrationName(req.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", req.Name)
	}
	helper := kv.NewTxnHelper(txn)
	if helper.ContainsOrchestration(req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"orchestration '%s' already exists",
			req.Name,
		)
	}
	return nil
}

// validateCreateAssetRequest verifies if all conditions to create an asset are
// met. It returns an error if a condition is not met and nil otherwise.
func validateCreateAssetRequest(
	txn *badger.Txn,
	req *api.CreateAssetRequest,
) error {
	if req == nil {
		return errdefs.InvalidArgumentWithMsg("'req' is nil")
	}
	if !IsValidAssetName(req.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			req.Name,
		)
	}
	helper := kv.NewTxnHelper(txn)
	if helper.ContainsAsset(req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			req.Name,
		)
	}
	return nil
}
