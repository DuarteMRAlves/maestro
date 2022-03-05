package server

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/create"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/dgraph-io/badger/v3"
)

type CreateAsset func(domain.CreateAssetRequest) domain.CreateAssetResponse

func CreateAssetWithTxn(txn *badger.Txn) CreateAsset {
	return func(req domain.CreateAssetRequest) domain.CreateAssetResponse {
		var res domain.AssetResult

		storageFunc := domain.Bind(asset.StoreAssetWithTxn(txn))

		res = create.RequestToResult(req)
		res = storageFunc(res)
		return create.ResultToResponse(res)
	}
}
