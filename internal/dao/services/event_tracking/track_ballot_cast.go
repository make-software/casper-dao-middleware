package event_tracking

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackBallotCast struct {
	di.EntityManagerAware
	di.DAOContractsMetadataAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackBallotCast() *TrackBallotCast {
	return &TrackBallotCast{}
}

func (s *TrackBallotCast) Execute() error {
	ballotCast, err := events.ParseBallotCastEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var voter *types.Hash
	if ballotCast.Voter.AccountHash != nil {
		voter = ballotCast.Voter.AccountHash
	} else {
		voter = ballotCast.Voter.ContractPackageHash
	}

	staked := ballotCast.Stake.Into().Int64()

	var isInFavor bool
	if ballotCast.Choice == events.ChoiceInFavor {
		isInFavor = true
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	vote := entities.NewVote(
		*voter,
		deployProcessedEvent.DeployProcessed.DeployHash,
		ballotCast.VotingID,
		uint64(staked),
		isInFavor,
		deployProcessedEvent.DeployProcessed.Timestamp)
	if err := s.GetEntityManager().VoteRepository().Save(vote); err != nil {
		return err
	}

	changes := []entities.ReputationChange{
		// one event represent negative reputation leaving from "Reputation" contract
		entities.NewReputationChange(
			*voter,
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			&ballotCast.VotingID,
			-staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			deployProcessedEvent.DeployProcessed.Timestamp),
		// second event represent positive reputation coming to "Voting" contract
		entities.NewReputationChange(
			*voter,
			s.GetDAOContractsMetadata().SimpleVoterContractPackageHash,
			&ballotCast.VotingID,
			staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
