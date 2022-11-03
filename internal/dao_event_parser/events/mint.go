package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/dao_event_parser/types"
	"casper-dao-middleware/pkg/casper/types"
)

const MintEventName = "Mint"

type Mint struct {
	Address dao_types.Address
	Amount  types.U256
}

func ParseMintEvent(bytes []byte) (Mint, error) {
	key, reminder, err := types.ParseKeyFromBytes(bytes)
	if err != nil {
		return Mint{}, err
	}

	event := Mint{}
	if key.AccountHash == nil && key.Hash == nil {
		return Mint{}, errors.New("expected Address in Mint event")
	}

	var address dao_types.Address
	if key.AccountHash != nil {
		address.AccountHash = key.AccountHash
	} else {
		address.ContractPackageHash = key.Hash
	}

	event.Address = address

	event.Amount, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return Mint{}, err
	}

	return event, nil
}
