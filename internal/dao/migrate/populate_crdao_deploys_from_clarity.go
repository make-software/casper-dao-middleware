package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/utils"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/config"
)

// MigrateCommand common migration command interface
// Maybe we can extract it in some common place to have unified interface for all Migration commands/scripts
type MigrateCommand interface {
	SetUp() error
	Execute() error
	TearDown() error
}
type Env struct {
	ClarityDBConfig config.DBConfig `envPrefix:"CLARITY_"`
	CrDAODBConfig   config.DBConfig `envPrefix:"CRDAO_"`
	NodeRPCURL      *url.URL
	DaoContracts    config.DaoContracts
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

type PopulateCrDAODeploysFromClarity struct {
	cfg                Env
	clarityDB, crDAODB *sqlx.DB
	casperClient       casper.RPCClient

	daoContractsMetadata utils.DAOContractsMetadata
}

func (c *PopulateCrDAODeploysFromClarity) SetUp() error {
	cfg := Env{}
	err := boot.ParseEnvConfig(&cfg)
	if err != nil {
		return err
	}
	c.cfg = cfg
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	c.clarityDB, err = boot.InitMySQL(ctx, cfg.ClarityDBConfig)
	if err != nil {
		return err
	}
	c.crDAODB, err = boot.InitMySQL(ctx, cfg.CrDAODBConfig)
	if err != nil {
		return err
	}

	c.casperClient = casper.NewRPCClient(cfg.NodeRPCURL.String())

	c.daoContractsMetadata, err = utils.NewDAOContractsMetadata(cfg.DaoContracts, c.casperClient)
	if err != nil {
		return err
	}
	return nil
}

func (c *PopulateCrDAODeploysFromClarity) Execute() error {
	daoDeploysCursor := c.createDAODeployCursor(c.clarityDB, c.daoContractsMetadata)

	crdaoEntityManager := persistence.NewEntityManager(c.crDAODB, c.daoContractsMetadata)

	cesParser, err := ces.NewParser(c.casperClient, c.daoContractsMetadata.ContractHashes())
	if err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to create CES Parser")
		return err
	}

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(crdaoEntityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(c.daoContractsMetadata)

	for daoDeploysCursor.Next() {
		var rawDeployHash string
		if err := daoDeploysCursor.Scan(&rawDeployHash); err != nil {
			return fmt.Errorf("failed to scan deploy hash: %s", err.Error())
		}
		deployHash, err := types.NewHashFromRawBytes([]byte(rawDeployHash))
		if err != nil {
			return err
		}
		deploy, err := c.casperClient.GetDeploy(deployHash.ToHex())
		if err != nil {
			return fmt.Errorf("failed to get deploy by hash: %s", err.Error())
		}
		if deploy.ExecutionResults[0].Result.Failure != nil {
			log.Println("Failed deploy, skipping: ", deploy.Deploy.Hash.ToHex())
			continue
		}
		processRawDeploy.SetDeployProcessedEvent(casper.DeployProcessedEvent{
			DeployProcessed: casper.DeployProcessed{
				DeployHash:      deploy.Deploy.Hash,
				Account:         deploy.Deploy.Header.Account.ToHex(),
				Timestamp:       deploy.Deploy.Header.Timestamp,
				TTL:             deploy.Deploy.Header.TTL,
				BlockHash:       deploy.Deploy.Header.TTL,
				ExecutionResult: deploy.ExecutionResults[0].Result,
			},
		})
		if err := processRawDeploy.Execute(); err != nil {
			log.Printf("failed to process rawDeploy %s\n", err.Error())
		}
		log.Println("Process DAO Deploy ", deploy.Deploy.Hash.ToHex())
	}
	log.Println("Processing finished successfully")
	return nil
}
func (c *PopulateCrDAODeploysFromClarity) TearDown() error {
	c.crDAODB.Close()
	c.clarityDB.Close()
	return nil
}
func (c *PopulateCrDAODeploysFromClarity) createDAODeployCursor(clarityDB *sqlx.DB, hashes utils.DAOContractsMetadata) *sqlx.Rows {
	daoContracts := hashes.ContractHashes()

	contractParams := make([]string, 0, len(daoContracts))
	contracts := make([]interface{}, 0, len(daoContracts))
	for _, contractHash := range daoContracts {
		contractParams = append(contractParams, `unhex(?)`)
		contracts = append(contracts, contractHash.ToHex())
	}
	query := fmt.Sprintf(`
		select deploy_hash from extended_deploys 
        where contract_hash in (%s) order by timestamp;`, strings.Join(contractParams, ","))
	daoDeploysCursor, err := clarityDB.Queryx(query, contracts...)
	if err != nil {
		log.Fatalln(err)
	}
	return daoDeploysCursor
}
func main() {
	runCommand(new(PopulateCrDAODeploysFromClarity))
}
func runCommand(command MigrateCommand) {
	if err := command.SetUp(); err != nil {
		log.Fatalf("Command initialization failed: %s", err.Error())
	}
	if err := command.Execute(); err != nil {
		log.Fatalf("Command execution failed: %s", err.Error())
	}
	if err := command.TearDown(); err != nil {
		log.Fatalf("Command teardown failed: %s", err.Error())
	}
}
