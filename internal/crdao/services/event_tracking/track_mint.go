package event_tracking

import (
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackMint struct {
	di.EntityManagerAware

	deployProcessed casper.DeployProcessed
	contractPackage types.Hash
	cesEvent        ces.Event
}

func NewTrackMint() *TrackMint {
	return &TrackMint{}
}

func (s *TrackMint) SetCESEvent(event ces.Event) {
	s.cesEvent = event
}

func (s *TrackMint) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackMint) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackMint) Execute() error {
	mintEvent, err := events.ParseMint(s.cesEvent)
	if err != nil {
		return err
	}

	var address *types.Hash
	if mintEvent.Address.AccountHash != nil {
		address = mintEvent.Address.AccountHash
	} else {
		address = mintEvent.Address.ContractPackageHash
	}

	changes := []entities.ReputationChange{
		entities.NewReputationChange(
			*address,
			s.contractPackage,
			nil,
			mintEvent.Amount.Into().Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonMint,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
