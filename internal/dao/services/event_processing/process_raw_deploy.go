package event_processing

import (
	"go.uber.org/zap"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/di"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware

	cesParser *ces.EventParser
}

func NewProcessRawDeploy() ProcessRawDeploy {
	return ProcessRawDeploy{}
}

func (c *ProcessRawDeploy) SetCESEventParser(parser *ces.EventParser) {
	c.cesParser = parser
}

func (c *ProcessRawDeploy) Execute() error {
	deployProcessedEvent := c.GetDeployProcessedEvent()
	daoContractsMetadata := c.GetDAOContractsMetadata()

	results, err := c.cesParser.ParseExecutionResults(deployProcessedEvent.DeployProcessed.ExecutionResult)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error != nil {
			zap.S().With(zap.Error(err)).Error("Failed to parse ces events")
			return err
		}
	}

	processContractEvents := NewProcessContractEvents()
	processContractEvents.SetDAOContractsMetadata(daoContractsMetadata)
	processContractEvents.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
	processContractEvents.SetEntityManager(c.GetEntityManager())

	for _, result := range results {
		processContractEvents.SetCESEvent(result.Event)
		if err := processContractEvents.Execute(); err != nil {
			zap.S().With(zap.Error(err)).With("event", result.Event.Name).Error("Failed to process ces event")
		}

		zap.S().With("event", result.Event.Name).Info("Successfully tracked event")
	}

	return nil
}
