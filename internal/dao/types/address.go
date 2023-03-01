package types

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
)

type Address struct {
	AccountHash         *casper_types.Hash
	ContractPackageHash *casper_types.Hash
}

func NewAddressFromCLValue(val casper_types.CLValue) (Address, error) {
	if val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return Address{}, errors.New("invalid CLTypeID")
	}

	if val.Key == nil {
		return Address{}, errors.New("nil Key in CLValue")
	}

	if val.Key.AccountHash != nil && val.Key.Hash != nil {
		return Address{}, errors.New("expected one value in Key")
	}

	return Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}, nil
}

func (a Address) ToHash() *casper_types.Hash {
	if a.AccountHash != nil {
		return a.AccountHash
	}
	return a.ContractPackageHash
}
