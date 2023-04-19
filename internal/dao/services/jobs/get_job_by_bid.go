package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
)

type GetJobByBid struct {
	di.EntityManagerAware

	bidID uint32
}

func NewGetJobByBid() *GetJobByBid {
	return &GetJobByBid{}
}

func (c *GetJobByBid) SetBidID(bidID uint32) {
	c.bidID = bidID
}

func (c *GetJobByBid) Execute() (*entities.Job, error) {
	return c.GetEntityManager().JobRepository().GetByBidID(c.bidID)
}
