package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const Transfer = "Transfer"

type VaTransfer struct {
	From    dao_types.Address
	To      dao_types.Address
	TokenID types.U512
}

func ParseVaTransferEvent(event ces.Event) (VaTransfer, error) {
	var vaTransfer VaTransfer

	val, ok := event.Data["from"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return VaTransfer{}, errors.New("invalid from value in event")
	}
	vaTransfer.From = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["to"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return VaTransfer{}, errors.New("invalid to value in event")
	}
	vaTransfer.From = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["token_id"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return VaTransfer{}, errors.New("invalid token_id value in event")
	}
	vaTransfer.TokenID = *val.U512

	return vaTransfer, nil
}
