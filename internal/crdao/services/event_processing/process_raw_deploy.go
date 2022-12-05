package event_processing

import (
	"errors"

	"casper-dao-middleware/internal/crdao/dao_event_parser"
	"casper-dao-middleware/internal/crdao/dao_event_parser/events"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/event_tracking"
	"casper-dao-middleware/pkg/casper"

	"go.uber.org/zap"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string

	deployProcessedEvent     *casper.DeployProcessedEvent
	daoEventsParser          *dao_event_parser.DaoEventParser
	daoContractPackageHashes dao_event_parser.DAOContractPackageHashes
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

func (c *ProcessRawDeploy) SetDAOEventParser(parser *dao_event_parser.DaoEventParser) {
	c.daoEventsParser = parser
}

func (c *ProcessRawDeploy) SetDAOContractPackageHashes(hashes dao_event_parser.DAOContractPackageHashes) {
	c.daoContractPackageHashes = hashes
}

func (c *ProcessRawDeploy) Execute() error {
	daoEvents, err := c.daoEventsParser.Parse(c.deployProcessedEvent)
	if err != nil {
		return err
	}

	var daoEventHandler interface {
		Execute() error

		SetEntityManager(manager persistence.EntityManager)
		SetEventBody(eventBody []byte)
	}

	for _, event := range daoEvents {
		switch event.EventName {
		case events.VotingCreatedEventName:
			trackVotingCreated := event_tracking.NewTrackVotingCreated()
			trackVotingCreated.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			daoEventHandler = trackVotingCreated
		case events.MintEventName:
			trackMintEvent := event_tracking.NewTrackMint()
			trackMintEvent.SetEventContractPackage(c.daoContractPackageHashes.ReputationContractPackageHash)
			trackMintEvent.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			daoEventHandler = trackMintEvent
		case events.BurnEventName:
			trackBurnEvent := event_tracking.NewTrackBurn()
			trackBurnEvent.SetEventContractPackage(c.daoContractPackageHashes.ReputationContractPackageHash)
			trackBurnEvent.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			daoEventHandler = trackBurnEvent
		case events.BallotCastName:
			trackBallotCast := event_tracking.NewTrackBallotCast()
			trackBallotCast.SetDAOContractPackageHashes(c.daoContractPackageHashes)
			trackBallotCast.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			daoEventHandler = trackBallotCast
		case events.VotingEndedEventName:
			trackVotingEnded := event_tracking.NewTrackVotingEnded()
			trackVotingEnded.SetEventContractPackage(c.daoContractPackageHashes.VoterContractPackageHash)
			trackVotingEnded.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
			daoEventHandler = trackVotingEnded
		case events.ValueUpdatedEventName:
			trackValueUpdated := event_tracking.NewTrackValueUpdated()
			trackValueUpdated.SetVariableRepositoryContractStorageUref(c.variableRepositoryContractStorageUref)
			trackValueUpdated.SetCasperClient(c.GetCasperClient())
			daoEventHandler = trackValueUpdated
		case "AddedToWhitelist":
			zap.S().Debug("new AddedToWhitelist event received")
			continue
		default:
			return errors.New("unsupported DAO event")
		}

		daoEventHandler.SetEventBody(event.EventBody)
		daoEventHandler.SetEntityManager(c.GetEntityManager())
		if err := daoEventHandler.Execute(); err != nil {
			zap.S().With(zap.Error(err)).With(zap.String("event-name", event.EventName)).Info("Failed to handle DAO event")
			return err
		}
	}

	return nil
}
