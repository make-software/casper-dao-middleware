package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/admin"
)

type TrackAdminVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackAdminVotingCreated() *TrackAdminVotingCreated {
	return &TrackAdminVotingCreated{}
}

func (s *TrackAdminVotingCreated) Execute() error {
	adminVotingCreatedEvent, err := admin.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"contract_to_update": adminVotingCreatedEvent.ContractToUpdate.ToHash(),
		"action":             adminVotingCreatedEvent.Action,
		"address":            adminVotingCreatedEvent.Address.ToHash(),
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(adminVotingCreatedEvent.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !adminVotingCreatedEvent.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(adminVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(adminVotingCreatedEvent.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*adminVotingCreatedEvent.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		adminVotingCreatedEvent.VotingID,
		entities.VotingTypeAdmin,
		metadataJSON,
		adminVotingCreatedEvent.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		adminVotingCreatedEvent.ConfigFormalQuorum,
		adminVotingCreatedEvent.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		adminVotingCreatedEvent.ConfigTotalOnboarded.Into().Uint64(),
		adminVotingCreatedEvent.ConfigVotingClearnessDelta.Into().Uint64(),
		adminVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
