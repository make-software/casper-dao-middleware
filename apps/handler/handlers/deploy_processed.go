package handlers

import (
	"context"

	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
	"go.uber.org/zap"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/utils"
)

type DeployProcessed struct {
	entityManager persistence.EntityManager
	casperClient  rpc.Client
	daoMetadata   utils.DAOContractsMetadata
	cesParser     *ces.EventParser
}

func NewDeployProcessed(
	entityManager persistence.EntityManager,
	casperClient rpc.Client,
	daoMetadata utils.DAOContractsMetadata,
	cesParser *ces.EventParser,
) *DeployProcessed {
	return &DeployProcessed{
		entityManager: entityManager,
		casperClient:  casperClient,
		daoMetadata:   daoMetadata,
		cesParser:     cesParser,
	}
}

func (h DeployProcessed) Handle(ctx context.Context, event sse.RawEvent) error {
	deployProcessedEvent, err := event.ParseAsDeployProcessedEvent()
	if err != nil {
		return err
	}

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(h.entityManager)
	processRawDeploy.SetCESEventParser(h.cesParser)
	processRawDeploy.SetDAOContractsMetadata(h.daoMetadata)
	processRawDeploy.SetDeployProcessedEvent(deployProcessedEvent)
	if err = processRawDeploy.Execute(); err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to handle DeployProcessedEvent")
	}
	return nil
}
