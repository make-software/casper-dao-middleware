package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type AggregatedReputationChange struct {
	Address        types.Hash `json:"address" db:"address"`
	EarnedAmount   uint64     `json:"earned_amount" db:"earned_amount"`
	LostAmount     uint64     `json:"lost_amount" db:"lost_amount"`
	StakedAmount   uint64     `json:"staked_amount" db:"staked_amount"`
	ReleasedAmount uint64     `json:"released_amount" db:"released_amount"`
	VotingID       *uint32    `json:"voting_id" db:"voting_id"`
	Timestamp      time.Time  `json:"timestamp" db:"timestamp"`
}

func NewAggregatedReputationChange(
	address types.Hash,
	earnedAmount, lostAmount uint64,
	stakedAmount, releasedAmount uint64,
	votingID *uint32,
	timestamp time.Time,
) AggregatedReputationChange {
	return AggregatedReputationChange{
		Address:        address,
		EarnedAmount:   earnedAmount,
		LostAmount:     lostAmount,
		StakedAmount:   stakedAmount,
		ReleasedAmount: releasedAmount,
		VotingID:       votingID,
		Timestamp:      timestamp,
	}
}
