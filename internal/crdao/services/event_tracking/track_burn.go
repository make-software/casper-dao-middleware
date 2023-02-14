package event_tracking

import (
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackBurn struct {
	di.EntityManagerAware

	contractPackage types.Hash
	deployProcessed casper.DeployProcessed
	cesEvent        ces.Event
}

func NewTrackBurn() *TrackBurn {
	return &TrackBurn{}
}

func (s *TrackBurn) SetCESEvent(event ces.Event) {
	s.cesEvent = event
}

func (s *TrackBurn) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackBurn) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackBurn) Execute() error {
	burnEvent, err := events.ParseBurn(s.cesEvent)
	if err != nil {
		return err
	}

	var address *types.Hash
	if burnEvent.Address.AccountHash != nil {
		address = burnEvent.Address.AccountHash
	} else {
		address = burnEvent.Address.ContractPackageHash
	}

	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*address,
			s.contractPackage,
			nil,
			burnEvent.Amount.Into().Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonBurn,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
