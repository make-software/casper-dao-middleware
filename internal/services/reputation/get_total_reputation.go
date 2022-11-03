package reputation

import (
	"casper-dao-middleware/internal/di"
	"casper-dao-middleware/internal/entities"
	"casper-dao-middleware/pkg/casper/types"
)

type GetTotalReputation struct {
	di.EntityManagerAware
	di.DAOContractPackageHashesAware

	addressHash types.Hash
}

func NewGetTotalReputation() *GetTotalReputation {
	return &GetTotalReputation{}
}

func (c *GetTotalReputation) SetAddressHash(hash types.Hash) {
	c.addressHash = hash
}

func (c *GetTotalReputation) Execute() (entities.TotalReputation, error) {
	return c.GetEntityManager().ReputationChangeRepository().CalculateTotalReputationForAddress(c.addressHash)
}
