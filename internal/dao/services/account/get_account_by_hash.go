package account

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/casper/types"
)

type GetAccountByHash struct {
	di.EntityManagerAware

	accountHash types.Hash
}

func NewGetAccountByHash() *GetAccountByHash {
	return &GetAccountByHash{}
}

func (c *GetAccountByHash) SetHash(hash types.Hash) {
	c.accountHash = hash
}

func (c *GetAccountByHash) Execute() (*entities.Account, error) {
	return c.GetEntityManager().AccountRepository().FindByHash(c.accountHash)
}
