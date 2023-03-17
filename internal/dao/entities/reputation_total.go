package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type LiquidStakeReputation struct {
	LiquidAmount *uint64     `json:"liquid_amount" db:"liquid_amount"`
	StakedAmount *uint64     `json:"staked_amount" db:"staked_amount"`
	Address      *types.Hash `json:"address" db:"address"`
}

type ReputationTotal struct {
	Address                types.Hash             `json:"address" db:"address"`
	TotalLiquidReputation  uint64                 `json:"total_liquid_reputation" db:"total_liquid_reputation"`
	TotalStakedReputation  uint64                 `json:"total_staked_reputation" db:"total_staked_reputation"`
	VotingLostReputation   uint64                 `json:"voting_lost_reputation" db:"voting_lost_reputation"`
	VotingEarnedReputation uint64                 `json:"voting_earned_reputation" db:"voting_earned_reputation"`
	VotingID               *uint32                `json:"voting_id" db:"voting_id"`
	DeployHash             types.Hash             `json:"deploy_hash" db:"deploy_hash"`
	Reason                 ReputationChangeReason `json:"reason" db:"reason"`
	Timestamp              time.Time              `json:"timestamp" db:"timestamp"`
}

func NewReputationTotal(
	address types.Hash,
	votingID *uint32,
	totalLiquidReputation, totalStakedReputation uint64,
	votingLostReputation, votingEarnedReputation uint64,
	deployHash types.Hash,
	reason ReputationChangeReason,
	timestamp time.Time) ReputationTotal {
	return ReputationTotal{
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
