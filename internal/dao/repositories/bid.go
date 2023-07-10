package repositories

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"
)

// Bid DB table interface
//
//go:generate mockgen -destination=../tests/mocks/bid_mock.go -package=mocks -source=./bid.go Bid
type Bid interface {
	Save(bid *entities.Bid) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Bid, error)
	UpdateIsPickedBy(bidID uint32, isPickedBy bool) error
}

type bid struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewBid(conn *sqlx.DB) *bid {
	return &bid{
		conn: conn,
		indexedFields: map[string]struct{}{
			"job_offer_id": {},
		},
	}
}

func (r *bid) Save(bid *entities.Bid) error {
	queryBuilder := query.Insert("bids").
		Options("IGNORE").
		Columns(
			"job_offer_id",
			"bid_id",
			"worker",
			"deploy_hash",
			"onboard",
			"proposed_time_frame",
			"proposed_payment",
			"picked_by_job_poster",
			"reputation_stake",
			"cspr_stake",
			"timestamp",
		).
		Values(
			bid.JobOfferID,
			bid.BidID,
			bid.Worker,
			bid.DeployHash,
			bid.Onboard,
			bid.ProposedTimeFrame,
			bid.ProposedPayment,
			bid.PickedByJobPoster,
			bid.ReputationStake,
			bid.CSPRStake,
			bid.Timestamp,
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

func (r *bid) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Bid, error) {
	queryBuilder := query.Select("*").
		From("bids").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	bids := make([]*entities.Bid, 0)
	if err := r.conn.Select(&bids, sql, args...); err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *bid) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("bids").
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

func (r *bid) UpdateIsPickedBy(bidID uint32, isPickedBy bool) error {
	queryBuilder := query.Update("bids").
		SetMap(map[string]interface{}{
			"picked_by_job_poster": isPickedBy,
		})

	queryBuilder = queryBuilder.
		Where(sq.Eq{
			"bid_id": bidID,
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
