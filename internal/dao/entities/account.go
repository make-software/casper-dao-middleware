package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Account struct {
	Hash      types.Hash `json:"creator" db:"creator"`
	IsKyc     bool       `json:"is_kyc" db:"is_formal"`
	IsVA      bool       `json:"is_va" db:"is_va"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
}

func NewAccount(hash types.Hash, isKyc, isVA bool, timestamp time.Time) Account {
	return Account{
		Hash:      hash,
		IsKyc:     isKyc,
		IsVA:      isVA,
		Timestamp: timestamp,
	}
}
