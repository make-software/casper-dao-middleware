package account

import (
	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
)

type GetAccountByHash struct {
	di.EntityManagerAware

	accountHash casper.Hash
}

func NewGetAccountByHash() *GetAccountByHash {
	return &GetAccountByHash{}
}

func (c *GetAccountByHash) SetHash(hash casper.Hash) {
	c.accountHash = hash
}

func (c *GetAccountByHash) Execute() (*entities.Account, error) {
	return c.GetEntityManager().AccountRepository().FindByHash(c.accountHash)
}
