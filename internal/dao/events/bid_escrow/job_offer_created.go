package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobOfferCreatedEventName = "JobOfferCreated"

type JobOfferCreatedEvent struct {
	JobOfferID        uint32
	JobPoster         casper_types.Key
	MaxBudget         casper_types.U512
	ExpectedTimeFrame uint64
}

func ParseJobOfferCreatedEvent(event ces.Event) (JobOfferCreatedEvent, error) {
	var jobOfferCreated JobOfferCreatedEvent

	val, ok := event.Data["job_offer_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobOfferCreatedEvent{}, errors.New("invalid job_offer_id value in event")
	}
	jobOfferCreated.JobOfferID = *val.U32

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobOfferCreatedEvent{}, errors.New("invalid job_poster value in event")
	}
	jobOfferCreated.JobPoster = *val.Key

	val, ok = event.Data["max_budget"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return JobOfferCreatedEvent{}, errors.New("invalid max_budget value in event")
	}
	jobOfferCreated.MaxBudget = *val.U512

	val, ok = event.Data["expected_timeframe"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return JobOfferCreatedEvent{}, errors.New("invalid expected_timeframe value in event")
	}
	jobOfferCreated.ExpectedTimeFrame = *val.U64 / 1000

	return jobOfferCreated, nil
}
