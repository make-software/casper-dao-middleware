package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/kyc_voter"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackKycVoterContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackKycVoterContract() *TrackKycVoterContract {
	return &TrackKycVoterContract{}
}

func (s *TrackKycVoterContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case kyc_voter.VotingCreatedEventName:
		return s.trackVotingCreated()
	case kyc_voter.BallotCastEventName:
		return s.trackBallotCast()
	}

	return nil
}

func (s *TrackKycVoterContract) trackVotingCreated() error {
	kycVotingCreated, err := kyc_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var creator casper_types.Hash
	if kycVotingCreated.Creator.AccountHash != nil {
		creator = *kycVotingCreated.Creator.AccountHash
	} else {
		creator = *kycVotingCreated.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum = kycVotingCreated.ConfigInformalQuorum
	var votingTime = kycVotingCreated.ConfigInformalVotingTime

	if kycVotingCreated.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = kycVotingCreated.ConfigFormalQuorum
		votingTime = kycVotingCreated.ConfigFormalVotingTime
	}

	var subjectAddress casper_types.Hash
	if kycVotingCreated.SubjectAddress.AccountHash != nil {
		subjectAddress = *kycVotingCreated.SubjectAddress.AccountHash
	} else {
		subjectAddress = *kycVotingCreated.SubjectAddress.ContractPackageHash
	}

	metadata := map[string]interface{}{
		"subject_address": subjectAddress.ToHex(),
		"document_hash":   kycVotingCreated.DocumentHash,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	voting := entities.NewVoting(
		creator,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		kycVotingCreated.VotingID,
		votingQuorum,
		votingTime,
		entities.VotingTypeKYC,
		metadataJSON,
		isFormal,
		kycVotingCreated.ConfigDoubleTimeBetweenVotings,
		kycVotingCreated.ConfigTotalOnboarded.Into().Uint64(),
		kycVotingCreated.ConfigVotingClearnessDelta.Into().Uint64(),
		kycVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
		s.GetDeployProcessedEvent().DeployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackKycVoterContract) trackBallotCast() error {
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
