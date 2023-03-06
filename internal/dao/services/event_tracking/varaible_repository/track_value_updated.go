package varaible_repository

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/variable_repository"
)

type TrackValueUpdated struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackValueUpdated() *TrackValueUpdated {
	return &TrackValueUpdated{}
}

func (s *TrackValueUpdated) Execute() error {
	valueUpdated, err := variable_repository.ParseValueUpdatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var activationTime *time.Time

	if valueUpdated.ActivationTime != nil {
		newTime := time.Unix(int64(*valueUpdated.ActivationTime), 0)
		activationTime = &newTime
	}

	setting := entities.NewSetting(valueUpdated.Key, valueUpdated.Value.String(), nil, activationTime)
	return s.GetEntityManager().SettingRepository().Upsert(setting)
}
