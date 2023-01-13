package event_tracking

import (
	"casper-dao-middleware/internal/crdao/dao_event_parser/events"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/services/settings"

	"go.uber.org/zap"
)

type TrackValueUpdated struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string
	eventBody                             []byte
}

func NewTrackValueUpdated() *TrackValueUpdated {
	return &TrackValueUpdated{}
}

func (s *TrackValueUpdated) SetVariableRepositoryContractStorageUref(uref string) {
	s.variableRepositoryContractStorageUref = uref
}

func (s *TrackValueUpdated) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackValueUpdated) Execute() error {
	valueUpdated, err := events.ParseValueUpdatedEvent(s.eventBody)
	if err != nil {
		return err
	}

	// if no activation time set we can just update current setting
	// https://github.com/make-software/dao-contracts/blob/develop/dao-modules/src/repository.rs#L34
	if valueUpdated.ActivationTime == nil {
		setting := entities.NewSetting(valueUpdated.Key, valueUpdated.Value.String(), nil, nil)
		return s.GetEntityManager().SettingRepository().Upsert(setting)
	}

	// in other case we need to update the setting stored in named keys
	syncDaoSetting := settings.NewSyncDAOSetting()
	syncDaoSetting.SetCasperClient(s.GetCasperClient())
	syncDaoSetting.SetVariableRepositoryContractStorageUref(s.variableRepositoryContractStorageUref)
	syncDaoSetting.SetEntityManager(s.GetEntityManager())
	syncDaoSetting.SetSetting(valueUpdated.Key)
	if err := syncDaoSetting.Execute(); err != nil {
		zap.S().With(zap.String("setting", valueUpdated.Key)).Info("failed to sync DAO setting")
		return err
	}

	return nil
}
