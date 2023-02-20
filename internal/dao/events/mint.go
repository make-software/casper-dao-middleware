package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/dao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const MintEventName = "Mint"

type Mint struct {
	Address dao_types.Address
	Amount  types.U512
}

func ParseMint(event ces.Event) (Mint, error) {
	var mint Mint

	val, ok := event.Data["address"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return Mint{}, errors.New("invalid address value in event")
	}
	mint.Address = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return Mint{}, errors.New("invalid amount value in event")
	}
	mint.Amount = *val.U512

	return mint, nil
}
