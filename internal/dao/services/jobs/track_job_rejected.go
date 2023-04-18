package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobRejected struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobRejected() *TrackJobRejected {
	return &TrackJobRejected{}
}

func (s *TrackJobRejected) Execute() error {
	jobRejected, err := bid_escrow.ParseJobRejectedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	job, err := s.GetEntityManager().JobRepository().GetByBidID(jobRejected.BidID)
	if err != nil {
		return err
	}

	caller, err := jobRejected.Caller.GetHashValue()
	if err != nil {
		return err
	}

	job.Caller = caller
	job.JobStatusID = entities.JobStatusIDRejected

	return s.GetEntityManager().JobRepository().Update(job)
}
