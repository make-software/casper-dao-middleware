package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/sse"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/apps/handler/config"
	"casper-dao-middleware/apps/handler/handlers"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/settings"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/assert"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	assert.OK(container.Provide(func(cfg *config.Env) (*sqlx.DB, error) {
		return boot.InitMySQL(ctx, cfg.DBConfig)
	}))

	defer container.Invoke(func(dbConn *sqlx.DB) {
		boot.CloseMySQL(dbConn)
	})

	assert.OK(container.Provide(func(cfg *config.Env) casper.RPCClient {
		handler := casper.NewRPCHandler(cfg.NodeRPCURL.String(), &http.Client{
			Timeout: 20 * time.Second,
		})

		return casper.NewRPCClient(handler)
	}))

	assert.OK(container.Provide(func(cfg *config.Env, rpcClient casper.RPCClient) (utils.DAOContractsMetadata, error) {
		return utils.NewDAOContractsMetadata(cfg.DaoContracts, rpcClient)
	}))

	//nolint:gocritic
	assert.OK(container.Provide(func(db *sqlx.DB, hashes utils.DAOContractsMetadata) persistence.EntityManager {
		return persistence.NewEntityManager(db, hashes)
	}))

	assert.OK(container.Invoke(func(env *config.Env, entityManager persistence.EntityManager, casperClient casper.RPCClient, metadata utils.DAOContractsMetadata) error {
		syncDaoSetting := settings.NewSyncDAOSettings()
		syncDaoSetting.SetCasperClient(casperClient)
		syncDaoSetting.SetVariableRepositoryContractStorageUref(metadata.VariableRepositoryContractStorageUref)
		syncDaoSetting.SetEntityManager(entityManager)
		syncDaoSetting.SetSettings(settings.VariableRepoSettings)
		syncDaoSetting.Execute()

		cesParser, err := ces.NewParser(casperClient, metadata.ContractHashes())
		if err != nil {
			zap.S().With(zap.Error(err)).Fatal("Failed to create CES Parser")
		}

		connection := sse.NewHttpConnection(&http.Client{Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 30,
		}}, env.NodeSSEURL.String()+env.EventStreamPath)
		streamReader := &sse.EventStreamReader{MaxBufferSize: 1024 * 1024 * 50} // 50 MB

		client := sse.NewClient(connection.URL)
		client.Streamer = sse.NewStreamer(connection, streamReader, 1*time.Minute)
		client.RegisterHandler(sse.DeployProcessedEventType, handlers.NewDeployProcessed(entityManager, casperClient, metadata, cesParser).Handle)

		client.EventStream = make(chan sse.RawEvent, 10)
		client.WorkersCount = 1

		defer client.Stop()
		return client.Start(ctx, 0)
	}))
}
