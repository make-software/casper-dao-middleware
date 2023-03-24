package account

import (
	"errors"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/kyc_nft"
)

type TrackKycTransfer struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackKycTransfer() *TrackKycTransfer {
	return &TrackKycTransfer{}
}

func (s *TrackKycTransfer) Execute() error {
	event, err := kyc_nft.ParseTransferEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if event.To == nil {
		return errors.New("expected not nil transfer receiver")
	}

	account := entities.NewAccount(*event.To.ToHash(), true, false, time.Now().UTC())

	return s.GetEntityManager().AccountRepository().UpsertIsKYC(account)
}
