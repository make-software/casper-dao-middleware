package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/pelletier/go-toml/v2"

	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/config"
)

type Contract struct {
	Name         string `toml:"name"`
	PackageHash  string `toml:"package_hash"`
	ContractHash string `toml:"contract_hash"`
}

type DeployResult struct {
	Contracts []Contract `toml:"contracts"`
}

type Env struct {
	NodeRPCURL         *url.URL
	TomlResultFilePath string `env:"DEPLOY_RESULT_TOML_FILE_PATH"`
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

type GatherContractHashesFromToml struct {
	cfg          Env
	casperClient casper.RPCClient
}

func (c *GatherContractHashesFromToml) SetUp() error {
	cfg := Env{}
	err := boot.ParseEnvConfig(&cfg)
	if err != nil {
		return err
	}
	c.cfg = cfg
	handler := casper.NewRPCHandler(cfg.NodeRPCURL.String(), &http.Client{
		Timeout: 20 * time.Second,
	})
	c.casperClient = casper.NewRPCClient(handler)
	return nil
}

func (c *GatherContractHashesFromToml) Execute() error {
	doc, err := os.ReadFile(c.cfg.TomlResultFilePath)
	if err != nil {
		panic(err)
	}

	var result DeployResult
	if err = toml.Unmarshal(doc, &result); err != nil {
		panic(err)
	}

	ctx := context.Background()
	for i := range result.Contracts {
		log.Printf("Searching for ContractPackageHash - %s \n", result.Contracts[i].PackageHash)
		res, err := c.casperClient.QueryGlobalStateByStateHash(ctx, nil, result.Contracts[i].PackageHash, nil)
		if err != nil {
			panic(err)
		}
		result.Contracts[i].ContractHash = res.StoredValue.ContractPackage.Versions[0].Hash.ToHex()
	}

	marshalled, err := toml.Marshal(result)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(c.cfg.TomlResultFilePath, marshalled, 06444); err != nil {
		panic(err)
	}

	log.Println("Processing finished successfully")
	return nil
}
func (c *GatherContractHashesFromToml) TearDown() error {
	return nil
}

func main() {
	command := new(GatherContractHashesFromToml)
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
