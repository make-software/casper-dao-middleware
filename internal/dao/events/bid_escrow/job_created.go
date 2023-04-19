package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobCreatedEventName = "JobCreated"

type JobCreatedEvent struct {
	BidID      uint32
	JobPoster  casper_types.Key
	Worker     casper_types.Key
	FinishTime uint64
	Payment    casper_types.U512
}

func ParseJobCreatedEvent(event ces.Event) (JobCreatedEvent, error) {
	var jobCreated JobCreatedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobCreatedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobCreated.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobCreatedEvent{}, errors.New("invalid worker value in event")
	}
	jobCreated.Worker = *val.Key

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobCreatedEvent{}, errors.New("invalid job_poster value in event")
	}
	jobCreated.JobPoster = *val.Key

	val, ok = event.Data["finish_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return JobCreatedEvent{}, errors.New("invalid finish_time value in event")
	}
	jobCreated.FinishTime = *val.U64 / 1000

	val, ok = event.Data["payment"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return JobCreatedEvent{}, errors.New("invalid payment value in event")
	}
	jobCreated.Payment = *val.U512

	return jobCreated, nil
}
