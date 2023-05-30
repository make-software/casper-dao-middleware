package repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/errors"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"
)

// Account DB table interface
//
//go:generate mockgen -destination=../tests/mocks/account_repo_mock.go -package=mocks -source=./account.go AccountRepository
type Account interface {
	UpsertIsKYC(account entities.Account) error
	UpsertIsVA(account entities.Account) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Account, error)
	FindByHash(hash casper.Hash) (*entities.Account, error)
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

func (r *account) UpsertIsKYC(account entities.Account) error {
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
		Suffix("ON DUPLICATE KEY UPDATE is_kyc = values(is_kyc)")

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

func (r *account) UpsertIsVA(account entities.Account) error {
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
		Suffix("ON DUPLICATE KEY UPDATE is_va = values(is_va)")

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

func (r *account) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("accounts").
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

func (r *account) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Account, error) {
	queryBuilder := query.Select("*").
		From("accounts").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	accounts := make([]*entities.Account, 0)
	if err := r.conn.Select(&accounts, sqlQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("not found account info by hash")
		}
		return nil, err
	}

	return accounts, nil
}

func (r *account) FindByHash(hash casper.Hash) (*entities.Account, error) {
	queryBuilder := query.Select("*").
		From("accounts").
		Where(sq.Eq{
			"hash": hash,
		})

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	account := entities.Account{}
	if err := r.conn.Get(&account, sqlQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("not found account info by hash")
		}
		return nil, err
	}

	return &account, nil
}
