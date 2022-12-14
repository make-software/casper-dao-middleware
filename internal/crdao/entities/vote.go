package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Vote struct {
	Address    types.Hash `json:"address" db:"address"`
	VotingID   uint32     `json:"voting_id" db:"voting_id"`
	Amount     uint64     `json:"amount" db:"amount"`
	IsInFavor  bool       `json:"is_in_favour" db:"is_in_favour"`
	DeployHash types.Hash `json:"deploy_hash" db:"deploy_hash"`
	Timestamp  time.Time  `json:"timestamp" db:"timestamp"`
}

func NewVote(address, deployHash types.Hash, votingID uint32, staked uint64, IsInFavor bool, timestamp time.Time) *Vote {
	return &Vote{
		Address:    address,
		VotingID:   votingID,
		Amount:     staked,
		DeployHash: deployHash,
		IsInFavor:  IsInFavor,
		Timestamp:  timestamp,
	}
}
