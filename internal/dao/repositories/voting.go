package repositories

import (
	sq "github.com/Masterminds/squirrel"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"

	"github.com/jmoiron/sqlx"
)

// VotingRepository DB table interface
//
//go:generate mockgen -destination=../tests/mocks/voting_repo_mock.go -package=mocks -source=./voting.go VotingRepository
type VotingRepository interface {
	Save(changes *entities.Voting) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Voting, error)
	GetByVotingID(votingID uint32) (*entities.Voting, error)
	Update(voting *entities.Voting) error
	UpdateIsCanceled(votingID uint32, isCanceled bool) error
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
		Options("IGNORE").
		Columns(
			"creator",
			"deploy_hash",
			"voting_id",
			"voting_type_id",
			"informal_voting_quorum",
			"informal_voting_starts_at",
			"informal_voting_ends_at",
			"formal_voting_quorum",
			"formal_voting_starts_at",
			"formal_voting_ends_at",
			"metadata",
			"is_canceled",
			"informal_voting_result",
			"formal_voting_result",
			"config_total_onboarded",
			"config_voting_clearness_delta",
			"config_time_between_informal_and_formal_voting",
		).
		Values(
			voting.Creator,
			voting.DeployHash,
			voting.VotingID,
			voting.VotingTypeID,
			voting.InformalVotingQuorum,
			voting.InformalVotingStartsAt,
			voting.InformalVotingEndsAt,
			voting.FormalVotingQuorum,
			voting.FormalVotingStartsAt,
			voting.FormalVotingEndsAt,
			voting.Metadata,
			voting.IsCanceled,
			voting.InformalVotingResult,
			voting.FormalVotingResult,
			voting.ConfigTotalOnboarded,
			voting.ConfigVotingClearnessDelta,
			voting.ConfigTimeBetweenInformalAndFormalVoting,
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

func (r *Voting) GetByVotingID(votingID uint32) (*entities.Voting, error) {
	queryBuilder := query.Select("*").
		From("votings").
		Where(sq.Eq{
			"voting_id": votingID,
		})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var voting entities.Voting
	if err := r.conn.Get(&voting, sql, args...); err != nil {
		return nil, err
	}

	return &voting, nil
}

func (r *Voting) Update(voting *entities.Voting) error {
	queryBuilder := query.Update("votings").
		SetMap(map[string]interface{}{
			"formal_voting_starts_at": voting.FormalVotingStartsAt,
			"formal_voting_ends_at":   voting.FormalVotingEndsAt,
			"informal_voting_result":  voting.InformalVotingResult,
			"formal_voting_result":    voting.FormalVotingResult,
		})

	queryBuilder = queryBuilder.
		Where(sq.Eq{
			"voting_id": voting.VotingID,
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

func (r *Voting) UpdateIsCanceled(votingID uint32, isCanceled bool) error {
	queryBuilder := query.Update("votings").
		Set("is_canceled", isCanceled).
		Where(sq.Eq{
			"voting_id": votingID,
		})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(sql, args...)
	return err
}
