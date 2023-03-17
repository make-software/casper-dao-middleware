package base

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/base"
)

type TrackBallotCanceled struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackBallotCanceled() *TrackBallotCanceled {
	return &TrackBallotCanceled{}
}

func (s *TrackBallotCanceled) Execute() error {
	ballotCanceled, err := base.ParseBallotCanceledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	return s.GetEntityManager().VoteRepository().UpdateIsCanceled(ballotCanceled.VotingID, *ballotCanceled.Voter.ToHash(), true)
}
