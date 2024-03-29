package jobs

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobCreated() *TrackJobCreated {
	return &TrackJobCreated{}
}

func (s *TrackJobCreated) Execute() error {
	jobCreated, err := bid_escrow.ParseJobCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.GetEntityManager().BidRepository().UpdateIsPickedBy(jobCreated.BidID, true); err != nil {
		return err
	}

	job := entities.NewJob(
		jobCreated.JobID,
		jobCreated.BidID,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		jobCreated.JobPoster,
		jobCreated.Worker,
		jobCreated.FinishTime,
		entities.JobStatusIDCreated,
		nil,
		nil,
		time.Now().UTC(),
	)

	return s.GetEntityManager().JobRepository().Save(&job)
}
