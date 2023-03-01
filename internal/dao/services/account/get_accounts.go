package account

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetAccounts struct {
	di.PaginationParamsAware
	di.EntityManagerAware
}

func NewGetAccounts() *GetAccounts {
	return &GetAccounts{}
}

func (c *GetAccounts) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	count, err := c.GetEntityManager().AccountRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	votes, err := c.GetEntityManager().AccountRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, votes), nil
}
