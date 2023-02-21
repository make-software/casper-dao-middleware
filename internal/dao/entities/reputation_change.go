package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type ReputationChange struct {
	Address             types.Hash             `json:"address" db:"address"`
	ContractPackageHash types.Hash             `json:"contract_package_hash" db:"contract_package_hash"`
	VotingID            *uint32                `json:"voting_id" db:"voting_id"`
	Amount              int64                  `json:"amount" db:"amount"`
	DeployHash          types.Hash             `json:"deploy_hash" db:"deploy_hash"`
	Reason              ReputationChangeReason `json:"reason" db:"reason"`
	Timestamp           time.Time              `json:"timestamp" db:"timestamp"`
}

func NewReputationChange(
	address, contractPackageHash types.Hash,
	votingID *uint32,
	amount int64,
	deployHash types.Hash,
	reason ReputationChangeReason,
	timestamp time.Time) ReputationChange {
	return ReputationChange{
		Address:             address,
		ContractPackageHash: contractPackageHash,
		VotingID:            votingID,
		Amount:              amount,
		DeployHash:          deployHash,
		Reason:              reason,
		Timestamp:           timestamp,
	}
}
