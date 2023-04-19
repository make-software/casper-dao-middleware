package reputation

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const UnstakeEventName = "Unstake"

type UnstakeEvent struct {
	BidID  uint32
	Worker casper_types.Key
	Amount casper_types.U512
}

func ParseUnstakeEvent(event ces.Event) (UnstakeEvent, error) {
	var ustake UnstakeEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return UnstakeEvent{}, errors.New("invalid bid_id value in event")
	}
	ustake.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return UnstakeEvent{}, errors.New("invalid worker value in event")
	}
	ustake.Worker = *val.Key

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return UnstakeEvent{}, errors.New("invalid amount value in event")
	}
	ustake.Amount = *val.U512

	return ustake, nil
}
