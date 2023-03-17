package repositories

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/casper/types"
)

// ReputationChange DB table interface
//
//go:generate mockgen -destination=../tests/mocks/reputation_change_repo_mock.go -package=mocks -source=./reputation_change.go ReputationChange
type ReputationChange interface {
	SaveBatch(changes []entities.ReputationChange) error
	CalculateLiquidStakeReputationForAddress(address types.Hash) (entities.LiquidStakeReputation, error)
	CalculateAggregatedLiquidStakeReputationForAddresses(addresses []types.Hash) ([]entities.LiquidStakeReputation, error)
}

type reputationChange struct {
	conn          *sqlx.DB
	indexedFields map[string]struct{}

	contractPackageHashes utils.DAOContractsMetadata
}

func NewReputationChange(conn *sqlx.DB, hashes utils.DAOContractsMetadata) *reputationChange {
	return &reputationChange{
		conn: conn,
		indexedFields: map[string]struct{}{
			"address": {},
		},
		contractPackageHashes: hashes,
	}
}

func (r *reputationChange) SaveBatch(changes []entities.ReputationChange) error {
	columns := []string{
		"address",
		"contract_package_hash",
		"voting_id",
		"amount",
		"deploy_hash",
		"reason",
		"timestamp",
	}

	insertQuery := `INSERT IGNORE INTO reputation_changes (` + strings.Join(columns, ",") + `) 
		VALUES (:` + strings.Join(columns, ",:") + `)`

	_, err := r.conn.NamedExec(insertQuery, changes)
	if err != nil {
		return err
	}

	return nil
}

func (r *reputationChange) CalculateLiquidStakeReputationForAddress(address types.Hash) (entities.LiquidStakeReputation, error) {
	query := `
	SELECT 
	    (SELECT SUM(amount) FROM reputation_changes WHERE contract_package_hash = ?) as liquid_amount, 
	    (SELECT SUM(amount)  FROM reputation_changes WHERE contract_package_hash != ?) as staked_amount  
	FROM reputation_changes WHERE address = ?;
`

	args := []interface{}{
		r.contractPackageHashes.ReputationContractPackageHash,
		r.contractPackageHashes.ReputationContractPackageHash,
		address,
	}

	liquidStakeReputation := entities.LiquidStakeReputation{}
	err := r.conn.Get(&liquidStakeReputation, query, args...)
	if err != nil {
		return entities.LiquidStakeReputation{}, err
	}

	return liquidStakeReputation, nil
}

func (r *reputationChange) CalculateAggregatedLiquidStakeReputationForAddresses(addresses []types.Hash) ([]entities.LiquidStakeReputation, error) {
	query := `
	SELECT 
	    (SELECT SUM(amount) FROM reputation_changes WHERE contract_package_hash = ?) as liquid_amount, 
	    (SELECT SUM(amount)  FROM reputation_changes WHERE contract_package_hash != ?) as staked_amount,
	    address
	FROM reputation_changes WHERE address in (?) GROUP BY address;
`

	args := []interface{}{
		r.contractPackageHashes.ReputationContractPackageHash,
		r.contractPackageHashes.ReputationContractPackageHash,
	}
	for _, address := range addresses {
		args = append(args, address)
	}

	liquidStakeReputations := make([]entities.LiquidStakeReputation, 0)
	err := r.conn.Select(&liquidStakeReputations, query, args...)
	if err != nil {
		return nil, err
	}

	return liquidStakeReputations, nil
}
