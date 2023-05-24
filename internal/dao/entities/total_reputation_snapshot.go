package entities

import (
	"time"

	"github.com/make-software/casper-go-sdk/casper"
)

type LiquidStakeReputation struct {
	LiquidAmount *uint64      `json:"liquid_amount" db:"liquid_amount"`
	StakedAmount *uint64      `json:"staked_amount" db:"staked_amount"`
	Address      *casper.Hash `json:"address" db:"address"`
}

type TotalReputationSnapshot struct {
	Address                casper.Hash            `json:"address" db:"address"`
	TotalLiquidReputation  uint64                 `json:"total_liquid_reputation" db:"total_liquid_reputation"`
	TotalStakedReputation  uint64                 `json:"total_staked_reputation" db:"total_staked_reputation"`
	VotingLostReputation   uint64                 `json:"voting_lost_reputation" db:"voting_lost_reputation"`
	VotingEarnedReputation uint64                 `json:"voting_earned_reputation" db:"voting_earned_reputation"`
	VotingID               *uint32                `json:"voting_id" db:"voting_id"`
	DeployHash             casper.Hash            `json:"deploy_hash" db:"deploy_hash"`
	Reason                 ReputationChangeReason `json:"reason" db:"reason"`
	Timestamp              time.Time              `json:"timestamp" db:"timestamp"`
}

func NewTotalReputationSnapshot(
	address casper.Hash,
	votingID *uint32,
	totalLiquidReputation, totalStakedReputation uint64,
	votingLostReputation, votingEarnedReputation uint64,
	deployHash casper.Hash,
	reason ReputationChangeReason,
	timestamp time.Time) TotalReputationSnapshot {
	return TotalReputationSnapshot{
		Address:                address,
		TotalLiquidReputation:  totalLiquidReputation,
		TotalStakedReputation:  totalStakedReputation,
		VotingLostReputation:   votingLostReputation,
		VotingEarnedReputation: votingEarnedReputation,
		VotingID:               votingID,
		DeployHash:             deployHash,
		Reason:                 reason,
		Timestamp:              timestamp,
	}
}
