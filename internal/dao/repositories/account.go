package repositories

import (
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/query"
)

// Account DB table interface
//
//go:generate mockgen -destination=../tests/mocks/account_repo_mock.go -package=mocks -source=./account.go AccountRepository
type Account interface {
	Upsert(account entities.Account) error
}

type account struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewAccount(conn *sqlx.DB) Account {
	return &account{
		conn: conn,
		indexedFields: map[string]struct{}{
			"hash": {},
		},
	}
}

func (r *account) Upsert(account entities.Account) error {
	queryBuilder := query.Insert("accounts").
		Columns(
			"hash",
			"is_kyc",
			"is_va",
			"timestamp",
		).
		Values(
			account.Hash,
			account.IsKyc,
			account.IsVA,
			account.Timestamp,
		).
		Suffix("ON DUPLICATE KEY UPDATE is_kyc = values(is_kyc), is_va = values(is_va)")

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
