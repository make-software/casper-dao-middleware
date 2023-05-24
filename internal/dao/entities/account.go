package entities

import (
	"time"

	"github.com/make-software/casper-go-sdk/casper"
)

type Account struct {
	Hash      casper.Hash `json:"hash" db:"hash"`
	IsKyc     bool        `json:"is_kyc" db:"is_kyc"`
	IsVA      bool        `json:"is_va" db:"is_va"`
	Timestamp time.Time   `json:"timestamp" db:"timestamp"`
}

func NewAccount(hash casper.Hash, isKyc, isVA bool, timestamp time.Time) Account {
	return Account{
		Hash:      hash,
		IsKyc:     isKyc,
		IsVA:      isVA,
		Timestamp: timestamp,
	}
}
