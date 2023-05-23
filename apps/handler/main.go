package main

import (
	"context"
	"log"

	"casper-dao-middleware/apps/handler/config"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/assert"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//        - name: VARIABLE_REPOSITORY_CONTRACT_HASH
//          value: "0dc284f29d97cd5d413c7cc2ca0365ad1a80ea14defd94b95e3fe75d5dd21f97"
//        - name: REPUTATION_CONTRACT_HASH
//          value: "30ce55dea23aceb7b2c71b8cf053753dc7faa41fc5cb3e9e22ad3dee36000165"
//        - name: SIMPLE_VOTER_CONTRACT_HASH
//          value: "f0464e837a6ef0e06737c973a51fec8729064b98d2dd00a85281655a7376a1b8"
//        - name: REPO_VOTER_CONTRACT_HASH
//          value: "6a3213fe5db928dd4bb3d1c5ecd3bfbc68656823c9486ef389a3080921d0d3ec"
//        - name: REPUTATION_VOTER_CONTRACT_HASH
//          value: "714c06db9845ea4bea4865428ed4c897e15b0f26b49b183dfd89f3b0e8cc0234"
//        - name: SLASHING_VOTER_CONTRACT_HASH
//          value: "ab10c368b10673faf913c8a07128681dd6fac422a45ac437a96539046c602044"
//        - name: KYC_VOTER_CONTRACT_HASH
//          value: "ab5c87665982b3dce75ea314089fd1b13f61da71f985d0414370de58a0678358"
//        - name: VA_NFT_CONTRACT_HASH
//          value: "92bdae12487fd0102fd8a9066829894c87e87672ab7c0f6b1d45fa46919ffdf1"
//        - name: KYC_NFT_CONTRACT_HASH
//          value: "7f2af2f8144e6142f397abbca6749f008766dede823d49ebddf3e44aece8890e"
//        - name: ONBOARDING_REQUEST_CONTRACT_HASH
//          value: "8f60c7d74e52502eb8975e87ea911e14638b0cd023ebe1b9dfb1aa32b1368fc1"
//        - name: ADMIN_CONTRACT_HASH
//          value: "a3e7a4772d8b3e741a9fd7e6c2a089443c791c1e8836b44eb5695faa9269981b"
//        - name: BID_ESCROW_CONTRACT_HASH
//          value: "245c39646071c3099dcfdc8f466f3e269c591bca1bcf72cbcf09bf3d96b2fcc1"

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

	assert.OK(container.Provide(func(cfg *config.Env) (utils.DAOContractsMetadata, error) {
		return utils.NewDAOContractsMetadata(cfg.DaoContracts, casper.NewRPCClient(cfg.NodeRPCURL.String()))
	}))

	//nolint:gocritic
	assert.OK(container.Provide(func(db *sqlx.DB, hashes utils.DAOContractsMetadata) persistence.EntityManager {
		return persistence.NewEntityManager(db, hashes)
	}))

	assert.OK(container.Provide(func(cfg *config.Env) casper.RPCClient {
		return casper.NewRPCClient(cfg.NodeRPCURL.String())
	}))

	assert.OK(container.Invoke(func(env *config.Env, entityManager persistence.EntityManager, casperClient casper.RPCClient, metadata utils.DAOContractsMetadata) error {
		processEventStream := event_processing.NewProcessEventStream()
		processEventStream.SetBaseStreamURL(env.NodeSSEURL)
		processEventStream.SetNodeStartFromEventID(env.NewNodeStartFromEventID)
		processEventStream.SetEventStreamPath(env.EventStreamPath)
		processEventStream.SetCasperClient(casperClient)
		processEventStream.SetEntityManager(entityManager)
		processEventStream.SetDAOContractsMetadata(metadata)

		return processEventStream.Execute(ctx)
	}))
}
