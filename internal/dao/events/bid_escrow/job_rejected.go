package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobRejectedEventName = "JobRejected"

type JobRejectedEvent struct {
	BidID      uint32
	JobPoster  casper_types.Key
	Caller     casper_types.Key
	Worker     casper_types.Key
	CSPRAmount casper_types.U512
}

func ParseJobRejectedEvent(event ces.Event) (JobRejectedEvent, error) {
	var jobRejected JobRejectedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobRejectedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobRejected.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobRejectedEvent{}, errors.New("invalid worker value in event")
	}
	jobRejected.Worker = *val.Key

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobRejectedEvent{}, errors.New("invalid job_poster value in event")
	}
	jobRejected.JobPoster = *val.Key

	val, ok = event.Data["caller"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobRejectedEvent{}, errors.New("invalid caller value in event")
	}
	jobRejected.Caller = *val.Key

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return JobRejectedEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobRejected.CSPRAmount = *val.U512

	return jobRejected, nil
}
