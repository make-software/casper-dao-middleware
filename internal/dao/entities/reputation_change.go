package entities

import (
	"time"

	"github.com/make-software/casper-go-sdk/casper"
)

type ReputationChange struct {
	Address             casper.Hash                `json:"address" db:"address"`
	ContractPackageHash casper.ContractPackageHash `json:"contract_package_hash" db:"contract_package_hash"`
	VotingID            *uint32                    `json:"voting_id" db:"voting_id"`
	Amount              int64                      `json:"amount" db:"amount"`
	DeployHash          casper.Hash                `json:"deploy_hash" db:"deploy_hash"`
	Reason              ReputationChangeReason     `json:"reason" db:"reason"`
	Timestamp           time.Time                  `json:"timestamp" db:"timestamp"`
}

func NewReputationChange(
	address casper.Hash,
	contractPackageHash casper.ContractPackageHash,
	votingID *uint32,
	amount int64,
	deployHash casper.Hash,
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
