package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"casper-dao-middleware/pkg/config"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap/zapcore"
)

type Env struct {
	LogLevel                  zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`
	EventStreamPath           string        `env:"EVENT_STREAM_PATH,required"`
	DictionarySetEventsBuffer uint32        `env:"DICTIONARY_SET_EVENTS_READ_BACK_BUFFER" envDefault:"100"`
	DaoContractHashes         map[string]string
	NewNodeStartFromEventID   uint64

	NodeSSEURL *url.URL
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

	e.NodeSSEURL, err = url.Parse(fmt.Sprintf("http://%s:%s", config.GetEnv("NODE_ADDRESS"),
		config.GetEnv("NODE_PORT")))
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

	eventID := os.Getenv(fmt.Sprintf("NEW_NODE_START_FROM_EVENT_ID_%s",
		strings.ReplaceAll(e.NodeSSEURL.Hostname(), ".", "_")))
	if eventID != "" {
		e.NewNodeStartFromEventID, err = strconv.ParseUint(eventID, 10, 0)
		if err != nil {
			return err
		}
	}

	return nil
}
