package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobDone struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobDone() *TrackJobDone {
	return &TrackJobDone{}
}

func (s *TrackJobDone) Execute() error {
	jobDone, err := bid_escrow.ParseJobDoneEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	job, err := s.GetEntityManager().JobRepository().GetByBidID(jobDone.BidID)
	if err != nil {
		return err
	}

	job.Caller = &jobDone.Caller
	job.JobStatusID = entities.JobStatusIDDone

	return s.GetEntityManager().JobRepository().Update(job)
}
