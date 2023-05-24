package settings

import (
	"casper-dao-middleware/internal/dao/di"
)

type TrackUpdatedSetting struct {
	di.EntityManagerAware
	di.CESEventAware
}

func NewTrackUpdatedSetting() *TrackUpdatedSetting {
	return &TrackUpdatedSetting{}
}

func (s *TrackUpdatedSetting) Execute() error {
	//valueUpdated, err := variable_repository.ParseValueUpdatedEvent(s.GetCESEvent())
	//if err != nil {
	//	return err
	//}

	//var activationTime *time.Time

	//if valueUpdated.ActivationTime != nil {
	//	newTime := time.Unix(int64(*valueUpdated.ActivationTime), 0)
	//	activationTime = &newTime
	//}

	//setting := entities.NewSetting(valueUpdated.Key, valueUpdated.Value.String(), nil, activationTime)
	//return s.GetEntityManager().SettingRepository().Upsert(setting)
	return nil
}
