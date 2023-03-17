package reputation

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/casper/types"
)

type GetTotalReputation struct {
	di.EntityManagerAware
	di.DAOContractsMetadataAware

	addressHash types.Hash
}

func NewGetTotalReputation() *GetTotalReputation {
	return &GetTotalReputation{}
}

func (c *GetTotalReputation) SetAddressHash(hash types.Hash) {
	c.addressHash = hash
}

//func (c *GetTotalReputation) Execute() (entities.LiquidStateReputation, error) {
//	return c.GetEntityManager().ReputationChangeRepository().CalculateLiquidStakeReputationForAddress(c.addressHash)
//}
