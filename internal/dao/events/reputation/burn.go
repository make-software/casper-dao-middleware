package reputation

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const BurnEventName = "Burn"

type Burn struct {
	Address types.Address
	Amount  clvalue.UInt512
}

func ParseBurn(event ces.Event) (Burn, error) {
	var burn Burn

	val, ok := event.Data["address"]
	if !ok || val.Type != cltype.Key {
		return Burn{}, errors.New("invalid address value in event")
	}
	burn.Address = types.Address{
		ContractPackageHash: val.Key.Hash,
	}

	if val.Key.Account != nil {
		burn.Address.AccountHash = &val.Key.Account.Hash
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type != cltype.UInt512 {
		return Burn{}, errors.New("invalid amount value in event")
	}
	burn.Amount = *val.UI512

	return burn, nil
}
