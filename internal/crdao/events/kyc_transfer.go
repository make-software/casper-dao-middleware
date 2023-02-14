package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const KYCTransfer = "Transfer"

type KycTransfer struct {
	From    dao_types.Address
	To      dao_types.Address
	TokenID types.U512
}

func ParseKYCTransferEvent(event ces.Event) (KycTransfer, error) {
	var kycTransfer KycTransfer

	val, ok := event.Data["from"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return KycTransfer{}, errors.New("invalid from value in event")
	}
	kycTransfer.From = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["to"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return KycTransfer{}, errors.New("invalid to value in event")
	}
	kycTransfer.From = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["token_id"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return KycTransfer{}, errors.New("invalid token_id value in event")
	}
	kycTransfer.TokenID = *val.U512

	return kycTransfer, nil
}
