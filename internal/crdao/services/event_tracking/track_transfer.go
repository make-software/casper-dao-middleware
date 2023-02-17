package event_tracking

import (
	"time"

	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackTransfer struct {
	di.EntityManagerAware
	di.DAOContractsMetadataAware

	deployProcessed casper.DeployProcessed
	cesEvent        ces.Event
}

func NewTrackTransfer() *TrackTransfer {
	return &TrackTransfer{}
}

func (s *TrackTransfer) SetCESEvent(event ces.Event) {
	s.cesEvent = event
}

func (s *TrackTransfer) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackTransfer) Execute() error {
	var daoMetadata = s.GetDAOContractsMetadata()
	var account entities.Account

	switch s.cesEvent.ContractPackageHash.ToHex() {
	case daoMetadata.VANFTContractPackageHash.ToHex():
		event, err := events.ParseKYCTransferEvent(s.cesEvent)
		if err != nil {
			return err
		}
		account = entities.NewAccount(*event.To.AccountHash, false, true, time.Now().UTC())

	case daoMetadata.KycNFTContractPackageHash.ToHex():
		event, err := events.ParseVaTransferEvent(s.cesEvent)
		if err != nil {
			return err
		}
		account = entities.NewAccount(*event.To.AccountHash, true, false, time.Now().UTC())
	}

	return s.GetEntityManager().AccountRepository().Upsert(account)
}
