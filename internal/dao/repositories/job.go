package repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/errors"
	"casper-dao-middleware/pkg/query"
)

// Job DB table interface
//
//go:generate mockgen -destination=../tests/mocks/job_mock.go -package=mocks -source=./job.go Job
type Job interface {
	Save(job *entities.Job) error
	GetByBidID(bidID uint32) (*entities.Job, error)
	Update(job *entities.Job) error
}

type job struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewJob(conn *sqlx.DB) *job {
	return &job{
		conn:          conn,
		indexedFields: map[string]struct{}{},
	}
}

func (r job) GetByBidID(bidID uint32) (*entities.Job, error) {
	queryBuilder := query.Select("*").
		From("jobs").
		Where(sq.Eq{
			"bid_id": bidID,
		})

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	job := entities.Job{}
	if err := r.conn.Get(&job, sqlQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("not found job info by bid_id")
		}
		return nil, err
	}

	return &job, nil
}

func (r job) Save(job *entities.Job) error {
	queryBuilder := query.Insert("jobs").
		Options("IGNORE").
		Columns(
			"bid_id",
			"job_poster",
			"worker",
			"caller",
			"result",
			"deploy_hash",
			"job_status",
			"finish_time",
			"timestamp",
		).
		Values(
			job.BidID,
			job.JobPoster,
			job.Worker,
			job.Caller,
			job.Result,
			job.DeployHash,
			job.JobStatus,
			job.FinishTime,
			job.Timestamp,
		)
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r job) Update(job *entities.Job) error {
	queryBuilder := query.Update("jobs").
		SetMap(map[string]interface{}{
			"caller":     job.Caller,
			"result":     job.Result,
			"job_status": job.JobStatus,
		})

	queryBuilder = queryBuilder.
		Where(sq.Eq{
			"bid_id": job.BidID,
		})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}
