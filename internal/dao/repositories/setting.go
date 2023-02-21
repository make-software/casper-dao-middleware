package repositories

import (
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"

	"github.com/jmoiron/sqlx"
)

// SettingRepository DB table interface
//
//go:generate mockgen -destination=../tests/mocks/setting_repo_mock.go -package=mocks -source=./setting.go SettingRepository
type SettingRepository interface {
	Upsert(setting entities.Setting) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Setting, error)
}

type Setting struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewSetting(conn *sqlx.DB) *Setting {
	return &Setting{
		conn: conn,
		indexedFields: map[string]struct{}{
			"name": {},
		},
	}
}

func (r *Setting) Upsert(setting entities.Setting) error {
	queryBuilder := query.Insert("settings").
		Columns(
			"name",
			"value",
			"next_value",
			"activation_time",
		).
		Values(
			setting.Name,
			setting.Value,
			setting.NextValue,
			setting.ActivationTime,
		).
		Suffix("ON DUPLICATE KEY UPDATE value = values(value), next_value = values(next_value), activation_time = values(activation_time)")

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

func (r *Setting) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Setting, error) {
	queryBuilder := query.Select("*").
		From("settings").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	settings := make([]*entities.Setting, 0)
	if err := r.conn.Select(&settings, sql, args...); err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *Setting) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("settings").
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
