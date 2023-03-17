package handlers

import (
	"net/http"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/http-response"
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

// HandleGetTotalReputation
//
//	@Summary	Calculate address TotalReputation
//
//	@Router		/accounts/{address}/total-reputation [GET]
//
//	@Param		address		path		string	true	"Hash or PublicKey"	maxlength(66)
//
//	@Success	200			{object}	http_response.SuccessResponse{data=entities.TotalReputation}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Reputation
func (h *Reputation) HandleGetTotalReputation(w http.ResponseWriter, r *http.Request) {
	//_, err := http_params.ParseOptionalHash("address", r)
	//if err != nil {
	//	accountPubKey, err := http_params.ParseOptionalPublicKey("address", r)
	//	if err != nil {
	//		http_response.Error(w, r, errors.NewInvalidInputError("Account address is not a valid account hash or public key"))
	//		return
	//	}
	//	addressHash = accountPubKey.AccountHash()
	//}

	//getTotalReputation := reputation.NewGetTotalReputation()
	//getTotalReputation.SetAddressHash(*addressHash)
	//getTotalReputation.SetEntityManager(h.entityManager)
	//getTotalReputation.SetDAOContractsMetadata(h.daoContractPackageHashes)
	//
	//http_response.FromFunction(getTotalReputation.Execute, w, r)
	http_response.WriteJSON(w, http.StatusOK, nil)
}

// HandleGetAggregatedReputationChange
//
//	@Summary	user AggregatedReputationChange
//
//	@Param		page			query	int			false	"Page number"										default(1)
//	@Param		page_size		query	string		false	"Number of items per page"							default(10)
//	@Param		order_direction	query	string		false	"Sorting direction"									Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query	[]string	false	"Comma-separated list of sorting fields (address)"	collectionFormat(csv)	default(date)
//	@Router		/accounts/{address}/aggregated-reputation-changes [GET]
//
//	@Param		address		path		string	true	"Hash or PublicKey"	maxlength(66)
//
//	@Success	200			{object}	http_response.PaginatedResponse{data=[]entities.AggregatedReputationChange}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Reputation
func (h *Reputation) HandleGetAggregatedReputationChange(w http.ResponseWriter, r *http.Request) {
	//addressHash, err := http_params.ParseOptionalHash("address", r)
	//if err != nil {
	//	accountPubKey, err := http_params.ParseOptionalPublicKey("address", r)
	//	if err != nil {
	//		http_response.Error(w, r, errors.NewInvalidInputError("Account address is not a valid account hash or public key"))
	//		return
	//	}
	//	addressHash = accountPubKey.AccountHash()
	//}
	//
	//paginationParams := pagination.NewParamsFromRequest(r)
	//
	//getAggregatedReputation := reputation.NewGetAggregatedReputationChanges()
	//getAggregatedReputation.SetAddressHash(*addressHash)
	//getAggregatedReputation.SetEntityManager(h.entityManager)
	//getAggregatedReputation.SetPaginationParams(paginationParams)
	//getAggregatedReputation.SetDAOContractsMetadata(h.daoContractPackageHashes)
	//
	//http_response.FromFunction(getAggregatedReputation.Execute, w, r)

	http_response.WriteJSON(w, http.StatusOK, nil)
}
