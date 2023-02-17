package event_processing

import (
	"errors"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/internal/crdao/services/event_tracking"
	"casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware

	cesParser            *ces.EventParser
	daoContractsMetadata types.DAOContractsMetadata
	deployProcessedEvent *casper.DeployProcessedEvent
}

func NewProcessRawDeploy() ProcessRawDeploy {
	return ProcessRawDeploy{}
}

func (c *ProcessRawDeploy) SetDeployProcessedEvent(event *casper.DeployProcessedEvent) {
	c.deployProcessedEvent = event
}

func (c *ProcessRawDeploy) SetDAOContractsMetadata(daoContractsMetadata types.DAOContractsMetadata) {
	c.daoContractsMetadata = daoContractsMetadata
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
		case events.SimpleVotingCreatedEventName, events.ReputationVotingCreatedEventName, events.RepoVotingCreated, events.KYCVotingCreated:
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
			trackBallotCast.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackBallotCast.SetCESEvent(result.Event)
			trackBallotCast.SetEntityManager(c.GetEntityManager())
			trackBallotCast.SetDAOContractsMetadata(c.daoContractsMetadata)
			if err := trackBallotCast.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case events.Transfer:
			trackTransfer := event_tracking.NewTrackTransfer()
			trackTransfer.SetDAOContractsMetadata(c.daoContractsMetadata)
			trackTransfer.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackTransfer.SetCESEvent(result.Event)
			if err := trackTransfer.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}

		case events.MintEventName:
			trackMintEvent := event_tracking.NewTrackMint()
			trackMintEvent.SetEventContractPackage(c.daoContractsMetadata.ReputationContractPackageHash)
			trackMintEvent.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackMintEvent.SetCESEvent(result.Event)
			if err := trackMintEvent.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case events.BurnEventName:
			trackBurnEvent := event_tracking.NewTrackBurn()
			trackBurnEvent.SetEventContractPackage(c.daoContractsMetadata.ReputationContractPackageHash)
			trackBurnEvent.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			trackBurnEvent.SetCESEvent(result.Event)
			if err := trackBurnEvent.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		default:
			return errors.New("unsupported DAO event")
		}
	}

	return nil
}
