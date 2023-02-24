package kyc_nft

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const TransferEventName = "Transfer"

type TransferEvent struct {
	From    *types.Address
	To      *types.Address
	TokenID casper_types.U512
}

func ParseTransferEvent(event ces.Event) (TransferEvent, error) {
	var kycTransfer TransferEvent

	val, ok := event.Data["from"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return TransferEvent{}, errors.New("invalid from value in event")
	}

	if val.Option != nil {
		if val.Option.Type.CLTypeID != casper_types.CLTypeIDKey {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}

		from, err := types.NewAddressFromCLValue(*val.Option)
		if err != nil {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}
		kycTransfer.From = &from
	}

	val, ok = event.Data["to"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return TransferEvent{}, errors.New("invalid to value in event")
	}

	if val.Option != nil {
		if val.Option.Type.CLTypeID != casper_types.CLTypeIDKey {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}

		to, err := types.NewAddressFromCLValue(*val.Option)
		if err != nil {
			return TransferEvent{}, errors.New("invalid value inside option of `from` value")
		}
		kycTransfer.To = &to
	}

	val, ok = event.Data["token_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return TransferEvent{}, errors.New("invalid token_id value in event")
	}
	kycTransfer.TokenID = *val.U512

	return kycTransfer, nil
}
