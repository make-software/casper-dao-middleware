package slashing_voter

import (
	base_events "casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/services/event_tracking/base"
)

type TrackBallotCast struct {
	base.TrackBallotCast
}

func NewTrackBallotCast() *TrackBallotCast {
	return &TrackBallotCast{}
}

func (s *TrackBallotCast) Execute() error {
	ballotCast, err := base_events.ParseBallotCastEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.SaveVote(ballotCast); err != nil {
		return err
	}

	if err := s.CollectReputationChanges(ballotCast, s.GetDAOContractsMetadata().SlashingVoterContractPackageHash); err != nil {
		return err
	}

	if err := s.AggregateReputationTotals(ballotCast); err != nil {
		return err
	}

	return nil
}
