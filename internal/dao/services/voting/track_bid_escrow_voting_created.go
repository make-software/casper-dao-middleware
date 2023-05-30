package voting

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/bid_escrow"
)

type TrackBidEscrowVotingCreated struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
}

func NewTrackBidEscrowVotingCreated() *TrackBidEscrowVotingCreated {
	return &TrackBidEscrowVotingCreated{}
}

func (s *TrackBidEscrowVotingCreated) Execute() error {
	bidEscrowVotingCreated, err := bid_escrow.ParseVotingCreatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"job_id":       bidEscrowVotingCreated.JobID,
		"bid_id":       bidEscrowVotingCreated.BidID,
		"job_offer_id": bidEscrowVotingCreated.JobOfferID,
		"worker":       bidEscrowVotingCreated.Worker.ToHex(),
		"job_poster":   bidEscrowVotingCreated.JobPoster.ToHex(),
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// starts the informal when the event was emitted
	informalVotingStartsAt := time.Now().UTC()
	informalVotingEndsAt := informalVotingStartsAt.Add(time.Millisecond * time.Duration(bidEscrowVotingCreated.ConfigInformalVotingTime))

	var formalVotingStartsAt, formalVotingEndsAt *time.Time

	// if the `config_double_time_between_votings` is false we can surely say when FormalVoting will start
	// as there is no need to have calculation of VotingEnded percentage based on `voting_clearness_delta`
	if !bidEscrowVotingCreated.ConfigDoubleTimeBetweenVotings {
		startsAt := informalVotingEndsAt.Add(time.Millisecond * time.Duration(bidEscrowVotingCreated.ConfigTimeBetweenInformalAndFormalVoting))
		formalVotingStartsAt = &startsAt

		endsAt := formalVotingStartsAt.Add(time.Millisecond * time.Duration(bidEscrowVotingCreated.ConfigFormalVotingTime))
		formalVotingEndsAt = &endsAt
	}

	voting := entities.NewVoting(
		*bidEscrowVotingCreated.Creator.ToHash(),
		s.GetDeployProcessedEvent().DeployProcessed.DeployHash,
		bidEscrowVotingCreated.VotingID,
		entities.VotingTypeBidEscrow,
		metadataJSON,
		bidEscrowVotingCreated.ConfigInformalQuorum,
		informalVotingStartsAt,
		informalVotingEndsAt,
		bidEscrowVotingCreated.ConfigFormalQuorum,
		bidEscrowVotingCreated.ConfigFormalVotingTime,
		formalVotingStartsAt, formalVotingEndsAt,
		bidEscrowVotingCreated.ConfigTotalOnboarded.Value().Uint64(),
		bidEscrowVotingCreated.ConfigVotingClearnessDelta.Value().Uint64(),
		bidEscrowVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
	)

	return s.GetEntityManager().VotingRepository().Save(&voting)
}
