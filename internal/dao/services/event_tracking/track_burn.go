package event_tracking

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackBurn struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware

	contractPackage types.Hash
}

func NewTrackBurn() *TrackBurn {
	return &TrackBurn{}
}

func (s *TrackBurn) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackBurn) Execute() error {
	burnEvent, err := events.ParseBurn(s.GetCESEvent())
	if err != nil {
		return err
	}

	var address *types.Hash
	if burnEvent.Address.AccountHash != nil {
		address = burnEvent.Address.AccountHash
	} else {
		address = burnEvent.Address.ContractPackageHash
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*address,
			s.contractPackage,
			nil,
			burnEvent.Amount.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonBurn,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
