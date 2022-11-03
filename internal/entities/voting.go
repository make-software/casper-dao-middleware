package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Voting struct {
	Creator          types.Hash `json:"creator" db:"creator"`
	DeployHash       types.Hash `json:"deploy_hash" db:"deploy_hash"`
	VotingID         uint32     `json:"voting_id" db:"voting_id"`
	InformalVotingID *uint32    `json:"informal_voting_id" db:"informal_voting_id"`
	IsFormal         bool       `json:"is_formal" db:"is_formal"`
	HasEnded         bool       `json:"has_ended" db:"has_ended"`
	VotingQuorum     uint64     `json:"voting_quorum" db:"voting_quorum"`
	VotingTime       uint64     `json:"voting_time" db:"voting_time"`
	Timestamp        time.Time  `json:"timestamp" db:"timestamp"`
}

func NewVoting(
	creator, deployHash types.Hash,
	votingID uint32,
	informalVotingID *uint32,
	votingTime, votingQuorum uint64,
	isFormal bool,
	timestamp time.Time,
) *Voting {
	return &Voting{
		Creator:          creator,
		DeployHash:       deployHash,
		VotingID:         votingID,
		InformalVotingID: informalVotingID,
		IsFormal:         isFormal,
		VotingQuorum:     votingQuorum,
		VotingTime:       votingTime,
		Timestamp:        timestamp,
		HasEnded:         false,
	}
}
