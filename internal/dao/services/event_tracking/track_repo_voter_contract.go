package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/repo_voter"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackRepoVoterContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackRepoVoterContract() *TrackRepoVoterContract {
	return &TrackRepoVoterContract{}
}

func (s *TrackRepoVoterContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case repo_voter.VotingCreatedEventName:
		return s.trackVotingCreated()
	case repo_voter.BallotCastEventName:
	}

	return nil
}

func (s *TrackRepoVoterContract) trackVotingCreated() error {
	repoVotingCreatedEvent, err := repo_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var creator casper_types.Hash
	if repoVotingCreatedEvent.Creator.AccountHash != nil {
		creator = *repoVotingCreatedEvent.Creator.AccountHash
	} else {
		creator = *repoVotingCreatedEvent.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum = repoVotingCreatedEvent.ConfigInformalQuorum
	var votingTime = repoVotingCreatedEvent.ConfigInformalVotingTime

	if repoVotingCreatedEvent.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = repoVotingCreatedEvent.ConfigFormalQuorum
		votingTime = repoVotingCreatedEvent.ConfigFormalVotingTime
	}

	var variableRepoToEdit casper_types.Hash
	if repoVotingCreatedEvent.VariableRepoToEdit.AccountHash != nil {
		variableRepoToEdit = *repoVotingCreatedEvent.VariableRepoToEdit.AccountHash
	} else {
		variableRepoToEdit = *repoVotingCreatedEvent.VariableRepoToEdit.ContractPackageHash
	}

	metadata := map[string]interface{}{
		"variable_repo_to_edit": variableRepoToEdit.ToHex(),
		"key":                   repoVotingCreatedEvent.Key,
		"value":                 repoVotingCreatedEvent.Value,
		"activation_time":       repoVotingCreatedEvent.ActivationTime,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	voting := entities.NewVoting(
		creator,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		repoVotingCreatedEvent.VotingID,
		votingQuorum,
		votingTime,
		entities.VotingTypeRepo,
		metadataJSON,
		isFormal,
		repoVotingCreatedEvent.ConfigDoubleTimeBetweenVotings,
		repoVotingCreatedEvent.ConfigTotalOnboarded.Into().Uint64(),
		repoVotingCreatedEvent.ConfigVotingClearnessDelta.Into().Uint64(),
		repoVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
		s.GetDeployProcessedEvent().DeployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackRepoVoterContract) trackBallotCast() error {
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
