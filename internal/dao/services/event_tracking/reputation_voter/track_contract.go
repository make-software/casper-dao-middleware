package reputation_voter

import (
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	base_events "casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/events/reputation_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/base"
)

type TrackContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackContract() *TrackContract {
	return &TrackContract{}
}

func (s *TrackContract) Execute() error {
	cesEvent := s.GetCESEvent()
	doaContractMetadata := s.GetDAOContractsMetadata()

	switch cesEvent.Name {
	case reputation_voter.VotingCreatedEventName:
		trackVotingCreated := NewTrackVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := NewTrackBallotCast()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := base.NewTrackBallotCanceled()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
