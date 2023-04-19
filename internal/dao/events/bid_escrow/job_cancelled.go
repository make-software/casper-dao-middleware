package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobCancelledEventName = "JobCancelled"

type JobCancelledEvent struct {
	BidID      uint32
	JobPoster  casper_types.Key
	Caller     casper_types.Key
	Worker     casper_types.Key
	CSPRAmount casper_types.U512
}

func ParseJobCancelledEvent(event ces.Event) (JobCancelledEvent, error) {
	var jobCancelled JobCancelledEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobCancelledEvent{}, errors.New("invalid bid_id value in event")
	}
	jobCancelled.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobCancelledEvent{}, errors.New("invalid worker value in event")
	}
	jobCancelled.Worker = *val.Key

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobCancelledEvent{}, errors.New("invalid job_poster value in event")
	}
	jobCancelled.JobPoster = *val.Key

	val, ok = event.Data["caller"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobCancelledEvent{}, errors.New("invalid caller value in event")
	}
	jobCancelled.Caller = *val.Key

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return JobCancelledEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobCancelled.CSPRAmount = *val.U512

	return jobCancelled, nil
}
