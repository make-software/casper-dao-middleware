package va_nft

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const TransferEventName = "Transfer"

type TransferEvent struct {
	From    *types.Address
	To      *types.Address
	TokenID clvalue.UInt512
}

func ParseTransferEvent(event ces.Event) (TransferEvent, error) {
	var vaTransfer TransferEvent

	val, ok := event.Data["from"]
	if !ok {
		return TransferEvent{}, errors.New("invalid from value in event")
	}

	if val.Option != nil {
		if val.Option.Type.Inner != cltype.Key {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}

		if val.Option.Inner != nil {
			from, err := types.NewAddressFromCLValue(*val.Option.Inner)
			if err != nil {
				return TransferEvent{}, errors.New("invalid value inside option of `from` value")
			}
			vaTransfer.From = &from
		}
	}

	val, ok = event.Data["to"]
	if !ok {
		return TransferEvent{}, errors.New("invalid to value in event")
	}

	if val.Option != nil {
		if val.Option.Inner.Type != cltype.Key {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}

		to, err := types.NewAddressFromCLValue(*val.Option.Inner)
		if err != nil {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}
		vaTransfer.To = &to
	}

	val, ok = event.Data["token_id"]
	if !ok || val.Type != cltype.UInt512 {
		return TransferEvent{}, errors.New("invalid token_id value in event")
	}
	vaTransfer.TokenID = *val.UI512

	return vaTransfer, nil
}
