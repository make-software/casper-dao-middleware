package event_tracking

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events"
)

type TrackTransfer struct {
	di.EntityManagerAware
	di.DAOContractsMetadataAware
	di.CESEventAware
}

func NewTrackTransfer() *TrackTransfer {
	return &TrackTransfer{}
}

func (s *TrackTransfer) Execute() error {
	var daoMetadata = s.GetDAOContractsMetadata()
	var account entities.Account
	var cesEvent = s.GetCESEvent()

	switch cesEvent.ContractPackageHash.ToHex() {
	case daoMetadata.VANFTContractPackageHash.ToHex():
		event, err := events.ParseVaTransferEvent(cesEvent)
		if err != nil {
			return err
		}
		account = entities.NewAccount(*event.To.AccountHash, false, true, time.Now().UTC())

	case daoMetadata.KycNFTContractPackageHash.ToHex():
		event, err := events.ParseKYCTransferEvent(cesEvent)
		if err != nil {
			return err
		}
		account = entities.NewAccount(*event.To.AccountHash, true, false, time.Now().UTC())
	}

	return s.GetEntityManager().AccountRepository().Upsert(account)
}
