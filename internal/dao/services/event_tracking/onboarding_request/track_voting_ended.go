package onboarding_request

import (
	base_events "casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/services/event_tracking/base"
)

type TrackVotingEnded struct {
	base.TrackVotingEnded
}

func NewTrackVotingEnded() *TrackVotingEnded {
	return &TrackVotingEnded{}
}

func (s *TrackVotingEnded) Execute() error {
	votingEnded, err := base_events.ParseVotingEndedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.UpdateVotingState(votingEnded); err != nil {
		return err
	}

	if err := s.CollectReputationChanges(votingEnded, s.GetDAOContractsMetadata().OnboardingRequestContractPackageHash); err != nil {
		return err
	}

	if err := s.AggregateReputationTotals(votingEnded); err != nil {
		return err
	}

	return nil
}
