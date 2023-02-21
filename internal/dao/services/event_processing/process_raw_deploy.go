package event_processing

import (
	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/config"
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events"
	"casper-dao-middleware/internal/dao/services/event_tracking"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware
	di.DeployProcessedEventAware

	cesParser            *ces.EventParser
	daoContractsMetadata config.DAOContractsMetadata
}

func NewProcessRawDeploy() ProcessRawDeploy {
	return ProcessRawDeploy{}
}

func (c *ProcessRawDeploy) SetDAOContractsMetadata(daoContractsMetadata config.DAOContractsMetadata) {
	c.daoContractsMetadata = daoContractsMetadata
}

func (c *ProcessRawDeploy) SetCESEventParser(parser *ces.EventParser) {
	c.cesParser = parser
}

func (c *ProcessRawDeploy) Execute() error {
	deployProcessedEvent := c.GetDeployProcessedEvent()
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

	for _, result := range results {
		switch result.Event.Name {
		case events.SimpleVotingCreatedEventName, events.ReputationVotingCreatedEventName, events.RepoVotingCreated, events.KYCVotingCreated:
			trackVotingCreated := event_tracking.NewTrackVotingCreated()
			trackVotingCreated.SetDeployProcessedEvent(deployProcessedEvent)
			trackVotingCreated.SetCESEvent(result.Event)
			trackVotingCreated.SetEntityManager(c.GetEntityManager())
			if err := trackVotingCreated.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case events.BallotCastName:
			trackBallotCast := event_tracking.NewTrackBallotCast()
			trackBallotCast.SetDeployProcessedEvent(deployProcessedEvent)
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
			trackTransfer.SetCESEvent(result.Event)
			if err := trackTransfer.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}

		case events.MintEventName:
			trackMintEvent := event_tracking.NewTrackMint()
			trackMintEvent.SetEventContractPackage(c.daoContractsMetadata.ReputationContractPackageHash)
			trackMintEvent.SetDeployProcessedEvent(deployProcessedEvent)
			trackMintEvent.SetCESEvent(result.Event)
			if err := trackMintEvent.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		case events.BurnEventName:
			trackBurnEvent := event_tracking.NewTrackBurn()
			trackBurnEvent.SetEventContractPackage(c.daoContractsMetadata.ReputationContractPackageHash)
			trackBurnEvent.SetDeployProcessedEvent(deployProcessedEvent)
			trackBurnEvent.SetCESEvent(result.Event)
			if err := trackBurnEvent.Execute(); err != nil {
				zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
				return err
			}
		}
	}

	return nil
}
