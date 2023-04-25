package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/onboarding_request"
)

type TrackOnboardingVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackOnboardingVotingCreated() *TrackOnboardingVotingCreated {
	return &TrackOnboardingVotingCreated{}
}

func (s *TrackOnboardingVotingCreated) Execute() error {
	onboardingRequestVotingCreatedEvent, err := onboarding_request.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"reason":       onboardingRequestVotingCreatedEvent.Reason,
		"cspr_deposit": onboardingRequestVotingCreatedEvent.CsprDeposit.Into().String(),
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(onboardingRequestVotingCreatedEvent.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !onboardingRequestVotingCreatedEvent.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(onboardingRequestVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(onboardingRequestVotingCreatedEvent.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*onboardingRequestVotingCreatedEvent.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		onboardingRequestVotingCreatedEvent.VotingID,
		entities.VotingTypeOnboarding,
		metadataJSON,
		onboardingRequestVotingCreatedEvent.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		onboardingRequestVotingCreatedEvent.ConfigFormalQuorum,
		onboardingRequestVotingCreatedEvent.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		onboardingRequestVotingCreatedEvent.ConfigTotalOnboarded.Into().Uint64(),
		onboardingRequestVotingCreatedEvent.ConfigVotingClearnessDelta.Into().Uint64(),
		onboardingRequestVotingCreatedEvent.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
