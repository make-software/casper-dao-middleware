package event_tracking

import (
	"casper-dao-middleware/internal/dao_event_parser/events"
	"casper-dao-middleware/internal/di"
	"casper-dao-middleware/internal/entities"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackMint struct {
	di.EntityManagerAware

	contractPackage types.Hash
	deployProcessed casper.DeployProcessed
	eventBody       []byte
}

func NewTrackMint() *TrackMint {
	return &TrackMint{}
}

func (s *TrackMint) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackMint) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackMint) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackMint) Execute() error {
	mintEvent, err := events.ParseMintEvent(s.eventBody)
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
			(*mintEvent.Amount).Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonMint,
			s.deployProcessed.Timestamp),
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
