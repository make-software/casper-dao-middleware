package types

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Address struct {
	AccountHash         *casper.Hash
	ContractPackageHash *casper.Hash
}

func NewAddressFromCLValue(val casper.CLValue) (Address, error) {
	if val.Type != cltype.Key {
		return Address{}, errors.New("invalid CLTypeID")
	}

	if val.Key == nil {
		return Address{}, errors.New("nil Key in CLValue")
	}

	if val.Key.Account != nil && val.Key.Hash != nil {
		return Address{}, errors.New("expected one value in Key")
	}

	address := Address{
		ContractPackageHash: val.Key.Hash,
	}

	if val.Key.Account != nil {
		address.AccountHash = &val.Key.Account.Hash
	}
	return address, nil
}

func (a Address) ToHash() *casper.Hash {
	if a.AccountHash != nil {
		return a.AccountHash
	}
	return a.ContractPackageHash
}
