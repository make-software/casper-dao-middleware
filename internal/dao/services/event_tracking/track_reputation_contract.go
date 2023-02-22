package event_tracking

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/reputation"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackReputationContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware

	contractPackage types.Hash
}

func NewTrackReputationContract(contractPackage types.Hash) *TrackReputationContract {
	return &TrackReputationContract{
		contractPackage: contractPackage,
	}
}

func (s *TrackReputationContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	}

	return nil
}

func (s *TrackReputationContract) trackBurnEvent() error {
	burnEvent, err := reputation.ParseBurn(s.GetCESEvent())
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

func (s *TrackReputationContract) trackMintEvent() error {
	mintEvent, err := reputation.ParseMint(s.GetCESEvent())
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
