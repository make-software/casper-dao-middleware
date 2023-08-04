package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"
)

const BidSubmittedEventName = "BidSubmitted"

type BidSubmittedEvent struct {
	JobOfferID        uint32
	BidID             uint32
	Worker            casper.Hash
	Onboard           bool
	ProposedTimeFrame uint64
	ProposedPayment   clvalue.UInt512
	ReputationStake   *clvalue.UInt512
	CSPRStake         *clvalue.UInt512
}

func ParseBidSubmittedEvent(event ces.Event) (BidSubmittedEvent, error) {
	var bidSubmitted BidSubmittedEvent

	val, ok := event.Data["job_offer_id"]
	if !ok || val.Type != cltype.UInt32 {
		return BidSubmittedEvent{}, errors.New("invalid job_offer_id value in event")
	}
	bidSubmitted.JobOfferID = val.UI32.Value()

	val, ok = event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return BidSubmittedEvent{}, errors.New("invalid bid_id value in event")
	}
	bidSubmitted.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return BidSubmittedEvent{}, errors.New("invalid worker value in event")
	}

	if val.Key.Account != nil {
		bidSubmitted.Worker = val.Key.Account.Hash
	} else {
		bidSubmitted.Worker = *val.Key.Hash
	}

	val, ok = event.Data["onboard"]
	if !ok || val.Type != cltype.Bool {
		return BidSubmittedEvent{}, errors.New("invalid onboard value in event")
	}
	bidSubmitted.Onboard = val.Bool.Value()

	val, ok = event.Data["proposed_timeframe"]
	if !ok || val.Type != cltype.UInt64 {
		return BidSubmittedEvent{}, errors.New("invalid proposed_timeframe value in event")
	}
	bidSubmitted.ProposedTimeFrame = val.UI64.Value() / 1000

	val, ok = event.Data["proposed_payment"]
	if !ok || val.Type != cltype.UInt512 {
		return BidSubmittedEvent{}, errors.New("invalid proposed_payment value in event")
	}
	bidSubmitted.ProposedPayment = *val.UI512

	val, ok = event.Data["reputation_stake"]
	if !ok {
		return BidSubmittedEvent{}, errors.New("invalid reputation_stake value in event")
	}

	if val.Option != nil && val.Option.Inner != nil {
		if val.Option.Inner.Type != cltype.UInt512 {
			return BidSubmittedEvent{}, errors.New("invalid value inside option of `reputation_stake` value")
		}

		bidSubmitted.ReputationStake = val.Option.Inner.UI512
	}

	val, ok = event.Data["cspr_stake"]
	if !ok {
		return BidSubmittedEvent{}, errors.New("invalid cspr_stake value in event")
	}

	if val.Option != nil && val.Option.Inner != nil {
		if val.Option.Inner.Type != cltype.UInt512 {
			return BidSubmittedEvent{}, errors.New("invalid value inside option of `cspr_stake` value")
		}

		bidSubmitted.CSPRStake = val.Option.Inner.UI512
	}

	return bidSubmitted, nil
}
