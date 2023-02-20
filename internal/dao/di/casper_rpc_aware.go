package di

import (
	"casper-dao-middleware/pkg/casper"
)

type CasperClientAware struct {
	casperClient casper.RPCClient
}

func (t *CasperClientAware) SetCasperClient(casperClient casper.RPCClient) {
	t.casperClient = casperClient
}

func (t *CasperClientAware) GetCasperClient() casper.RPCClient {
	return t.casperClient
}
