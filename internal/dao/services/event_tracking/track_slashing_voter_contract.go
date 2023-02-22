package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackSlashingVoterContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackSlashingVoterContract() *TrackSlashingVoterContract {
	return &TrackSlashingVoterContract{}
}

func (s *TrackSlashingVoterContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case slashing_voter.VotingCreatedEventName:
		return s.trackVotingCreated()
	case slashing_voter.BallotCastEventName:
		return s.trackBallotCast()
	}

	return nil
}

func (s *TrackSlashingVoterContract) trackVotingCreated() error {
	slashingVotingCreatedEvent, err := slashing_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var creator casper_types.Hash
	if slashingVotingCreatedEvent.Creator.AccountHash != nil {
		creator = *slashingVotingCreatedEvent.Creator.AccountHash
	} else {
		creator = *slashingVotingCreatedEvent.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum = slashingVotingCreatedEvent.ConfigInformalQuorum
	var votingTime = slashingVotingCreatedEvent.ConfigInformalVotingTime

	if slashingVotingCreatedEvent.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = slashingVotingCreatedEvent.ConfigFormalQuorum
		votingTime = slashingVotingCreatedEvent.ConfigFormalVotingTime
	}

	var addressToSlash casper_types.Hash
	if slashingVotingCreatedEvent.AddressToSlash.AccountHash != nil {
		addressToSlash = *slashingVotingCreatedEvent.AddressToSlash.AccountHash
	} else {
		addressToSlash = *slashingVotingCreatedEvent.AddressToSlash.ContractPackageHash
	}

	metadata := map[string]interface{}{
		"address_to_slash": addressToSlash.ToHex(),
		"slash_ration":     slashingVotingCreatedEvent.SlashRation,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	voting := entities.NewVoting(
		creator,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		slashingVotingCreatedEvent.VotingID,
		votingQuorum,
		votingTime,
		entities.VotingTypeSlashing,
		metadataJSON,
		isFormal,
		slashingVotingCreatedEvent.ConfigDoubleTimeBetweenVotings,
		slashingVotingCreatedEvent.ConfigTotalOnboarded.Into().Uint64(),
		slashingVotingCreatedEvent.ConfigVotingClearnessDelta.Into().Uint64(),
		slashingVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
		s.GetDeployProcessedEvent().DeployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackSlashingVoterContract) trackBallotCast() error {
	ballotCast, err := slashing_voter.ParseBallotCastEvent(s.GetCESEvent())
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
