package bid_escrow

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const JobDoneEventName = "JobDone"

type JobDoneEvent struct {
	BidID      uint32
	JobPoster  casper_types.Key
	Caller     casper_types.Key
	Worker     casper_types.Key
	CSPRAmount casper_types.U512
}

func ParseJobDoneEvent(event ces.Event) (JobDoneEvent, error) {
	var jobDone JobDoneEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return JobDoneEvent{}, errors.New("invalid bid_id value in event")
	}
	jobDone.BidID = *val.U32

	val, ok = event.Data["worker"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobDoneEvent{}, errors.New("invalid worker value in event")
	}
	jobDone.Worker = *val.Key

	val, ok = event.Data["job_poster"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobDoneEvent{}, errors.New("invalid job_poster value in event")
	}
	jobDone.JobPoster = *val.Key

	val, ok = event.Data["caller"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return JobDoneEvent{}, errors.New("invalid caller value in event")
	}
	jobDone.Caller = *val.Key

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return JobDoneEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobDone.CSPRAmount = *val.U512

	return jobDone, nil
}
