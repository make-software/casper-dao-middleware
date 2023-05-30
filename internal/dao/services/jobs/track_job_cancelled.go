package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobCancelled struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobCancelled() *TrackJobCancelled {
	return &TrackJobCancelled{}
}

func (s *TrackJobCancelled) Execute() error {
	jobCancelled, err := bid_escrow.ParseJobCancelledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	job, err := s.GetEntityManager().JobRepository().GetByBidID(jobCancelled.BidID)
	if err != nil {
		return err
	}

	job.Caller = &jobCancelled.Caller
	job.JobStatusID = entities.JobStatusIDCancelled

	return s.GetEntityManager().JobRepository().Update(job)
}
