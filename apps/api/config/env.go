package config

import (
	"fmt"
	"net/url"

	"casper-dao-middleware/pkg/config"
	"casper-dao-middleware/pkg/http"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap/zapcore"
)

type Env struct {
	Addr     http.ServerAddress `env:"ADDRESS,required"`
	LogLevel zapcore.Level      `env:"LOG_LEVEL" envDefault:"info"`

	NodeRPCURL *url.URL

	DBConfig     config.DBConfig
	DaoContracts config.DaoContracts
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

	return nil
}
