package repo_voter

import (
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/repo_voter"
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
	case repo_voter.VotingCreatedEventName:
		trackVotingCreated := NewTrackVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case repo_voter.VotingEndedEventName:
		trackVotingEnded := NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case repo_voter.VotingCanceledEventName:
		trackVotingCanceled := NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case repo_voter.BallotCastEventName:
		trackBallotCast := NewTrackBallotCast()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case repo_voter.BallotCanceledEventName:
		trackBallotCanceled := NewTrackBallotCanceled()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
