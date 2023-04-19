package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
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
	BidID       uint32      `json:"bid_id" db:"bid_id"`
	JobPoster   types.Hash  `json:"job_poster" db:"job_poster"`
	Worker      types.Hash  `json:"worker" db:"worker"`
	Caller      *types.Hash `json:"caller" db:"caller"`
	Result      *string     `json:"result" db:"result"`
	DeployHash  types.Hash  `json:"deploy_hash" db:"deploy_hash"`
	JobStatusID JobStatusID `json:"job_status_id" db:"job_status_id"`
	FinishTime  uint64      `json:"finish_time"  db:"finish_time"`
	Timestamp   time.Time   `json:"timestamp"  db:"timestamp"`
}

func NewJob(
	bidID uint32,
	deployHash, jobPoster, worker types.Hash,
	finishTime uint64,
	jobStatusID JobStatusID,
	result *string,
	caller *types.Hash,
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
