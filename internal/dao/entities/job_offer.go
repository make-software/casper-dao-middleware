package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type AuctionType byte

const (
	AuctionTypeInternal AuctionType = iota + 1
	AuctionTypeExternal
)

type JobOffer struct {
	JobOfferID        uint32      `json:"job_offer_id" db:"job_offer_id"`
	JobPoster         types.Hash  `json:"job_poster" db:"job_poster"`
	DeployHash        types.Hash  `json:"deploy_hash" db:"deploy_hash"`
	MaxBudget         uint64      `json:"max_budget" db:"max_budget"`
	AuctionType       AuctionType `json:"auction_type" db:"auction_type"`
	ExpectedTimeFrame uint64      `json:"expected_time_frame"  db:"expected_time_frame"`
	Timestamp         time.Time   `json:"timestamp"  db:"timestamp"`
}

func NewJobOffer(
	jobOfferID uint32,
	deployHash, jobPoster types.Hash,
	maxBudget uint64,
	auctionType AuctionType,
	expectedTimeFrame uint64,
	timestamp time.Time) JobOffer {
	return JobOffer{
		JobOfferID:        jobOfferID,
		JobPoster:         jobPoster,
		DeployHash:        deployHash,
		MaxBudget:         maxBudget,
		AuctionType:       auctionType,
		ExpectedTimeFrame: expectedTimeFrame,
		Timestamp:         timestamp,
	}
}
