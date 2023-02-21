package event_tracking

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackMint struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	contractPackage types.Hash
}

func NewTrackMint() *TrackMint {
	return &TrackMint{}
}

func (s *TrackMint) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackMint) Execute() error {
	mintEvent, err := events.ParseMint(s.GetCESEvent())
	if err != nil {
		return err
	}

	var address *types.Hash
	if mintEvent.Address.AccountHash != nil {
		address = mintEvent.Address.AccountHash
	} else {
		address = mintEvent.Address.ContractPackageHash
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*address,
			s.contractPackage,
			nil,
			mintEvent.Amount.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonMint,
			deployProcessedEvent.DeployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
