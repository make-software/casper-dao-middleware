package repositories

import (
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// VotingRepository DB table interface
//
//go:generate mockgen -destination=../tests/mocks/voting_repo_mock.go -package=mocks -source=./voting.go VotingRepository
type VotingRepository interface {
	Save(changes *entities.Voting) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Voting, error)
	UpdateHasEnded(votingID uint32, hasEnded bool) error
}

type Voting struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewVoting(conn *sqlx.DB) *Voting {
	return &Voting{
		conn: conn,
		indexedFields: map[string]struct{}{
			"voting_id": {},
			"is_formal": {},
			"has_ended": {},
			"address":   {},
		},
	}
}

func (r *Voting) Save(voting *entities.Voting) error {
	queryBuilder := query.Insert("votings").
		Columns(
			"creator",
			"deploy_hash",
			"voting_id",
			"voting_type_id",
			"is_formal",
			"has_ended",
			"metadata",
			"config_double_time_between_votings",
			"config_total_onboarded",
			"config_voting_clearness_delta",
			"config_time_between_informal_and_formal_voting",
			"voting_quorum",
			"voting_time",
			"timestamp",
		).
		Values(
			voting.Creator,
			voting.DeployHash,
			voting.VotingID,
			voting.VotingTypeID,
			voting.IsFormal,
			voting.HasEnded,
			voting.Metadata,
			voting.ConfigDoubleTimeBetweenVotings,
			voting.ConfigTotalOnboarded,
			voting.ConfigVotingClearnessDelta,
			voting.ConfigTimeBetweenInformalAndFormalVoting,
			voting.VotingQuorum,
			voting.VotingTime,
			voting.Timestamp,
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

func (r *Voting) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("votings").
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

func (r *Voting) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Voting, error) {
	queryBuilder := query.Select("*").
		From("votings").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	votings := make([]*entities.Voting, 0)
	if err := r.conn.Select(&votings, sql, args...); err != nil {
		return nil, err
	}

	return votings, nil
}

func (r *Voting) UpdateHasEnded(votingID uint32, hasEnded bool) error {
	queryBuilder := query.Update("votings").
		Set("has_ended", hasEnded).
		Where(sq.Eq{"voting_id": votingID})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(sql, args...)
	return err
}
