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
		SetDeployProcessed(deployProcessed casper.DeployProcessed)
	}

	for _, event := range daoEvents {
		switch event.EventName {
		case events.VotingCreatedEventName:
			daoEventHandler = event_tracking.NewTrackVotingCreated()
		case events.MintEventName:
			trackMintEvent := event_tracking.NewTrackMint()
			trackMintEvent.SetEventContractPackage(c.daoContractPackageHashes.ReputationContractPackageHash)
			daoEventHandler = trackMintEvent
		case events.BurnEventName:
			trackBurnEvent := event_tracking.NewTrackBurn()
			trackBurnEvent.SetEventContractPackage(c.daoContractPackageHashes.ReputationContractPackageHash)
			daoEventHandler = trackBurnEvent
		case events.BallotCastName:
			trackBallotCast := event_tracking.NewTrackBallotCast()
			trackBallotCast.SetDAOContractPackageHashes(c.daoContractPackageHashes)
			daoEventHandler = trackBallotCast
		case events.VotingEndedEventName:
			trackVotingEnded := event_tracking.NewTrackVotingEnded()
			trackVotingEnded.SetEventContractPackage(c.daoContractPackageHashes.VoterContractPackageHash)
			daoEventHandler = trackVotingEnded
		case "AddedToWhitelist":
			zap.S().Debug("new AddedToWhitelist event received")
			continue
		default:
			return errors.New("unsupported DAO event")
		}

		daoEventHandler.SetDeployProcessed(c.deployProcessedEvent.DeployProcessed)
		daoEventHandler.SetEventBody(event.EventBody)
		daoEventHandler.SetEntityManager(c.GetEntityManager())
		if err := daoEventHandler.Execute(); err != nil {
			zap.S().With(zap.Error(err)).With(zap.String("event-name", event.EventName)).Info("Failed to handle DAO event")
			return err
		}
	}

	return nil
}
