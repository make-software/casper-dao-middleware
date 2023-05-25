package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/apps/rate-oracle/config"
	"casper-dao-middleware/internal/dao/services/rate"
	"casper-dao-middleware/pkg/assert"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/exec"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	container := dig.New()

	ctx, cancel := context.WithCancel(context.Background())
	exec.RunGracefulShutDownListener(ctx, cancel)

	assert.OK(container.Provide(func() *config.Env {
		cfg := config.Env{}

		if err := boot.ParseEnvConfig(&cfg); err != nil {
			log.Fatal(err)
		}
		return &cfg
	}))

	// we should provide log level to invoke deps.InitLogger method
	assert.OK(container.Provide(func(cfg *config.Env) zapcore.Level {
		return cfg.LogLevel
	}))

	assert.OK(container.Invoke(boot.NewLogger))
	defer zap.S().Sync()

	assert.OK(container.Provide(func(cfg *config.Env) casper.RPCClient {
		handler := casper.NewRPCHandler(cfg.NodeRPCURL.String(), &http.Client{
			Timeout: 20 * time.Second,
		})

		return casper.NewRPCClient(handler)
	}))

	assert.OK(container.Invoke(func(env *config.Env, casperClient casper.RPCClient) error {
		syncRate := rate.NewSyncRates()
		syncRate.SetSyncDuration(env.RateSyncDuration)
		syncRate.SetCasperClient(casperClient)
		syncRate.SetRateDeployerPrivateKey(env.SetRateDeployerPrivateKey)
		syncRate.SetRateAPIUrl(env.RateApiURL)
		syncRate.SetContractExecutionAmount(env.SetRateCallPaymentAmount.Int64())
		syncRate.SetNetworkName(env.NetworkName)
		syncRate.SetCSPRRatesProviderContractHash(env.CSPRRateProviderContractHash)

		return syncRate.Execute(ctx)
	}))
}
