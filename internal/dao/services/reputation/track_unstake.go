package reputation

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation"
)

type TrackUnstake struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackUnstake() *TrackUnstake {
	return &TrackUnstake{}
}

func (s *TrackUnstake) Execute() error {
	unstake, err := reputation.ParseUnstakeEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	worker, err := unstake.Worker.GetHashValue()
	if err != nil {
		return err
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*worker,
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			unstake.Amount.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonUnstaked,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	if err := s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes); err != nil {
		return err
	}

	liquidStakeReputation, err := s.GetEntityManager().
		ReputationChangeRepository().
		CalculateLiquidStakeReputationForAddress(*worker)
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
		*worker,
		nil,
		liquidReputation,
		stakedReputation,
		0,
		0,
		deployProcessedEvent.DeployProcessed.DeployHash,
		entities.ReputationChangeReasonUnstaked,
		deployProcessedEvent.DeployProcessed.Timestamp)

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch([]entities.TotalReputationSnapshot{reputationTotal})
}
