package voting

import (
	"casper-dao-middleware/internal/dao/di"
)

type GetVotesNumber struct {
	di.EntityManagerAware

	votingIDs []uint32
}

func NewGetVotesNumber() *GetVotesNumber {
	return &GetVotesNumber{}
}

func (c *GetVotesNumber) SetVotingIDs(ids []uint32) {
	c.votingIDs = ids
}

func (c *GetVotesNumber) Execute() (map[uint32]uint32, error) {
	return c.GetEntityManager().VoteRepository().CountVotesNumberForVotings(c.votingIDs)
}
