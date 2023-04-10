package job_offer

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetJobOffers struct {
	di.PaginationParamsAware
	di.EntityManagerAware
}

func NewGetJobOffers() *GetJobOffers {
	return &GetJobOffers{}
}

func (c *GetJobOffers) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	count, err := c.GetEntityManager().JobOfferRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	offers, err := c.GetEntityManager().JobOfferRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, offers), nil
}
