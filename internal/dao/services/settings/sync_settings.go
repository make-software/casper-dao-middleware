package settings

import (
	"casper-dao-middleware/internal/dao/di"

	"go.uber.org/zap"
)

type SyncDAOSettings struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string
	settings                              []string
}

func NewSyncDAOSettings() SyncDAOSettings {
	return SyncDAOSettings{}
}

func (c *SyncDAOSettings) SetVariableRepositoryContractStorageUref(uref string) {
	c.variableRepositoryContractStorageUref = uref
}

func (c *SyncDAOSettings) SetSettings(settings []string) {
	c.settings = settings
}

func (c *SyncDAOSettings) Execute() {
	syncDaoSetting := NewSyncDAOSetting()
	syncDaoSetting.SetCasperClient(c.GetCasperClient())
	syncDaoSetting.SetVariableRepositoryContractStorageUref(c.variableRepositoryContractStorageUref)
	syncDaoSetting.SetEntityManager(c.GetEntityManager())

	for _, setting := range c.settings {
		syncDaoSetting.SetSetting(setting)
		if err := syncDaoSetting.Execute(); err != nil {
			zap.S().With(zap.String("setting", setting)).Info("failed to sync DAO setting")
		}
	}
}
