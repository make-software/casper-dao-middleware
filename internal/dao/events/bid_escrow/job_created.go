package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const JobCreatedEventName = "JobCreated"

type JobCreatedEvent struct {
	BidID      uint32
	JobPoster  casper.Hash
	Worker     casper.Hash
	FinishTime uint64
	Payment    clvalue.UInt512
}

func ParseJobCreatedEvent(event ces.Event) (JobCreatedEvent, error) {
	var jobCreated JobCreatedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobCreatedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobCreated.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return JobCreatedEvent{}, errors.New("invalid worker value in event")
	}

	if val.Key.Account != nil {
		jobCreated.Worker = val.Key.Account.Hash
	} else {
		jobCreated.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobCreatedEvent{}, errors.New("invalid job_poster value in event")
	}

	if val.Key.Account != nil {
		jobCreated.JobPoster = val.Key.Account.Hash
	} else {
		jobCreated.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["finish_time"]
	if !ok || val.Type != cltype.UInt64 {
		return JobCreatedEvent{}, errors.New("invalid finish_time value in event")
	}
	jobCreated.FinishTime = val.UI64.Value() / 1000

	val, ok = event.Data["payment"]
	if !ok || val.Type != cltype.UInt512 {
		return JobCreatedEvent{}, errors.New("invalid payment value in event")
	}
	jobCreated.Payment = *val.UI512

	return jobCreated, nil
}
