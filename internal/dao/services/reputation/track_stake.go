package reputation

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation"
)

type TrackStake struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackStake() *TrackStake {
	return &TrackStake{}
}

func (s *TrackStake) Execute() error {
	stake, err := reputation.ParseStakeEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			stake.Worker,
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			-stake.Amount.Value().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonStaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	if err := s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes); err != nil {
		return err
	}

	liquidStakeReputation, err := s.GetEntityManager().
		ReputationChangeRepository().
		CalculateLiquidStakeReputationForAddress(stake.Worker)
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
		stake.Worker,
		nil,
		liquidReputation,
		stakedReputation,
		0,
		0,
		deployProcessedEvent.DeployProcessed.DeployHash,
		entities.ReputationChangeReasonStaked,
		deployProcessedEvent.DeployProcessed.Timestamp)

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch([]entities.TotalReputationSnapshot{reputationTotal})
}
