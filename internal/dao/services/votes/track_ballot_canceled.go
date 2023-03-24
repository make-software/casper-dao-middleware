package votes

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/base"
)

type TrackCanceledVote struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackCanceledVote() *TrackCanceledVote {
	return &TrackCanceledVote{}
}

func (s *TrackCanceledVote) Execute() error {
	ballotCanceled, err := base.ParseBallotCanceledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	return s.GetEntityManager().VoteRepository().UpdateIsCanceled(ballotCanceled.VotingID, *ballotCanceled.Voter.ToHash(), true)
}
