package reputation

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const MintEventName = "Mint"

type Mint struct {
	Address types.Address
	Amount  clvalue.UInt512
}

func ParseMint(event ces.Event) (Mint, error) {
	var mint Mint

	val, ok := event.Data["address"]
	if !ok || val.Type != cltype.Key {
		return Mint{}, errors.New("invalid address value in event")
	}
	mint.Address = types.Address{
		ContractPackageHash: val.Key.Hash,
	}

	if val.Key.Account != nil {
		mint.Address.AccountHash = &val.Key.Account.Hash
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type != cltype.UInt512 {
		return Mint{}, errors.New("invalid amount value in event")
	}
	mint.Amount = *val.UI512

	return mint, nil
}
