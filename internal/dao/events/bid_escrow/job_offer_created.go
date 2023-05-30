package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const JobOfferCreatedEventName = "JobOfferCreated"

type JobOfferCreatedEvent struct {
	JobOfferID        uint32
	JobPoster         casper.Hash
	MaxBudget         clvalue.UInt512
	ExpectedTimeFrame uint64
}

func ParseJobOfferCreatedEvent(event ces.Event) (JobOfferCreatedEvent, error) {
	var jobOfferCreated JobOfferCreatedEvent

	val, ok := event.Data["job_offer_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobOfferCreatedEvent{}, errors.New("invalid job_offer_id value in event")
	}
	jobOfferCreated.JobOfferID = val.UI32.Value()

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobOfferCreatedEvent{}, errors.New("invalid job_poster value in event")
	}
	if val.Key.Account != nil {
		jobOfferCreated.JobPoster = val.Key.Account.Hash
	} else {
		jobOfferCreated.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["max_budget"]
	if !ok || val.Type != cltype.UInt512 {
		return JobOfferCreatedEvent{}, errors.New("invalid max_budget value in event")
	}
	jobOfferCreated.MaxBudget = *val.UI512

	val, ok = event.Data["expected_timeframe"]
	if !ok || val.Type != cltype.UInt64 {
		return JobOfferCreatedEvent{}, errors.New("invalid expected_timeframe value in event")
	}
	jobOfferCreated.ExpectedTimeFrame = val.UI64.Value() / 1000

	return jobOfferCreated, nil
}
