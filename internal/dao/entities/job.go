package entities

import (
	"time"

	"github.com/make-software/casper-go-sdk/casper"
)

type JobStatusID byte

const (
	JobStatusIDCreated JobStatusID = iota + 1
	JobStatusIDSubmitted
	JobStatusIDCancelled
	JobStatusIDDone
	JobStatusIDRejected
)

type JobStatus struct {
	ID   JobStatusID `json:"id"`
	Name string      `json:"name"`
}

type Job struct {
	BidID       uint32       `json:"bid_id" db:"bid_id"`
	JobPoster   casper.Hash  `json:"job_poster" db:"job_poster"`
	Worker      casper.Hash  `json:"worker" db:"worker"`
	Caller      *casper.Hash `json:"caller" db:"caller"`
	Result      *string      `json:"result" db:"result"`
	DeployHash  casper.Hash  `json:"deploy_hash" db:"deploy_hash"`
	JobStatusID JobStatusID  `json:"job_status_id" db:"job_status_id"`
	FinishTime  uint64       `json:"finish_time"  db:"finish_time"`
	Timestamp   time.Time    `json:"timestamp"  db:"timestamp"`
}

func NewJob(
	bidID uint32,
	deployHash, jobPoster, worker casper.Hash,
	finishTime uint64,
	jobStatusID JobStatusID,
	result *string,
	caller *casper.Hash,
	timestamp time.Time) Job {
	return Job{
		BidID:       bidID,
		JobPoster:   jobPoster,
		Worker:      worker,
		DeployHash:  deployHash,
		JobStatusID: jobStatusID,
		Result:      result,
		Caller:      caller,
		FinishTime:  finishTime,
		Timestamp:   timestamp,
	}
}
