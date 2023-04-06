package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const BidSubmittedEventName = "BidSubmitted"

type BidSubmittedEvent struct {
	JobOfferID        uint32
	BidID             uint32
	Worker            casper_types.Key
	Onboard           bool
	ProposedTimeFrame uint64
	ProposedPayment   casper_types.U512
	ReputationStake   *casper_types.U512
	CSPRStake         *casper_types.U512
}

func ParseBidSubmittedEvent(event ces.Event) (BidSubmittedEvent, error) {
	var bidSubmitted BidSubmittedEvent

	val, ok := event.Data["job_offer_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return BidSubmittedEvent{}, errors.New("invalid job_offer_id value in event")
	}
	bidSubmitted.JobOfferID = *val.U32

	val, ok = event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return BidSubmittedEvent{}, errors.New("invalid bid_id value in event")
	}
	bidSubmitted.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return BidSubmittedEvent{}, errors.New("invalid worker value in event")
	}
	bidSubmitted.Worker = *val.Key

	val, ok = event.Data["onboard"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDBool {
		return BidSubmittedEvent{}, errors.New("invalid onboard value in event")
	}
	bidSubmitted.Onboard = *val.Bool

	val, ok = event.Data["proposed_timeframe"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return BidSubmittedEvent{}, errors.New("invalid proposed_timeframe value in event")
	}
	bidSubmitted.ProposedTimeFrame = *val.U64 / 1000

	val, ok = event.Data["proposed_payment"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return BidSubmittedEvent{}, errors.New("invalid proposed_payment value in event")
	}
	bidSubmitted.ProposedPayment = *val.U512

	val, ok = event.Data["reputation_stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return BidSubmittedEvent{}, errors.New("invalid reputation_stake value in event")
	}

	if val.Option != nil {
		if val.Option.Type.CLTypeID != casper_types.CLTypeIDU512 {
			return BidSubmittedEvent{}, errors.New("invalid value inside option of `reputation_stake` value")
		}

		bidSubmitted.ReputationStake = val.Option.U512
	}

	val, ok = event.Data["cspr_stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return BidSubmittedEvent{}, errors.New("invalid cspr_stake value in event")
	}

	if val.Option != nil {
		if val.Option.Type.CLTypeID != casper_types.CLTypeIDU512 {
			return BidSubmittedEvent{}, errors.New("invalid value inside option of `cspr_stake` value")
		}

		bidSubmitted.CSPRStake = val.Option.U512
	}

	return bidSubmitted, nil
}
