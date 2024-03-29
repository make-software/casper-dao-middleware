package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation_voter"
)

type TrackReputationVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackReputationVotingCreated() *TrackReputationVotingCreated {
	return &TrackReputationVotingCreated{}
}

func (s *TrackReputationVotingCreated) Execute() error {
	reputationVotingCreated, err := reputation_voter.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"document_hash": reputationVotingCreated.DocumentHash,
		"account":       reputationVotingCreated.Account.ToHash().ToHex(),
		"action":        reputationVotingCreated.Action,
		"amount":        reputationVotingCreated.Amount.Value().Uint64(),
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(reputationVotingCreated.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !reputationVotingCreated.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(reputationVotingCreated.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(reputationVotingCreated.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*reputationVotingCreated.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		reputationVotingCreated.VotingID,
		entities.VotingTypeReputation,
		metadataJSON,
		reputationVotingCreated.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		reputationVotingCreated.ConfigFormalQuorum,
		reputationVotingCreated.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		reputationVotingCreated.ConfigTotalOnboarded.Value().Uint64(),
		reputationVotingCreated.ConfigVotingClearnessDelta.Value().Uint64(),
		reputationVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
