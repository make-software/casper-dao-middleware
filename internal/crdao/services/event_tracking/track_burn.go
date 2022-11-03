package event_tracking

import (
	"casper-dao-middleware/internal/crdao/dao_event_parser/events"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackBurn struct {
	di.EntityManagerAware

	contractPackage types.Hash
	deployProcessed casper.DeployProcessed
	eventBody       []byte
}

func NewTrackBurn() *TrackBurn {
	return &TrackBurn{}
}

func (s *TrackBurn) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackBurn) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackBurn) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackBurn) Execute() error {
	burnEvent, err := events.ParseBurnEvent(s.eventBody)
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
			(*burnEvent.Amount).Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonBurn,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
