package job_offer

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackJobOfferCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackJobOfferCreated() *TrackJobOfferCreated {
	return &TrackJobOfferCreated{}
}

func (s *TrackJobOfferCreated) Execute() error {
	jobOfferCreated, err := bid_escrow.ParseJobOfferCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	jobPoster, err := jobOfferCreated.JobPoster.GetHashValue()
	if err != nil {
		return err
	}

	jobOffer := entities.NewJobOffer(
		jobOfferCreated.JobOfferID,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		*jobPoster,
		jobOfferCreated.MaxBudget.Into().Uint64(),
		entities.AuctionTypeIDInternal,
		jobOfferCreated.ExpectedTimeFrame,
		time.Now().UTC(),
	)

	return s.GetEntityManager().JobOfferRepository().Save(&jobOffer)
}
