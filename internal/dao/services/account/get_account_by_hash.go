package account

import (
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
)

//go:generate mockgen -destination=../../tests/mocks/rpc_client_mock.go -package=mocks -source=./get_account_by_hash.go RPCClient
type RPCClient interface {
	rpc.Client
}

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
