package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/simple_voter"
)

type TrackSimpleVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackSimpleVotingCreated() *TrackSimpleVotingCreated {
	return &TrackSimpleVotingCreated{}
}

func (s *TrackSimpleVotingCreated) Execute() error {
	simpleVotingCreated, err := simple_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"document_hash": simpleVotingCreated.DocumentHash,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(simpleVotingCreated.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !simpleVotingCreated.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(simpleVotingCreated.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(simpleVotingCreated.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*simpleVotingCreated.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		simpleVotingCreated.VotingID,
		entities.VotingTypeSimple,
		metadataJSON,
		simpleVotingCreated.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		simpleVotingCreated.ConfigFormalQuorum,
		simpleVotingCreated.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		simpleVotingCreated.ConfigTotalOnboarded.Value().Uint64(),
		simpleVotingCreated.ConfigVotingClearnessDelta.Value().Uint64(),
		simpleVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
