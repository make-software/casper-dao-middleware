package repositories

import (
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Vote DB table interface
//
//go:generate mockgen -destination=../tests/mocks/vote_repo_mock.go -package=mocks -source=./vote.go Vote
type Vote interface {
	Save(changes *entities.Vote) error
	Count(filters map[string]interface{}) (uint64, error)
	Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Vote, error)
	CountVotesNumberForVotings(votingIDs []uint32) (map[uint32]uint32, error)
	UpdateIsCanceled(votingID uint32, address types.Hash, isCanceled bool) error
}

type vote struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}
}

func NewVote(conn *sqlx.DB) *vote {
	return &vote{
		conn: conn,
		indexedFields: map[string]struct{}{
			"address":   {},
			"voting_id": {},
			"is_formal": {},
		},
	}
}

func (r *vote) Save(vote *entities.Vote) error {
	queryBuilder := query.Insert("votes").
		Options("IGNORE").
		Columns(
			"deploy_hash",
			"voting_id",
			"address",
			"amount",
			"is_in_favour",
			"is_canceled",
			"timestamp",
		).
		Values(
			vote.DeployHash,
			vote.VotingID,
			vote.Address,
			vote.Amount,
			vote.IsInFavor,
			vote.IsCanceled,
			vote.Timestamp,
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

func (r *vote) Find(params *pagination.Params, filters map[string]interface{}) ([]*entities.Vote, error) {
	queryBuilder := query.Select("*").
		From("votes").
		FilterBy(filters, r.indexedFields).
		Paginate(params, r.indexedFields)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	infos := make([]*entities.Vote, 0)
	if err := r.conn.Select(&infos, sql, args...); err != nil {
		return nil, err
	}

	return infos, nil
}

func (r *vote) Count(filters map[string]interface{}) (uint64, error) {
	queryBuilder := query.Select("COUNT(*)").
		From("votes").
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

func (r *vote) CountVotesNumberForVotings(votingIDs []uint32) (map[uint32]uint32, error) {
	queryBuilder := query.Select("voting_id, COUNT(*)").
		From("votes").
		Where(sq.Eq{"voting_id": votingIDs}).
		GroupBy("voting_id")

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	result := make(map[uint32]uint32)

	rows, err := r.conn.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var votingID, votesNumber uint32
		if err := rows.Scan(&votingID, &votesNumber); err != nil {
			return nil, err
		}
		result[votingID] = votesNumber
	}

	return result, nil
}

func (r *vote) UpdateIsCanceled(votingID uint32, address types.Hash, isCanceled bool) error {
	queryBuilder := query.Update("votes").
		Set("is_canceled", isCanceled).
		Where(sq.Eq{
			"voting_id": votingID,
			"address":   address,
		})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(sql, args...)
	return err
}
