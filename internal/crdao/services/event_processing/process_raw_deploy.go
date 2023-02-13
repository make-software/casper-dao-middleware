package event_processing

import (
	"errors"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/internal/crdao/services/event_tracking"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string
	cesParser                             *ces.EventParser

	deployProcessedEvent *casper.DeployProcessedEvent
}

func NewProcessRawDeploy() ProcessRawDeploy {
	return ProcessRawDeploy{}
}

func (c *ProcessRawDeploy) SetDeployProcessedEvent(event *casper.DeployProcessedEvent) {
	c.deployProcessedEvent = event
}

func (c *ProcessRawDeploy) SetVariableRepositoryContractStorageUref(uref string) {
	c.variableRepositoryContractStorageUref = uref
}

func (c *ProcessRawDeploy) SetCESEventParser(parser *ces.EventParser) {
	c.cesParser = parser
}

func (c *ProcessRawDeploy) Execute() error {
	results, err := c.cesParser.ParseExecutionResults(c.deployProcessedEvent.DeployProcessed.ExecutionResult)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error != nil {
			zap.S().With(zap.Error(err)).Error("Failed to parse ces events")
			return err
		}
	}

	for _, result := range results {
		switch result.Event.Name {
		case events.SimpleVotingCreatedEventName:
			trackVotingCreated := event_tracking.NewTrackVotingCreated()
			trackVotingCreated.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackVotingCreated.SetCESEvent(result.Event)
			trackVotingCreated.SetEntityManager(c.GetEntityManager())
			if err := trackVotingCreated.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case events.BallotCastName:
			trackBallotCast := event_tracking.NewTrackBallotCast()
			//trackBallotCast.SetDAOContractsMetadata(c.daoContractPackageHashes)
			trackBallotCast.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackBallotCast.SetCESEvent(result.Event)
			trackBallotCast.SetEntityManager(c.GetEntityManager())
			if err := trackBallotCast.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case "AddedToWhitelist":
			zap.S().Debug("new AddedToWhitelist event received")
			continue
		default:
			return errors.New("unsupported DAO event")
		}
	}

	return nil
}
