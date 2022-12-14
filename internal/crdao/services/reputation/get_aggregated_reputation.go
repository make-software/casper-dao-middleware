package reputation

import (
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"
)

type GetAggregatedReputationChanges struct {
	di.PaginationParamsAware
	di.EntityManagerAware
	di.DAOContractPackageHashesAware

	addressHash types.Hash
}

func NewGetAggregatedReputationChanges() *GetAggregatedReputationChanges {
	return &GetAggregatedReputationChanges{}
}

func (c *GetAggregatedReputationChanges) SetAddressHash(hash types.Hash) {
	c.addressHash = hash
}

func (c *GetAggregatedReputationChanges) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{
		"address": c.addressHash,
	}

	paginationParams := c.GetPaginationParams()

	count, err := c.GetEntityManager().ReputationChangeRepository().CountAggregatedReputationChanges(filters)
	if err != nil {
		return nil, err
	}

	reputations, err := c.GetEntityManager().ReputationChangeRepository().FindAggregatedReputationChanges(paginationParams, filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, paginationParams.PageSize, reputations), nil
}
