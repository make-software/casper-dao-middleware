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
	TotalReputationSnapshotRepository() repositories.TotalReputationSnapshot
	VoteRepository() repositories.Vote
	VotingRepository() repositories.Voting
	JobOfferRepository() repositories.JobOffer
	BidRepository() repositories.Bid
	JobRepository() repositories.Job
	SettingRepository() repositories.Setting
	AccountRepository() repositories.Account
}

type entityManager struct {
	reputationChangesRepo       repositories.ReputationChange
	totalReputationSnapshotRepo repositories.TotalReputationSnapshot
	voteRepository              repositories.Vote
	votingRepository            repositories.Voting
	settingRepository           repositories.Setting
	accountRepo                 repositories.Account
	jobOfferRepo                repositories.JobOffer
	bidRepo                     repositories.Bid
	jobRepo                     repositories.Job
}

func NewEntityManager(db *sqlx.DB, hashes utils.DAOContractsMetadata) EntityManager {
	return &entityManager{
		reputationChangesRepo:       repositories.NewReputationChange(db, hashes),
		totalReputationSnapshotRepo: repositories.NewTotalReputationSnapshot(db),
		voteRepository:              repositories.NewVote(db),
		votingRepository:            repositories.NewVoting(db),
		settingRepository:           repositories.NewSetting(db),
		accountRepo:                 repositories.NewAccount(db),
		jobOfferRepo:                repositories.NewJobOffer(db),
		bidRepo:                     repositories.NewBid(db),
		jobRepo:                     repositories.NewJob(db),
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

func (e entityManager) TotalReputationSnapshotRepository() repositories.TotalReputationSnapshot {
	return e.totalReputationSnapshotRepo
}

func (e entityManager) JobOfferRepository() repositories.JobOffer {
	return e.jobOfferRepo
}

func (e entityManager) BidRepository() repositories.Bid {
	return e.bidRepo
}

func (e entityManager) JobRepository() repositories.Job {
	return e.jobRepo
}
