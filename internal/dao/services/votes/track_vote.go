package votes

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackVote struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware

	voterContractPackageHash casper_types.Hash
}

func NewTrackVote() *TrackVote {
	return &TrackVote{}
}

func (s *TrackVote) SetVoterContractPackageHash(hash casper_types.Hash) {
	s.voterContractPackageHash = hash
}

func (s *TrackVote) Execute() error {
	ballotCast, err := base.ParseBallotCastEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.saveVote(ballotCast); err != nil {
		return err
	}

	if err := s.collectReputationChanges(ballotCast, s.voterContractPackageHash); err != nil {
		return err
	}

	if err := s.aggregateReputationTotals(ballotCast); err != nil {
		return err
	}

	return nil
}

func (s *TrackVote) saveVote(ballotCast base.BallotCastEvent) error {
	staked := ballotCast.Stake.Into().Int64()

	var isInFavor bool
	if ballotCast.Choice == types.ChoiceInFavor {
		isInFavor = true
	}

	votingID := ballotCast.VotingID
	voting, err := s.GetEntityManager().VotingRepository().GetByVotingID(votingID)
	if err != nil {
		return err
	}

	var isFormal bool
	if voting.FormalVotingStartsAt != nil {
		if time.Now().After(*voting.FormalVotingStartsAt) {
			isFormal = true
		}
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	vote := entities.NewVote(
		*ballotCast.Voter.ToHash(),
		deployProcessedEvent.DeployProcessed.DeployHash,
		ballotCast.VotingID,
		uint64(staked),
		isInFavor,
		isFormal,
		deployProcessedEvent.DeployProcessed.Timestamp)

	return s.GetEntityManager().VoteRepository().Save(vote)
}

func (s *TrackVote) collectReputationChanges(ballotCast base.BallotCastEvent, voterContractPackageHash casper_types.Hash) error {
	deployProcessedEvent := s.GetDeployProcessedEvent()
	staked := ballotCast.Stake.Into().Int64()

	changes := []entities.ReputationChange{
		// one event represent negative reputation leaving from "Reputation" contract
		entities.NewReputationChange(
			*ballotCast.Voter.ToHash(),
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			&ballotCast.VotingID,
			-staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonStaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
		// second event represent positive reputation coming to "Voting" contract
		entities.NewReputationChange(
			*ballotCast.Voter.ToHash(),
			voterContractPackageHash,
			&ballotCast.VotingID,
			staked,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonStaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}

func (s *TrackVote) aggregateReputationTotals(ballotCast base.BallotCastEvent) error {
	deployProcessedEvent := s.GetDeployProcessedEvent()

	liquidStakeReputation, err := s.GetEntityManager().
		ReputationChangeRepository().
		CalculateLiquidStakeReputationForAddress(*ballotCast.Voter.ToHash())
	if err != nil {
		return err
	}

	var liquidReputation uint64
	if liquidStakeReputation.LiquidAmount != nil {
		liquidReputation = *liquidStakeReputation.LiquidAmount
	}

	var stakedReputation uint64
	if liquidStakeReputation.StakedAmount != nil {
		stakedReputation = *liquidStakeReputation.StakedAmount
	}

	reputationTotal := entities.NewTotalReputationSnapshot(
		*ballotCast.Voter.ToHash(),
		&ballotCast.VotingID,
		liquidReputation,
		stakedReputation,
		0,
		0,
		deployProcessedEvent.DeployProcessed.DeployHash,
		entities.ReputationChangeReasonStaked,
		deployProcessedEvent.DeployProcessed.Timestamp)

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch([]entities.TotalReputationSnapshot{reputationTotal})
}
