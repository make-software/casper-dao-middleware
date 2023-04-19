package bid

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetBids struct {
	di.PaginationParamsAware
	di.EntityManagerAware

	jobOfferID uint32
}

func NewGetBids() *GetBids {
	return &GetBids{}
}

func (c *GetBids) SetJobOfferID(jobOfferID uint32) {
	c.jobOfferID = jobOfferID
}

func (c *GetBids) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{
		"job_offer_id": c.jobOfferID,
	}

	count, err := c.GetEntityManager().BidRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	bids, err := c.GetEntityManager().BidRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, bids), nil
}
