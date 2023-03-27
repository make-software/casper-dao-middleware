package handlers

import (
	"net/http"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/reputation"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/errors"
	http_params "casper-dao-middleware/pkg/http-params"
	"casper-dao-middleware/pkg/http-response"
	"casper-dao-middleware/pkg/pagination"
)

type Reputation struct {
	entityManager            persistence.EntityManager
	daoContractPackageHashes utils.DAOContractsMetadata
}

func NewReputation(entityManager persistence.EntityManager, packageHashes utils.DAOContractsMetadata) *Reputation {
	return &Reputation{
		entityManager:            entityManager,
		daoContractPackageHashes: packageHashes,
	}
}

// HandleGetTotalReputationSnapshots
//
//	@Summary	Return paginated list of total-reputation-snapshots for account
//
//	@Router		/accounts/{address}/total-reputation-snapshots [GET]
//
//	@Param		address		path		string	true	"Hash or PublicKey"	maxlength(66)
//
//	@Success	200			{object}	http_response.PaginatedResponse{data=entities.TotalReputationSnapshot}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Reputation
func (h *Reputation) HandleGetTotalReputationSnapshots(w http.ResponseWriter, r *http.Request) {
	addressHash, err := http_params.ParseOptionalHash("address", r)
	if err != nil {
		accountPubKey, err := http_params.ParseOptionalPublicKey("address", r)
		if err != nil {
			http_response.Error(w, r, errors.NewInvalidInputError("Account address is not a valid account hash or public key"))
			return
		}
		addressHash = accountPubKey.AccountHash()
	}

	paginationParams := pagination.NewParamsFromRequest(r)
	paginationParams.SetDefaultOrder("timestamp", pagination.OrderDirectionASC)

	getTotalReputation := reputation.NewGetTotalReputationSnapshots()
	getTotalReputation.SetAddress(addressHash)
	getTotalReputation.SetEntityManager(h.entityManager)
	getTotalReputation.SetPaginationParams(paginationParams)

	http_response.FromFunction(getTotalReputation.Execute, w, r)
}
