package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/repo_voter"
)

type TrackRepoVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackRepoVotingCreated() *TrackRepoVotingCreated {
	return &TrackRepoVotingCreated{}
}

func (s *TrackRepoVotingCreated) Execute() error {
	repoVotingCreatedEvent, err := repo_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"variable_repo_to_edit": repoVotingCreatedEvent.VariableRepoToEdit.ToHash().ToHex(),
		"key":                   repoVotingCreatedEvent.Key,
		"value":                 string(repoVotingCreatedEvent.Value),
		"activation_time":       repoVotingCreatedEvent.ActivationTime,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(repoVotingCreatedEvent.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !repoVotingCreatedEvent.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(repoVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(repoVotingCreatedEvent.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*repoVotingCreatedEvent.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		repoVotingCreatedEvent.VotingID,
		entities.VotingTypeRepo,
		metadataJSON,
		repoVotingCreatedEvent.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		repoVotingCreatedEvent.ConfigFormalQuorum,
		repoVotingCreatedEvent.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		repoVotingCreatedEvent.ConfigTotalOnboarded.Value().Uint64(),
		repoVotingCreatedEvent.ConfigVotingClearnessDelta.Value().Uint64(),
		repoVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
