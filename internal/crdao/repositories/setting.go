package repositories

import (
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/pkg/query"

	"github.com/jmoiron/sqlx"
)

// SettingRepository DB table interface
//
//go:generate mockgen -destination=../tests/mocks/setting_repo_mock.go -package=mocks -source=./setting.go SettingRepository
type SettingRepository interface {
	Upsert(setting entities.Setting) error
}

type Setting struct {
	conn *sqlx.DB
}

func NewSetting(conn *sqlx.DB) *Setting {
	return &Setting{
		conn: conn,
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
