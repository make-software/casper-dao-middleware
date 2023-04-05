package repositories

import (
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"
)

// JobOffer DB table interface
//
//go:generate mockgen -destination=../tests/mocks/job_offer_mock.go -package=mocks -source=./job_offer.go JobOffer
type JobOffer interface {
	Save(jobOffer *entities.JobOffer) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.JobOffer, error)
}

type jobOffer struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewJobOffer(conn *sqlx.DB) *jobOffer {
	return &jobOffer{
		conn:          conn,
		indexedFields: map[string]struct{}{},
	}
}

func (r *jobOffer) Save(jobOffer *entities.JobOffer) error {
	queryBuilder := query.Insert("job_offers").
		Options("IGNORE").
		Columns(
			"job_offer_id",
			"job_poster",
			"deploy_hash",
			"max_budget",
			"status",
			"auction_type",
			"expected_time_frame",
			"timestamp",
		).
		Values(
			jobOffer.JobOfferID,
			jobOffer.JobPoster,
			jobOffer.DeployHash,
			jobOffer.MaxBudget,
			jobOffer.Status,
			jobOffer.AuctionType,
			jobOffer.ExpectedTimeFrame,
			jobOffer.Timestamp,
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

func (r *jobOffer) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.JobOffer, error) {
	queryBuilder := query.Select("*").
		From("job_offers").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	offers := make([]*entities.JobOffer, 0)
	if err := r.conn.Select(&offers, sql, args...); err != nil {
		return nil, err
	}

	return offers, nil
}

func (r *jobOffer) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("job_offers").
		FilterBy(filters, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64

	row := r.conn.QueryRow(sql, args...)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
