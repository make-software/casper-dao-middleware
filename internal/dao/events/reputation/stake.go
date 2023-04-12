package reputation

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const StakeEventName = "Stake"

type StakeEvent struct {
	BidID  uint32
	Worker casper_types.Key
	Amount casper_types.U512
}

func ParseStakeEvent(event ces.Event) (StakeEvent, error) {
	var stake StakeEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return StakeEvent{}, errors.New("invalid bid_id value in event")
	}
	stake.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return StakeEvent{}, errors.New("invalid worker value in event")
	}
	stake.Worker = *val.Key

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return StakeEvent{}, errors.New("invalid amount value in event")
	}
	stake.Amount = *val.U512

	return stake, nil
}
