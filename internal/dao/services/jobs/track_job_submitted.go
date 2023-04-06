package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobSubmitted struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobSubmitted() *TrackJobSubmitted {
	return &TrackJobSubmitted{}
}

func (s *TrackJobSubmitted) Execute() error {
	jobSubmitted, err := bid_escrow.ParseJobSubmittedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	job, err := s.GetEntityManager().JobRepository().GetByBidID(jobSubmitted.BidID)
	if err != nil {
		return err
	}

	job.Result = &jobSubmitted.Result
	job.JobStatus = entities.JobStatusSubmitted

	return s.GetEntityManager().JobRepository().Update(job)
}
