package event_processing

import (
	"casper-dao-middleware/internal/crdao/services/settings"
	"context"

	"casper-dao-middleware/internal/crdao/dao_event_parser"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/pkg/casper"

	"go.uber.org/zap"
)

// ProcessEventStream command start number of concurrent worker to process event from synchronous event stream
type ProcessEventStream struct {
	di.BaseStreamURLAware
	di.CasperClientAware
	di.EntityManagerAware
	di.DAOContractPackageHashesAware

	daoContractHashes                     map[string]string
	eventStreamPath                       string
	variableRepositoryContractStorageUref string
	nodeStartFromEventID                  uint64
	dictionarySetEventsBuffer             uint32
}

func NewProcessEventStream() *ProcessEventStream {
	return &ProcessEventStream{}
}

func (c *ProcessEventStream) SetNodeStartFromEventID(eventID uint64) *ProcessEventStream {
	c.nodeStartFromEventID = eventID
	return c
}

func (c *ProcessEventStream) SetDAOContractHashes(daoContractHashes map[string]string) *ProcessEventStream {
	c.daoContractHashes = daoContractHashes
	return c
}

func (c *ProcessEventStream) SetDictionarySetEventsBuffer(buffer uint32) *ProcessEventStream {
	c.dictionarySetEventsBuffer = buffer
	return c
}

func (c *ProcessEventStream) SetVariableRepositoryContractStorageUref(uref string) {
	c.variableRepositoryContractStorageUref = uref
}

func (c *ProcessEventStream) SetEventStreamPath(eventPath string) *ProcessEventStream {
	c.eventStreamPath = eventPath
	return c
}

func (c *ProcessEventStream) Execute(ctx context.Context) error {
	eventListener, err := casper.NewEventListener(c.GetBaseStreamURL(), c.eventStreamPath, &c.nodeStartFromEventID)
	if err != nil {
		return err
	}

	daoEventParser, err := dao_event_parser.NewDaoEventParser(c.GetCasperClient(), c.daoContractHashes, c.dictionarySetEventsBuffer)
	if err != nil {
		return err
	}

	syncDaoSetting := settings.NewSyncDAOSettings()
	syncDaoSetting.SetCasperClient(c.GetCasperClient())
	syncDaoSetting.SetVariableRepositoryContractStorageUref(c.variableRepositoryContractStorageUref)
	syncDaoSetting.SetEntityManager(c.GetEntityManager())

	for _, setting := range settings.DaoSettings {
		syncDaoSetting.SetSetting(setting)
		if err := syncDaoSetting.Execute(); err != nil {
			zap.S().With(zap.String("setting", setting)).Info("failed to sync DAO setting")
		}
	}

	processRawDeploy := NewProcessRawDeploy()
	processRawDeploy.SetDAOEventParser(daoEventParser)
	processRawDeploy.SetCasperClient(c.GetCasperClient())
	processRawDeploy.SetEntityManager(c.GetEntityManager())
	processRawDeploy.SetDAOContractPackageHashes(c.GetDAOContractPackageHashes())

	stopListening := func() {
		eventListener.Close()
		zap.S().Info("Finish ProcessEvents command successfully")
	}
	// in case of blocking on eventListener.ReadEvent(), shutdown will happen on next event/ loop iteration
	for {
		select {
		case <-ctx.Done():
			stopListening()
			return nil
		default:
			rawEventData, err := eventListener.ReadEvent()
			if err != nil {
				zap.S().With(zap.Error(err)).Error("Error on event listening")
				stopListening()
				return err
			}

			if rawEventData.EventType != casper.DeployProcessedEventType {
				zap.S().Info("Skip not supported event type, expect DeployProcessedEvent")
				continue
			}

			deployProcessedEvent, err := rawEventData.Data.ParseAsDeployProcessedEvent()
			if err != nil {
				zap.S().With(zap.Error(err)).Info("Failed to parse rawEvent as DeployProcessedEvent")
				return err
			}

			if deployProcessedEvent.DeployProcessed.ExecutionResult.Success == nil {
				zap.S().With(zap.Error(err)).Info("Failed to parse rawEvent as DeployProcessedEvent")
				continue
			}

			processRawDeploy.SetDeployProcessedEvent(deployProcessedEvent)
			if err = processRawDeploy.Execute(); err != nil {
				zap.S().With(zap.Error(err)).Error("Failed to handle DeployProcessedEvent")
			}
		}
	}
}
