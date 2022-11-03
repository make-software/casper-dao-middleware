package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/dao_event_parser/types"
	"casper-dao-middleware/pkg/casper/types"
)

const BurnEventName = "Burn"

type Burn struct {
	Address dao_types.Address
	Amount  types.U256
}

func ParseBurnEvent(bytes []byte) (Burn, error) {
	key, reminder, err := types.ParseKeyFromBytes(bytes)
	if err != nil {
		return Burn{}, err
	}

	event := Burn{}
	if key.AccountHash == nil && key.Hash == nil {
		return Burn{}, errors.New("expected Address in Burn event")
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
		return Burn{}, err
	}

	return event, nil
}
