package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/pkg/config"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap/zapcore"
)

type Env struct {
	LogLevel                      zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`
	EventStreamPath               string        `env:"EVENT_STREAM_PATH,required"`
	DictionarySetEventsBuffer     uint32        `env:"DICTIONARY_SET_EVENTS_READ_BACK_BUFFER" envDefault:"100"`
	NewNodeStartFromEventID       uint64
	VariableRepoInstallDeployHash casper.Hash `env:"VARIABLE_REPOSITORY_INSTALL_DEPLOY_HASH,required" `

	NodeSSEURL *url.URL
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

	e.NodeSSEURL, err = url.Parse(fmt.Sprintf("http://%s:%s", config.GetEnv("NODE_ADDRESS"),
		config.GetEnv("NODE_PORT")))
	if err != nil {
		return err
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
