package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Bid struct {
	JobOfferID        uint32     `json:"job_offer_id" db:"job_offer_id"`
	BidID             uint32     `json:"bid_id" db:"bid_id"`
	Worker            types.Hash `json:"worker" db:"worker"`
	DeployHash        types.Hash `json:"deploy_hash" db:"deploy_hash"`
	Onboard           bool       `json:"onboard" db:"onboard"`
	ProposedTimeFrame uint64     `json:"proposed_time_frame"  db:"proposed_time_frame"`
	ProposedPayment   uint64     `json:"proposed_payment"  db:"proposed_payment"`
	PickedByJobPoster bool       `json:"picked_by_job_poster" db:"picked_by_job_poster"`
	ReputationStake   *uint64    `json:"reputation_stake"  db:"reputation_stake"`
	CSPRStake         *uint64    `json:"cspr_stake"  db:"cspr_stake"`
	Timestamp         time.Time  `json:"timestamp"  db:"timestamp"`
}

func NewBid(
	jobOfferID uint32,
	bidID uint32,
	deployHash, worker types.Hash,
	onboard bool,
	proposedTimeFrame uint64,
	proposedPayment uint64,
	pickedByJobPoster bool,
	reputationStake, csprStake *uint64,
	timestamp time.Time) Bid {
	return Bid{
		JobOfferID:        jobOfferID,
		BidID:             bidID,
		Worker:            worker,
		DeployHash:        deployHash,
		Onboard:           onboard,
		ProposedTimeFrame: proposedTimeFrame,
		PickedByJobPoster: pickedByJobPoster,
		ProposedPayment:   proposedPayment,
		ReputationStake:   reputationStake,
		CSPRStake:         csprStake,
		Timestamp:         timestamp,
	}
}
