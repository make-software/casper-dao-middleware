package persistence

import (
	"casper-dao-middleware/internal/crdao/repositories"
	"casper-dao-middleware/internal/crdao/types"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../tests/mocks/entity_manager_mock.go -package=mocks -source=./entity_manager.go EntityManager

// EntityManager main persistence interface
type EntityManager interface {
	ReputationChangeRepository() repositories.ReputationChangeRepository
	VoteRepository() repositories.VoteRepository
	VotingRepository() repositories.VotingRepository
	SettingRepository() repositories.SettingRepository
}

type entityManager struct {
	reputationChangesRepo repositories.ReputationChangeRepository
	voteRepository        repositories.VoteRepository
	votingRepository      repositories.VotingRepository
	settingRepository     repositories.SettingRepository
}

func NewEntityManager(db *sqlx.DB, hashes types.DAOContractsMetadata) EntityManager {
	return &entityManager{
		reputationChangesRepo: repositories.NewReputationChange(db, hashes),
		voteRepository:        repositories.NewVote(db),
		votingRepository:      repositories.NewVoting(db),
		settingRepository:     repositories.NewSetting(db),
	}
}

func (e entityManager) ReputationChangeRepository() repositories.ReputationChangeRepository {
	return e.reputationChangesRepo
}

func (e entityManager) VoteRepository() repositories.VoteRepository {
	return e.voteRepository
}

func (e entityManager) VotingRepository() repositories.VotingRepository {
	return e.votingRepository
}

func (e entityManager) SettingRepository() repositories.SettingRepository {
	return e.settingRepository
}
