package bid

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackBidSubmitted struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackBidSubmitted() *TrackBidSubmitted {
	return &TrackBidSubmitted{}
}

func (s *TrackBidSubmitted) Execute() error {
	bidSubmitted, err := bid_escrow.ParseBidSubmittedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var reputationStake *uint64
	if bidSubmitted.ReputationStake != nil {
		stake := bidSubmitted.ReputationStake.Value().Uint64()
		reputationStake = &stake
	} else {
		// if the reputation stake is missing it means the bid contains stake in cspr which is possible only in External auction
		if err := s.GetEntityManager().JobOfferRepository().UpdateAuctionType(bidSubmitted.JobOfferID, entities.AuctionTypeIDExternal); err != nil {
			return err
		}
	}

	var csprStake *uint64
	if bidSubmitted.CSPRStake != nil {
		stake := bidSubmitted.CSPRStake.Value().Uint64()
		csprStake = &stake
	}

	bid := entities.NewBid(
		bidSubmitted.JobOfferID,
		bidSubmitted.BidID,
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		bidSubmitted.Worker,
		bidSubmitted.Onboard,
		bidSubmitted.ProposedTimeFrame,
		bidSubmitted.ProposedPayment.Value().Uint64(),
		false,
		reputationStake,
		csprStake,
		time.Now().UTC(),
	)

	return s.GetEntityManager().BidRepository().Save(&bid)
}
