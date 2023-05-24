package reputation

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const UnstakeEventName = "Unstake"

type UnstakeEvent struct {
	BidID  uint32
	Worker casper.Hash
	Amount clvalue.UInt512
}

func ParseUnstakeEvent(event ces.Event) (UnstakeEvent, error) {
	var ustake UnstakeEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return UnstakeEvent{}, errors.New("invalid bid_id value in event")
	}
	ustake.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return UnstakeEvent{}, errors.New("invalid worker value in event")
	}
	if val.Key.Account != nil {
		ustake.Worker = val.Key.Account.Hash
	} else {
		ustake.Worker = *val.Key.Hash
	}

	val, ok = event.Data["amount"]
	if !ok || val.Type != cltype.UInt512 {
		return UnstakeEvent{}, errors.New("invalid amount value in event")
	}
	ustake.Amount = *val.UI512

	return ustake, nil
}
