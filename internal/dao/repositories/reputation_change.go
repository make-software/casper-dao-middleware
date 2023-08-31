package repositories

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/utils"
)

// ReputationChange DB table interface
//
//go:generate mockgen -destination=../tests/mocks/reputation_change_repo_mock.go -package=mocks -source=./reputation_change.go ReputationChange
type ReputationChange interface {
	SaveBatch(changes []entities.ReputationChange) error
	CalculateLiquidStakeReputationForAddress(address casper.Hash) (entities.LiquidStakeReputation, error)
	CalculateAggregatedLiquidStakeReputationForAddresses(addresses []casper.Hash) ([]entities.LiquidStakeReputation, error)
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

func (r *reputationChange) CalculateLiquidStakeReputationForAddress(address casper.Hash) (entities.LiquidStakeReputation, error) {
	query := `
	SELECT 
	    (SELECT ABS(SUM(amount)) FROM reputation_changes WHERE contract_package_hash = ? and address = ?) as liquid_amount, 
	    (SELECT ABS(SUM(amount))  FROM reputation_changes WHERE contract_package_hash != ? and address = ?) as staked_amount  
	FROM reputation_changes WHERE address = ?;
`

	args := []interface{}{
		r.contractPackageHashes.ReputationContractPackageHash,
		address,
		r.contractPackageHashes.ReputationContractPackageHash,
		address,
		address,
	}

	liquidStakeReputation := entities.LiquidStakeReputation{}
	err := r.conn.Get(&liquidStakeReputation, query, args...)
	if err != nil {
		return entities.LiquidStakeReputation{}, err
	}

	return liquidStakeReputation, nil
}

func (r *reputationChange) CalculateAggregatedLiquidStakeReputationForAddresses(addresses []casper.Hash) ([]entities.LiquidStakeReputation, error) {
	addressesParams := make([]string, 0, len(addresses))
	for range addresses {
		addressesParams = append(addressesParams, "?")
	}

	// liquid_amount
	query := fmt.Sprintf(`
	SELECT 
	    ABS(SUM(amount)) as liquid_amount,
		address
	FROM reputation_changes  WHERE contract_package_hash = ? and address in (%s) GROUP BY address;
`, strings.Join(addressesParams, ","))

	args := []interface{}{
		r.contractPackageHashes.ReputationContractPackageHash,
	}
	for _, address := range addresses {
		args = append(args, address)
	}

	liquidReputations := make([]entities.LiquidStakeReputation, 0)
	if err := r.conn.Select(&liquidReputations, query, args...); err != nil {
		return nil, err
	}
	liquidReputationsMap := make(map[string]*uint64)
	for _, entry := range liquidReputations {
		liquidReputationsMap[entry.Address.String()] = entry.LiquidAmount
	}

	// staked_amount
	query = fmt.Sprintf(`
	SELECT 
	    ABS(SUM(amount)) as staked_amount,
	    address
	FROM reputation_changes  WHERE contract_package_hash != ? and address in (%s) GROUP BY address;
`, strings.Join(addressesParams, ","))

	stakeReputations := make([]entities.LiquidStakeReputation, 0)
	if err := r.conn.Select(&stakeReputations, query, args...); err != nil {
		return nil, err
	}

	for i := range stakeReputations {
		liquidReputation := liquidReputationsMap[stakeReputations[i].Address.String()]
		stakeReputations[i].LiquidAmount = liquidReputation
	}

	return stakeReputations, nil
}
