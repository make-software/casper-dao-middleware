package settings

import (
	"context"
	"errors"
	"time"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/ces-go-parser"
	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/variable_repository"
)

type SyncInitialDAOSettings struct {
	di.EntityManagerAware
	di.CasperClientAware
	di.DAOContractsMetadataAware
	di.CESParserAware

	variableRepoInstallDeployHash casper.Hash
}

func NewSyncInitialDAOSettings() SyncInitialDAOSettings {
	return SyncInitialDAOSettings{}
}

func (c *SyncInitialDAOSettings) SetVariableRepoInstallDeployHash(hash casper.Hash) {
	c.variableRepoInstallDeployHash = hash
}

func (c *SyncInitialDAOSettings) Execute() error {
	deployResult, err := c.GetCasperClient().GetDeploy(context.Background(), c.variableRepoInstallDeployHash.ToHex())
	if err != nil {
		return err
	}

	if len(deployResult.ExecutionResults) == 0 {
		return errors.New("error: unsuccessful variable repository contract install deploy")
	}

	results, err := c.GetCESParser().ParseExecutionResults(deployResult.ExecutionResults[0].Result)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error != nil {
			zap.S().With(zap.Error(err)).Error("Failed to parse ces events")
			continue
		}

		if result.Event.Name != variable_repository.ValueUpdatedEventName {
			continue
		}

		if err := c.trackValueUpdatedEvent(result.Event); err != nil {
			zap.S().With(zap.Error(err)).Error("Failed to track ValueUpdated event")
		}
	}

	return nil
}

func (c *SyncInitialDAOSettings) trackValueUpdatedEvent(event ces.Event) error {
	valueUpdated, err := variable_repository.ParseValueUpdatedEvent(event)
	if err != nil {
		return err
	}

	var activationTime *time.Time

	if valueUpdated.ActivationTime != nil {
		newTime := time.Unix(int64(*valueUpdated.ActivationTime), 0)
		activationTime = &newTime
	}

	setting := entities.NewSetting(valueUpdated.Key, valueUpdated.Value.String(), nil, activationTime)
	if err := c.GetEntityManager().SettingRepository().Save(setting); err != nil {
		return err
	}

	zap.S().With(zap.String("setting", setting.Name)).Info("ValueUpdated event tracked successfully!")
	return nil
}
