package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobSubmittedEventName = "JobSubmitted"

type JobSubmittedEvent struct {
	BidID     uint32
	JobPoster casper_types.Key
	Worker    casper_types.Key
	Result    string
}

func ParseJobSubmittedEvent(event ces.Event) (JobSubmittedEvent, error) {
	var jobSubmitted JobSubmittedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobSubmittedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobSubmitted.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobSubmittedEvent{}, errors.New("invalid worker value in event")
	}
	jobSubmitted.Worker = *val.Key

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobSubmittedEvent{}, errors.New("invalid job_poster value in event")
	}
	jobSubmitted.JobPoster = *val.Key

	val, ok = event.Data["result"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDString {
		return JobSubmittedEvent{}, errors.New("invalid result value in event")
	}
	jobSubmitted.Result = *val.String

	return jobSubmitted, nil
}
