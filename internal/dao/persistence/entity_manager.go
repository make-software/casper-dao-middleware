package persistence

import (
	"github.com/jmoiron/sqlx"

	"casper-dao-middleware/internal/dao/repositories"
	"casper-dao-middleware/internal/dao/utils"
)

//go:generate mockgen -destination=../tests/mocks/entity_manager_mock.go -package=mocks -source=./entity_manager.go EntityManager

// EntityManager main persistence interface
type EntityManager interface {
	ReputationChangeRepository() repositories.ReputationChange
	ReputationTotalRepository() repositories.ReputationTotal
	VoteRepository() repositories.Vote
	VotingRepository() repositories.Voting
	SettingRepository() repositories.Setting
	AccountRepository() repositories.Account
}

type entityManager struct {
	reputationChangesRepo repositories.ReputationChange
	reputationTotalRepo   repositories.ReputationTotal
	voteRepository        repositories.Vote
	votingRepository      repositories.Voting
	settingRepository     repositories.Setting
	accountRepo           repositories.Account
}

func NewEntityManager(db *sqlx.DB, hashes utils.DAOContractsMetadata) EntityManager {
	return &entityManager{
		reputationChangesRepo: repositories.NewReputationChange(db, hashes),
		reputationTotalRepo:   repositories.NewReputationTotal(db),
		voteRepository:        repositories.NewVote(db),
		votingRepository:      repositories.NewVoting(db),
		settingRepository:     repositories.NewSetting(db),
		accountRepo:           repositories.NewAccount(db),
	}
}

func (e entityManager) ReputationChangeRepository() repositories.ReputationChange {
	return e.reputationChangesRepo
}

func (e entityManager) VoteRepository() repositories.Vote {
	return e.voteRepository
}

func (e entityManager) VotingRepository() repositories.Voting {
	return e.votingRepository
}

func (e entityManager) SettingRepository() repositories.Setting {
	return e.settingRepository
}

func (e entityManager) AccountRepository() repositories.Account {
	return e.accountRepo
}

func (e entityManager) ReputationTotalRepository() repositories.ReputationTotal {
	return e.reputationTotalRepo
}
