package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation_voter"
	"casper-dao-middleware/internal/dao/events/simple_voter"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackSimpleVoterContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackSimpleVoterContract() *TrackSimpleVoterContract {
	return &TrackSimpleVoterContract{}
}

func (s *TrackSimpleVoterContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case simple_voter.VotingCreatedEventName:
		return s.trackVotingCreated()
	case simple_voter.BallotCastEventName:
		return s.trackBallotCast()
	}

	return nil
}

func (s *TrackSimpleVoterContract) trackVotingCreated() error {
	simpleVotingCreated, err := simple_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var creator casper_types.Hash
	if simpleVotingCreated.Creator.AccountHash != nil {
		creator = *simpleVotingCreated.Creator.AccountHash
	} else {
		creator = *simpleVotingCreated.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum = simpleVotingCreated.ConfigInformalQuorum
	var votingTime = simpleVotingCreated.ConfigInformalVotingTime

	if simpleVotingCreated.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = simpleVotingCreated.ConfigFormalQuorum
		votingTime = simpleVotingCreated.ConfigFormalVotingTime
	}

	metadata := map[string]interface{}{
		"document_hash": simpleVotingCreated.DocumentHash,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	voting := entities.NewVoting(
		creator,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		simpleVotingCreated.VotingID,
		votingQuorum,
		votingTime,
		entities.VotingTypeSimple,
		metadataJSON,
		isFormal,
		simpleVotingCreated.ConfigDoubleTimeBetweenVotings,
		simpleVotingCreated.ConfigTotalOnboarded.Into().Uint64(),
		simpleVotingCreated.ConfigVotingClearnessDelta.Into().Uint64(),
		simpleVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
		s.GetDeployProcessedEvent().DeployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackSimpleVoterContract) trackBallotCast() error {
	ballotCast, err := reputation_voter.ParseBallotCastEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var voter *casper_types.Hash
	if ballotCast.Voter.AccountHash != nil {
		voter = ballotCast.Voter.AccountHash
	} else {
		voter = ballotCast.Voter.ContractPackageHash
	}

	staked := ballotCast.Stake.Into().Int64()

	var isInFavor bool
	if ballotCast.Choice == types.ChoiceInFavor {
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
