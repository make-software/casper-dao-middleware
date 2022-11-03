package persistence

import (
	"casper-dao-middleware/internal/dao_event_parser"
	"casper-dao-middleware/internal/repositories"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../tests/mocks/entity_manager_mock.go -package=mocks -source=./entity_manager.go EntityManager

// EntityManager main persistence interface
type EntityManager interface {
	ReputationChangeRepository() repositories.ReputationChangeRepository
	VoteRepository() repositories.VoteRepository
	VotingRepository() repositories.VotingRepository
}

type entityManager struct {
	reputationChangesRepo repositories.ReputationChangeRepository
	voteRepository        repositories.VoteRepository
	votingRepository      repositories.VotingRepository
}

func NewEntityManager(db *sqlx.DB, hashes dao_event_parser.DAOContractPackageHashes) EntityManager {
	return &entityManager{
		reputationChangesRepo: repositories.NewReputationChange(db, hashes),
		voteRepository:        repositories.NewVote(db),
		votingRepository:      repositories.NewVoting(db),
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
