package voting

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetVotings struct {
	di.PaginationParamsAware
	di.EntityManagerAware

	votingIDs []uint32
	isFormal  *bool
	hasEnded  *bool
}

func NewGetVotings() *GetVotings {
	return &GetVotings{}
}

func (c *GetVotings) SetVotingIDs(ids []uint32) {
	c.votingIDs = ids
}

func (c *GetVotings) SetIsFormal(isFormal *bool) {
	c.isFormal = isFormal
}

func (c *GetVotings) SetHasEnded(hasEnded *bool) {
	c.hasEnded = hasEnded
}

func (c *GetVotings) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	if c.hasEnded != nil {
		filters["has_ended"] = *c.hasEnded
	}

	if c.isFormal != nil {
		filters["is_formal"] = *c.isFormal
	}

	if len(c.votingIDs) != 0 {
		filters["voting_id"] = c.votingIDs
	}

	count, err := c.GetEntityManager().VotingRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	votings, err := c.GetEntityManager().VotingRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, votings), nil
}
