package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation_voter"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackReputationVoterContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackReputationVoterContract() *TrackReputationVoterContract {
	return &TrackReputationVoterContract{}
}

func (s *TrackReputationVoterContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case reputation_voter.VotingCreatedEventName:
		return s.trackVotingCreated()
	case reputation_voter.BallotCastEventName:
		return s.trackBallotCast()
	}

	return nil
}

func (s *TrackReputationVoterContract) trackVotingCreated() error {
	reputationVotingCreated, err := reputation_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var creator casper_types.Hash
	if reputationVotingCreated.Creator.AccountHash != nil {
		creator = *reputationVotingCreated.Creator.AccountHash
	} else {
		creator = *reputationVotingCreated.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum = reputationVotingCreated.ConfigInformalQuorum
	var votingTime = reputationVotingCreated.ConfigInformalVotingTime

	if reputationVotingCreated.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = reputationVotingCreated.ConfigFormalQuorum
		votingTime = reputationVotingCreated.ConfigFormalVotingTime
	}

	var account casper_types.Hash
	if reputationVotingCreated.Account.AccountHash != nil {
		account = *reputationVotingCreated.Account.AccountHash
	} else {
		account = *reputationVotingCreated.Account.ContractPackageHash
	}

	metadata := map[string]interface{}{
		"document_hash": reputationVotingCreated.DocumentHash,
		"account":       account.ToHex(),
		"action":        reputationVotingCreated.Action,
		"amount":        reputationVotingCreated.Amount.Into().Uint64(),
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	voting := entities.NewVoting(
		creator,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		reputationVotingCreated.VotingID,
		votingQuorum,
		votingTime,
		entities.VotingTypeReputation,
		metadataJSON,
		isFormal,
		reputationVotingCreated.ConfigDoubleTimeBetweenVotings,
		reputationVotingCreated.ConfigTotalOnboarded.Into().Uint64(),
		reputationVotingCreated.ConfigVotingClearnessDelta.Into().Uint64(),
		reputationVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
		s.GetDeployProcessedEvent().DeployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackReputationVoterContract) trackBallotCast() error {
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
