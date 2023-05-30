package reputation

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const StakeEventName = "Stake"

type StakeEvent struct {
	BidID  uint32
	Worker casper.Hash
	Amount clvalue.UInt512
}

func ParseStakeEvent(event ces.Event) (StakeEvent, error) {
	var stake StakeEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return StakeEvent{}, errors.New("invalid bid_id value in event")
	}
	stake.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return StakeEvent{}, errors.New("invalid worker value in event")
	}

	if val.Key.Account != nil {
		stake.Worker = val.Key.Account.Hash
	} else {
		stake.Worker = *val.Key.Hash
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type != cltype.UInt512 {
		return StakeEvent{}, errors.New("invalid amount value in event")
	}
	stake.Amount = *val.UI512

	return stake, nil
}
