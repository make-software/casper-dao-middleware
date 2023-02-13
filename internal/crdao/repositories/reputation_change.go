package repositories

import (
	"fmt"
	"strings"

	"casper-dao-middleware/internal/crdao/entities"
	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/query"

	"github.com/jmoiron/sqlx"
)

// ReputationChangeRepository DB table interface
//
//go:generate mockgen -destination=../tests/mocks/reputation_change_repo_mock.go -package=mocks -source=./reputation_change.go ReputationChangeRepository
type ReputationChangeRepository interface {
	SaveBatch(changes []entities.ReputationChange) error
	CalculateTotalReputationForAddress(address types.Hash) (entities.TotalReputation, error)
	FindAggregatedReputationChanges(params *pagination.Params, filters map[string]interface{}) ([]entities.AggregatedReputationChange, error)
	CountAggregatedReputationChanges(filters map[string]interface{}) (uint64, error)
}

type ReputationChange struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}

	contractPackageHashes dao_types.DAOContractsMetadata
}

func NewReputationChange(conn *sqlx.DB, hashes dao_types.DAOContractsMetadata) *ReputationChange {
	return &ReputationChange{
		conn: conn,
		indexedFields: map[string]struct{}{
			"address": {},
		},
		contractPackageHashes: hashes,
	}
}

func (r *ReputationChange) SaveBatch(changes []entities.ReputationChange) error {
	columns := []string{
		"address",
		"contract_package_hash",
		"voting_id",
		"amount",
		"deploy_hash",
		"reason",
		"timestamp",
	}

	insertQuery := `INSERT INTO reputation_changes (` + strings.Join(columns, ",") + `) 
		VALUES (:` + strings.Join(columns, ",:") + `)`

	_, err := r.conn.NamedExec(insertQuery, changes)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReputationChange) CalculateTotalReputationForAddress(address types.Hash) (entities.TotalReputation, error) {
	query := `
	SELECT 
	    (SELECT SUM(amount) FROM reputation_changes WHERE address = ? and contract_package_hash = ?) as available_amount, 
	    (SELECT SUM(amount)  FROM reputation_changes WHERE address = ? and contract_package_hash = ?) as staked_amount  
	FROM reputation_changes;
`

	totalReputation := entities.TotalReputation{}
	err := r.conn.Get(&totalReputation, query,
		address, r.contractPackageHashes.ReputationContractPackageHash,
		address, r.contractPackageHashes.VoterContractPackageHash)
	if err != nil {
		return entities.TotalReputation{}, err
	}

	return totalReputation, nil
}

func (r *ReputationChange) FindAggregatedReputationChanges(params *pagination.Params, filters map[string]interface{}) ([]entities.AggregatedReputationChange, error) {
	args := make([]interface{}, 0)
	aggregatedBuilder := query.Select(
		`address, 
				  voting_id, 
				  sum(IF(reason = 3 and contract_package_hash = ?, -amount, 0)) as staked_amount,  
			      sum(IF(reason = 4 and contract_package_hash = ?, amount, 0)) as released_amount`,
	).
		From("reputation_changes").
		GroupBy("address, IFNULL(voting_id, deploy_hash)").
		Paginate(params, r.indexedFields).
		FilterBy(filters, r.indexedFields)

	aggregationSql, aggregationArgs, err := aggregatedBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	args = append(args, r.contractPackageHashes.ReputationContractPackageHash, r.contractPackageHashes.VoterContractPackageHash)
	args = append(args, aggregationArgs...)

	queryBuilder := query.Select(
		`aggregated_changes.address, 
				  aggregated_changes.voting_id, 
                  CURRENT_TIMESTAMP() as timestamp,
                  aggregated_changes.staked_amount, 
                  aggregated_changes.released_amount,
                  IF(aggregated_changes.staked_amount > aggregated_changes.released_amount, aggregated_changes.staked_amount - aggregated_changes.released_amount, 0) as lost_amount, 
                  IF(aggregated_changes.released_amount > aggregated_changes.staked_amount,  aggregated_changes.released_amount - aggregated_changes.staked_amount, 0) as earned_amount`,
	).From(fmt.Sprintf("(%s) as aggregated_changes", aggregationSql))

	sql, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	aggregatedChanges := make([]entities.AggregatedReputationChange, 0)
	if err := r.conn.Select(&aggregatedChanges, sql, args...); err != nil {
		return nil, err
	}

	return aggregatedChanges, nil
}

func (r *ReputationChange) CountAggregatedReputationChanges(filters map[string]interface{}) (uint64, error) {
	args := make([]interface{}, 0)
	aggregatedBuilder := query.Select("voting_id").
		From("reputation_changes").
		GroupBy("IFNULL(voting_id, deploy_hash)").
		FilterBy(filters, r.indexedFields)

	aggregationSql, aggregationArgs, err := aggregatedBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	args = append(args, aggregationArgs...)

	queryBuilder := query.Select("COUNT(*)").From(fmt.Sprintf("(%s) as aggregated_changes", aggregationSql))

	sql, _, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64
	row := r.conn.QueryRow(sql, args...)
	if err = row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
