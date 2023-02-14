package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const BurnEventName = "Burn"

type Burn struct {
	Address dao_types.Address
	Amount  types.U512
}

func ParseBurn(event ces.Event) (Burn, error) {
	var burn Burn

	val, ok := event.Data["address"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return Burn{}, errors.New("invalid address value in event")
	}
	burn.Address = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return Burn{}, errors.New("invalid amount value in event")
	}
	burn.Amount = *val.U512

	return burn, nil
}
