package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
)

type TrackSlashingVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackSlashingVotingCreated() *TrackSlashingVotingCreated {
	return &TrackSlashingVotingCreated{}
}

func (s *TrackSlashingVotingCreated) Execute() error {
	slashingVotingCreatedEvent, err := slashing_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"address_to_slash": slashingVotingCreatedEvent.AddressToSlash.ToHash().ToHex(),
		"slash_ration":     slashingVotingCreatedEvent.SlashRation,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(slashingVotingCreatedEvent.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !slashingVotingCreatedEvent.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(slashingVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(slashingVotingCreatedEvent.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*slashingVotingCreatedEvent.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		slashingVotingCreatedEvent.VotingID,
		entities.VotingTypeSlashing,
		metadataJSON,
		slashingVotingCreatedEvent.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		slashingVotingCreatedEvent.ConfigFormalQuorum,
		slashingVotingCreatedEvent.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		slashingVotingCreatedEvent.ConfigTotalOnboarded.Into().Uint64(),
		slashingVotingCreatedEvent.ConfigVotingClearnessDelta.Into().Uint64(),
		slashingVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
