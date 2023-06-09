package config

import (
	"fmt"
	"math/big"
	"net/url"
	"path/filepath"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/pkg/errors"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap/zapcore"

	"casper-dao-middleware/pkg/config"
)

type Env struct {
	NodeRPCURL                *url.URL
	SetRateDeployerPrivateKey casper.PrivateKey

	LogLevel                      zapcore.Level       `env:"LOG_LEVEL" envDefault:"info"`
	SetRateCallPaymentAmount      *big.Int            `env:"SET_RATE_CALL_PAYMENT_AMOUNT,required"`
	SetRateDeployerPrivateKeyPath string              `env:"SET_RATE_DEPLOYER_PRIVATE_KEY_PATH,required"`
	NetworkName                   string              `env:"NETWORK_NAME,required"`
	RateApiURL                    string              `env:"RATE_API_URL,required"`
	CSPRRateProviderContractHash  casper.ContractHash `env:"CSPR_RATE_PROVIDER_CONTRACT_HASH,required"`
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

	e.SetRateDeployerPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(e.SetRateDeployerPrivateKeyPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse private key from %s", filepath.Dir(e.SetRateDeployerPrivateKeyPath))
	}

	return nil
}
