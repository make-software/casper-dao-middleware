package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/event_processing"
	"casper-dao-middleware/internal/crdao/services/settings"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/config"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	NodeRPCURL        *url.URL
	DaoContractHashes map[string]string
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

type PopulateCrDAODeploysFromClarity struct {
	cfg                Env
	clarityDB, crDAODB *sqlx.DB
	casperClient       casper.RPCClient

	daoEventParser           *dao.DaoEventParser
	daoContractPackageHashes dao.DAOContractsMetadata
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

	c.daoEventParser, err = dao.NewDaoEventParser(c.casperClient, cfg.DaoContractHashes, 100)
	if err != nil {
		return err
	}

	c.daoContractPackageHashes, err = dao.NewDAOContractsMetadataFromHashesMap(cfg.DaoContractHashes, c.casperClient)
	return err
}

func (c *PopulateCrDAODeploysFromClarity) Execute() error {
	daoDeploysCursor := c.createDAODeployCursor(c.clarityDB, c.cfg.DaoContractHashes)

	crdaoEntityManager := persistence.NewEntityManager(c.crDAODB, c.daoContractPackageHashes)

	syncDaoSetting := settings.NewSyncDAOSettings()
	syncDaoSetting.SetCasperClient(c.casperClient)
	syncDaoSetting.SetVariableRepositoryContractStorageUref(c.daoContractPackageHashes.VariableRepositoryContractStorageUref)
	syncDaoSetting.SetEntityManager(crdaoEntityManager)
	syncDaoSetting.SetSettings(settings.DaoSettings)
	syncDaoSetting.Execute()

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetDAOEventParser(c.daoEventParser)
	processRawDeploy.SetCasperClient(c.casperClient)
	processRawDeploy.SetEntityManager(crdaoEntityManager)
	processRawDeploy.SetDAOContractPackageHashes(c.daoContractPackageHashes)

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

		processRawDeploy.SetDeployProcessedEvent(&casper.DeployProcessedEvent{
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
			return fmt.Errorf("failed to process rawDeploy %s", err.Error())
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

func (c *PopulateCrDAODeploysFromClarity) createDAODeployCursor(clarityDB *sqlx.DB, daoContractHashes map[string]string) *sqlx.Rows {
	contractParams := make([]string, 0, len(daoContractHashes))
	contracts := make([]interface{}, 0, len(daoContractHashes))
	for _, contractHash := range daoContractHashes {
		contractParams = append(contractParams, `unhex(?)`)
		contracts = append(contracts, contractHash)
	}

	query := fmt.Sprintf(`
		select deploy_hash from extended_deploys 
        where contract_hash in (%s) order by timestamp desc;`, strings.Join(contractParams, ","))

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
