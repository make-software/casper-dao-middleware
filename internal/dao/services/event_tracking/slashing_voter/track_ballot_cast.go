package slashing_voter

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/types"
)

type TrackBallotCast struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackBallotCast() *TrackBallotCast {
	return &TrackBallotCast{}
}

func (s *TrackBallotCast) Execute() error {
	ballotCast, err := slashing_voter.ParseBallotCastEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	staked := ballotCast.Stake.Into().Int64()

	var isInFavor bool
	if ballotCast.Choice == types.ChoiceInFavor {
		isInFavor = true
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	vote := entities.NewVote(
		*ballotCast.Voter.ToHash(),
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
			*ballotCast.Voter.ToHash(),
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			&ballotCast.VotingID,
			-staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonStaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
		// second event represent positive reputation coming to "Voting" contract
		entities.NewReputationChange(
			*ballotCast.Voter.ToHash(),
			s.GetDAOContractsMetadata().SlashingVoterContractPackageHash,
			&ballotCast.VotingID,
			staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonStaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
