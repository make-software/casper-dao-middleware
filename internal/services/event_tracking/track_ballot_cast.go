package event_tracking

import (
	"casper-dao-middleware/internal/dao_event_parser/events"
	"casper-dao-middleware/internal/di"
	"casper-dao-middleware/internal/entities"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackBallotCast struct {
	di.EntityManagerAware
	di.DAOContractPackageHashesAware

	deployProcessed casper.DeployProcessed
	eventBody       []byte
}

func NewTrackBallotCast() *TrackBallotCast {
	return &TrackBallotCast{}
}

func (s *TrackBallotCast) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackBallotCast) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackBallotCast) Execute() error {
	ballotCast, err := events.ParseBallotCastEvent(s.eventBody)
	if err != nil {
		return err
	}

	var address *types.Hash
	if ballotCast.Address.AccountHash != nil {
		address = ballotCast.Address.AccountHash
	} else {
		address = ballotCast.Address.ContractPackageHash
	}

	staked := (*ballotCast.Stake).Int64()
	votingID := uint32((*ballotCast.VotingID).Uint64())

	var isInFavor bool
	if ballotCast.Choice == events.ChoiceInFavor {
		isInFavor = true
	}

	vote := entities.NewVote(
		*address,
		s.deployProcessed.DeployHash,
		votingID,
		uint64(staked),
		isInFavor,
		s.deployProcessed.Timestamp)
	if err := s.GetEntityManager().VoteRepository().Save(vote); err != nil {
		return err
	}

	changes := []entities.ReputationChange{
		// one event represent negative reputation leaving from "Reputation" contract
		entities.NewReputationChange(
			*address,
			s.GetDAOContractPackageHashes().ReputationContractPackageHash,
			&votingID,
			-staked,
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			s.deployProcessed.Timestamp),
		// second event represent positive reputation coming to "Voting" contract
		entities.NewReputationChange(
			*address,
			s.GetDAOContractPackageHashes().VoterContractPackageHash,
			&votingID,
			staked,
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVote,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
