package di

import "github.com/make-software/casper-go-sdk/casper"

type CasperClientAware struct {
	casperClient casper.RPCClient
}

func (t *CasperClientAware) SetCasperClient(casperClient casper.RPCClient) {
	t.casperClient = casperClient
}

func (t *CasperClientAware) GetCasperClient() casper.RPCClient {
	return t.casperClient
}
