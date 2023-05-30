package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const JobRejectedEventName = "JobRejected"

type JobRejectedEvent struct {
	BidID      uint32
	JobPoster  casper.Hash
	Caller     casper.Hash
	Worker     casper.Hash
	CSPRAmount clvalue.UInt512
}

func ParseJobRejectedEvent(event ces.Event) (JobRejectedEvent, error) {
	var jobRejected JobRejectedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobRejectedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobRejected.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return JobRejectedEvent{}, errors.New("invalid worker value in event")
	}
	if val.Key.Account != nil {
		jobRejected.Worker = val.Key.Account.Hash
	} else {
		jobRejected.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobRejectedEvent{}, errors.New("invalid job_poster value in event")
	}
	if val.Key.Account != nil {
		jobRejected.JobPoster = val.Key.Account.Hash
	} else {
		jobRejected.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["caller"]
	if !ok || val.Type != cltype.Key {
		return JobRejectedEvent{}, errors.New("invalid caller value in event")
	}
	if val.Key.Account != nil {
		jobRejected.Caller = val.Key.Account.Hash
	} else {
		jobRejected.Caller = *val.Key.Hash
	}

	val, ok = event.Data["cspr_amount"]
	if !ok || val.Type != cltype.UInt512 {
		return JobRejectedEvent{}, errors.New("invalid cspr_amount value in event")
	}
	jobRejected.CSPRAmount = *val.UI512

	return jobRejected, nil
}
