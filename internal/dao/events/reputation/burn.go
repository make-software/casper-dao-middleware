package reputation

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const BurnEventName = "Burn"

type Burn struct {
	Address types.Address
	Amount  casper_types.U512
}

func ParseBurn(event ces.Event) (Burn, error) {
	var burn Burn

	val, ok := event.Data["address"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return Burn{}, errors.New("invalid address value in event")
	}
	burn.Address = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return Burn{}, errors.New("invalid amount value in event")
	}
	burn.Amount = *val.U512

	return burn, nil
}
