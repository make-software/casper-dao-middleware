package varaible_repository

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/variable_repository"
)

type TrackVariableRepositoryContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackVariableRepositoryContract() *TrackVariableRepositoryContract {
	return &TrackVariableRepositoryContract{}
}

func (s *TrackVariableRepositoryContract) Execute() error {
	cesEvent := s.GetCESEvent()

	switch cesEvent.Name {
	case variable_repository.ValueUpdatedEventName:
		return s.trackValueUpdated()
	}

	return nil
}

func (s *TrackVariableRepositoryContract) trackValueUpdated() error {
	valueUpdated, err := variable_repository.ParseValueUpdatedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	var activationTime time.Time

	if valueUpdated.ActivationTime != nil {
		activationTime = time.Unix(int64(*valueUpdated.ActivationTime), 0)
	}

	setting := entities.NewSetting(valueUpdated.Key, valueUpdated.Value.String(), nil, &activationTime)
	return s.GetEntityManager().SettingRepository().Upsert(setting)
}
