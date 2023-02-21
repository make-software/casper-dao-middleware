package event_tracking

import (
	"encoding/json"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackVotingCreated() *TrackVotingCreated {
	return &TrackVotingCreated{}
}

func (s *TrackVotingCreated) Execute() error {
	var voting entities.Voting
	var err error

	switch s.GetCESEvent().Name {
	case events.SimpleVotingCreatedEventName:
		voting, err = s.newVotingFromSimpleVotingCreated()
	case events.ReputationVotingCreatedEventName:
		voting, err = s.newVotingFromReputationVotingCreated()
	case events.RepoVotingCreated:
		voting, err = s.newVotingFromRepoVotingCreated()
	case events.SlashingVotingCreated:
		voting, err = s.newVotingFromSlashingVotingCreated()
	case events.KYCVotingCreated:
		voting, err = s.newVotingFromKYCVotingCreated()
	}
	if err != nil {
		return err
	}

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackVotingCreated) newVotingFromSimpleVotingCreated() (entities.Voting, error) {
	simpleVotingCreated, err := events.ParseSimpleVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
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
		return entities.Voting{}, err
	}

	return entities.NewVoting(
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
	), nil
}

func (s *TrackVotingCreated) newVotingFromReputationVotingCreated() (entities.Voting, error) {
	reputationVotingCreated, err := events.ParseReputationVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
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

	var account types.Hash
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
		return entities.Voting{}, err
	}

	return entities.NewVoting(
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
	), nil
}

func (s *TrackVotingCreated) newVotingFromRepoVotingCreated() (entities.Voting, error) {
	repoVotingCreatedEvent, err := events.ParseRepoVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
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

	var variableRepoToEdit types.Hash
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
		return entities.Voting{}, err
	}

	return entities.NewVoting(
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
	), nil
}

func (s *TrackVotingCreated) newVotingFromSlashingVotingCreated() (entities.Voting, error) {
	slashingVotingCreatedEvent, err := events.ParseSlashingVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
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

	var addressToSlash types.Hash
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
		return entities.Voting{}, err
	}

	return entities.NewVoting(
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
	), nil
}

func (s *TrackVotingCreated) newVotingFromKYCVotingCreated() (entities.Voting, error) {
	kycVotingCreated, err := events.ParseKYCVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
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

	var subjectAddress types.Hash
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
		return entities.Voting{}, err
	}

	return entities.NewVoting(
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
	), nil
}
