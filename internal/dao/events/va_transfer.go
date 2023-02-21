package events

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const Transfer = "Transfer"

type VaTransfer struct {
	From    types.Address
	To      types.Address
	TokenID casper_types.U512
}

func ParseVaTransferEvent(event ces.Event) (VaTransfer, error) {
	var vaTransfer VaTransfer

	val, ok := event.Data["from"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return VaTransfer{}, errors.New("invalid from value in event")
	}
	vaTransfer.From = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["to"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return VaTransfer{}, errors.New("invalid to value in event")
	}
	vaTransfer.From = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["token_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VaTransfer{}, errors.New("invalid token_id value in event")
	}
	vaTransfer.TokenID = *val.U512

	return vaTransfer, nil
}
