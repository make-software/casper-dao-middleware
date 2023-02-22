package event_tracking

import (
	"errors"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/kyc_nft"
)

type TrackKycNFTContract struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackKycNFTContract() *TrackKycNFTContract {
	return &TrackKycNFTContract{}
}

func (s *TrackKycNFTContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case kyc_nft.TransferEventName:
		return s.trackTransfer()
	default:
		return errors.New("unsupported contract event")
	}
}

func (s *TrackKycNFTContract) trackTransfer() error {
	event, err := kyc_nft.ParseTransferEvent(s.GetCESEvent())
	if err != nil {
		return err
	}
	account := entities.NewAccount(*event.To.AccountHash, true, false, time.Now().UTC())

	return s.GetEntityManager().AccountRepository().Upsert(account)
}
