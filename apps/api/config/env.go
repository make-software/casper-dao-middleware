package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"casper-dao-middleware/pkg/config"
	"casper-dao-middleware/pkg/http"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap/zapcore"
)

type Env struct {
	Addr              http.ServerAddress `env:"ADDRESS,required"`
	LogLevel          zapcore.Level      `env:"LOG_LEVEL" envDefault:"info"`
	DaoContractHashes map[string]string

	NodeRPCURL *url.URL
	DBConfig   config.DBConfig
}

func (e *Env) Parse() error {
	err := env.Parse(e)
	if err != nil {
		return err
	}

	e.NodeRPCURL, err = url.Parse(fmt.Sprintf("http://%s:%s/rpc", config.GetEnv("NODE_ADDRESS"),
		config.GetEnv("NODE_RPC_PORT")))
	if err != nil {
		return err
	}

	e.DaoContractHashes = make(map[string]string, 0)
	for _, contract := range strings.Split(config.GetEnv("DAO_CONTRACT_HASHES"), ",") {
		definitions := strings.Split(contract, ":")
		//expect contract_name:contract_hash
		if len(definitions) != 2 {
			return errors.New("invalid DAO_CONTRACT_HASHES format provided")
		}
		e.DaoContractHashes[definitions[0]] = definitions[1]
	}

	return nil
}
