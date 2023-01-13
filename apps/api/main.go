package main

import (
	"context"
	"log"

	"casper-dao-middleware/apps/api/config"
	"casper-dao-middleware/internal/crdao/dao_event_parser"
	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/pkg/assert"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/exec"
	"casper-dao-middleware/pkg/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title   Casper-CRDao API
// @version 0.0.1
// @Produce json

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

	assert.OK(container.Provide(func(cfg *config.Env) (dao_event_parser.DAOContractsMetadata, error) {
		return dao_event_parser.NewDAOContractsMetadataFromHashesMap(cfg.DaoContractHashes, casper.NewRPCClient(cfg.NodeRPCURL.String()))
	}))

	assert.OK(container.Provide(func(db *sqlx.DB, hashes dao_event_parser.DAOContractsMetadata) persistence.EntityManager {
		return persistence.NewEntityManager(db, hashes)
	}))

	assert.OK(container.Provide(NewRouter))
	assert.OK(container.Provide(func(cfg *config.Env) http.ServerAddress { return cfg.Addr }))
	assert.OK(container.Provide(http.NewServer))

	assert.OK(container.Invoke(func(server http.Server) {
		if err := server.ListenAndServe(ctx); err != nil {
			log.Printf("Failed to serve api %s", err.Error())
		}
	}))
}
