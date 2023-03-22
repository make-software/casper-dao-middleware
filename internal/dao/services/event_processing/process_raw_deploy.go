package event_processing

import (
	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/services/event_tracking"
	"casper-dao-middleware/pkg/go-ces-parser"
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

	trackContract := event_tracking.NewTrackContract()
	trackContract.SetDAOContractsMetadata(daoContractsMetadata)
	trackContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
	trackContract.SetEntityManager(c.GetEntityManager())

	for _, result := range results {
		trackContract.SetCESEvent(result.Event)
		if err := trackContract.Execute(); err != nil {
			return err
		}

		zap.S().With("event", result.Event.Name).Info("Successfully tracked event")
	}

	return nil
}
