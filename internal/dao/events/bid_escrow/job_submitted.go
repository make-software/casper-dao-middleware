package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"
)

const JobSubmittedEventName = "JobSubmitted"

type JobSubmittedEvent struct {
	BidID     uint32
	JobPoster casper.Hash
	Worker    casper.Hash
	Result    string
}

func ParseJobSubmittedEvent(event ces.Event) (JobSubmittedEvent, error) {
	var jobSubmitted JobSubmittedEvent

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return JobSubmittedEvent{}, errors.New("invalid bid_id value in event")
	}
	jobSubmitted.BidID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return JobSubmittedEvent{}, errors.New("invalid worker value in event")
	}
	if val.Key.Account != nil {
		jobSubmitted.Worker = val.Key.Account.Hash
	} else {
		jobSubmitted.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return JobSubmittedEvent{}, errors.New("invalid job_poster value in event")
	}
	if val.Key.Account != nil {
		jobSubmitted.JobPoster = val.Key.Account.Hash
	} else {
		jobSubmitted.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["result"]
	if !ok || val.Type != cltype.String {
		return JobSubmittedEvent{}, errors.New("invalid result value in event")
	}
	jobSubmitted.Result = val.StringVal.String()

	return jobSubmitted, nil
}
