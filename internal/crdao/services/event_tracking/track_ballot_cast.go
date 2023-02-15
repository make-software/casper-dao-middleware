package event_tracking

import (
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackBallotCast struct {
	di.EntityManagerAware
	di.DAOContractsMetadataAware

	cesEvent        ces.Event
	deployProcessed casper.DeployProcessed
}

func NewTrackBallotCast() *TrackBallotCast {
	return &TrackBallotCast{}
}

func (s *TrackBallotCast) SetCESEvent(event ces.Event) {
	s.cesEvent = event
}

func (s *TrackBallotCast) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackBallotCast) Execute() error {
	ballotCast, err := events.ParseBallotCastEvent(s.cesEvent)
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

	vote := entities.NewVote(
		*voter,
		s.deployProcessed.DeployHash,
		ballotCast.VotingID,
		uint64(staked),
		isInFavor,
		s.deployProcessed.Timestamp)
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
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			s.deployProcessed.Timestamp),
		// second event represent positive reputation coming to "Voting" contract
		entities.NewReputationChange(
			*voter,
			s.GetDAOContractsMetadata().SimpleVoterContractPackageHash,
			&ballotCast.VotingID,
			staked,
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
