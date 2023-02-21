package reputation

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const MintEventName = "Mint"

type Mint struct {
	Address types.Address
	Amount  casper_types.U512
}

func ParseMint(event ces.Event) (Mint, error) {
	var mint Mint

	val, ok := event.Data["address"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return Mint{}, errors.New("invalid address value in event")
	}
	mint.Address = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return Mint{}, errors.New("invalid amount value in event")
	}
	mint.Amount = *val.U512

	return mint, nil
}
