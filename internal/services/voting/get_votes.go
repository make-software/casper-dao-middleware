package voting

import (
	"casper-dao-middleware/internal/di"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"
)

type GetVotes struct {
	di.PaginationParamsAware
	di.EntityManagerAware

	votingIDs []uint32
	address   *types.Hash
}

func NewGetVotes() *GetVotes {
	return &GetVotes{}
}

func (c *GetVotes) SetVotingIDs(votingIDs []uint32) {
	c.votingIDs = votingIDs
}

func (c *GetVotes) SetAddress(address *types.Hash) {
	c.address = address
}

func (c *GetVotes) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	if len(c.votingIDs) != 0 {
		filters["voting_id"] = c.votingIDs
	}

	if c.address != nil {
		filters["address"] = *c.address
	}

	count, err := c.GetEntityManager().VoteRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	votes, err := c.GetEntityManager().VoteRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, votes), nil
}
