package repositories

import (
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
)

// Job DB table interface
//
//go:generate mockgen -destination=../tests/mocks/job_mock.go -package=mocks -source=./job.go Job
type Job interface {
	Save(job *entities.Job) error
	GetByBidID(bidID uint32) (*entities.Job, error)
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Job, error)
	Update(job *entities.Job) error
}

type job struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func (j job) Save(job *entities.Job) error {
	//TODO implement me
	panic("implement me")
}

func (j job) Count(filters map[string]interface{}) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (j job) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Job, error) {
	//TODO implement me
	panic("implement me")
}

func (j job) Update(job *entities.Job) error {
	//TODO implement me
	panic("implement me")
}

func NewJob(conn *sqlx.DB) *job {
	return &job{
		conn:          conn,
		indexedFields: map[string]struct{}{},
	}
}
