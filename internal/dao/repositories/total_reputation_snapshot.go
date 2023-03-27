package repositories

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"
)

// TotalReputationSnapshot DB table interface
//
//go:generate mockgen -destination=../tests/mocks/reputation__repo_mock.go -package=mocks -source=./total_reputation_snapshot.go TotalReputationSnapshot
type TotalReputationSnapshot interface {
	SaveBatch(snapshots []entities.TotalReputationSnapshot) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]entities.TotalReputationSnapshot, error)
}

type totalReputationSnapshot struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewTotalReputationSnapshot(conn *sqlx.DB) TotalReputationSnapshot {
	return &totalReputationSnapshot{
		conn: conn,
		indexedFields: map[string]struct{}{
			"address":   {},
			"timestamp": {},
		},
	}
}

func (r *totalReputationSnapshot) SaveBatch(snapshots []entities.TotalReputationSnapshot) error {
	columns := []string{
		"address",
		"total_liquid_reputation",
		"total_staked_reputation",
		"voting_lost_reputation",
		"voting_earned_reputation",
		"voting_id",
		"deploy_hash",
		"reason",
		"timestamp",
	}

	insertQuery := `INSERT IGNORE INTO total_reputation_snapshots (` + strings.Join(columns, ",") + `) 
		VALUES (:` + strings.Join(columns, ",:") + `)`

	_, err := r.conn.NamedExec(insertQuery, snapshots)
	if err != nil {
		return err
	}

	return nil
}

func (r *totalReputationSnapshot) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("total_reputation_snapshots").
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

func (r *totalReputationSnapshot) Find(params *pagination.Params, filters map[string]interface{}) ([]entities.TotalReputationSnapshot, error) {
	queryBuilder := query.Select("*").
		From("total_reputation_snapshots").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	snapshots := make([]entities.TotalReputationSnapshot, 0)
	if err := r.conn.Select(&snapshots, sql, args...); err != nil {
		return nil, err
	}

	return snapshots, nil
}
