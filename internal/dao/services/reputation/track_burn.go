package reputation

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation"
)

type TrackBurn struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackBurn() *TrackBurn {
	return &TrackBurn{}
}

func (s *TrackBurn) Execute() error {
	burnEvent, err := reputation.ParseBurn(s.GetCESEvent())
	if err != nil {
		return err
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	burnedAmount := burnEvent.Amount.Value().Int64()

	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*burnEvent.Address.ToHash(),
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			-burnedAmount,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonBurned,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	if err := s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes); err != nil {
		return err
	}

	liquidStakeReputation, err := s.GetEntityManager().
		ReputationChangeRepository().
		CalculateLiquidStakeReputationForAddress(*burnEvent.Address.ToHash())
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
		*burnEvent.Address.ToHash(),
		nil,
		liquidReputation,
		stakedReputation,
		uint64(burnedAmount),
		0,
		deployProcessedEvent.DeployProcessed.DeployHash,
		entities.ReputationChangeReasonBurned,
		deployProcessedEvent.DeployProcessed.Timestamp)

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch([]entities.TotalReputationSnapshot{reputationTotal})
}
