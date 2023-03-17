package repositories

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
)

// ReputationTotal DB table interface
//
//go:generate mockgen -destination=../tests/mocks/reputation__repo_mock.go -package=mocks -source=./reputation_change.go ReputationTotal
type ReputationTotal interface {
	SaveBatch(totals []entities.ReputationTotal) error
}

type reputationTotal struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewReputationTotal(conn *sqlx.DB) *reputationTotal {
	return &reputationTotal{
		conn: conn,
		indexedFields: map[string]struct{}{
			"address": {},
		},
	}
}

func (r reputationTotal) SaveBatch(totals []entities.ReputationTotal) error {
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

	insertQuery := `INSERT IGNORE INTO reputation_totals (` + strings.Join(columns, ",") + `) 
		VALUES (:` + strings.Join(columns, ",:") + `)`

	_, err := r.conn.NamedExec(insertQuery, totals)
	if err != nil {
		return err
	}

	return nil
}
