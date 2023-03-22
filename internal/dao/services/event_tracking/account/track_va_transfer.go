package account

import (
	"errors"
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/va_nft"
)

type TrackVATransfer struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackVATransfer() *TrackVATransfer {
	return &TrackVATransfer{}
}

func (s *TrackVATransfer) Execute() error {
	event, err := va_nft.ParseTransferEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if event.To == nil {
		return errors.New("expected not nil transfer receiver")
	}

	account := entities.NewAccount(*event.To.ToHash(), false, true, time.Now().UTC())

	return s.GetEntityManager().AccountRepository().UpsertIsVA(account)
}
