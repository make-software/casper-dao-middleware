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
	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*burnEvent.Address.ToHash(),
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			burnEvent.Amount.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonBurn,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
