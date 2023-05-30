package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"
)

const JobCancelledEventName = "JobCancelled"

type JobCancelledEvent struct {
	BidID      uint32
	JobPoster  casper.Hash
	Caller     casper.Hash
	Worker     casper.Hash
	CSPRAmount clvalue.UInt512
}

func ParseJobCancelledEvent(event ces.Event) (JobCancelledEvent, error) {
	var jobCancelled JobCancelledEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobCancelledEvent{}, errors.New("invalid bid_id value in event")
	}
	jobCancelled.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return JobCancelledEvent{}, errors.New("invalid worker value in event")
	}
	if val.Key.Account != nil {
		jobCancelled.Worker = val.Key.Account.Hash
	} else {
		jobCancelled.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobCancelledEvent{}, errors.New("invalid job_poster value in event")
	}

	if val.Key.Account != nil {
		jobCancelled.JobPoster = val.Key.Account.Hash
	} else {
		jobCancelled.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["caller"]
	if !ok || val.Type != cltype.Key {
		return JobCancelledEvent{}, errors.New("invalid caller value in event")
	}

	if val.Key.Account != nil {
		jobCancelled.Caller = val.Key.Account.Hash
	} else {
		jobCancelled.Caller = *val.Key.Hash
	}

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type != cltype.UInt512 {
		return JobCancelledEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobCancelled.CSPRAmount = *val.UI512

	return jobCancelled, nil
}
