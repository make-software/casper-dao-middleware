package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const JobDoneEventName = "JobDone"

type JobDoneEvent struct {
	BidID      uint32
	JobPoster  casper.Key
	Caller     casper.Hash
	Worker     casper.Hash
	CSPRAmount clvalue.UInt512
}

func ParseJobDoneEvent(event ces.Event) (JobDoneEvent, error) {
	var jobDone JobDoneEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobDoneEvent{}, errors.New("invalid bid_id value in event")
	}
	jobDone.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return JobDoneEvent{}, errors.New("invalid worker value in event")
	}
	if val.Key.Account != nil {
		jobDone.Worker = val.Key.Account.Hash
	} else {
		jobDone.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobDoneEvent{}, errors.New("invalid job_poster value in event")
	}
	jobDone.JobPoster = *val.Key

	val, ok = event.Data["caller"]
	if !ok || val.Type != cltype.Key {
		return JobDoneEvent{}, errors.New("invalid caller value in event")
	}
	if val.Key.Account != nil {
		jobDone.Caller = val.Key.Account.Hash
	} else {
		jobDone.Caller = *val.Key.Hash
	}

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type != cltype.UInt512 {
		return JobDoneEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobDone.CSPRAmount = *val.UI512

	return jobDone, nil
}
