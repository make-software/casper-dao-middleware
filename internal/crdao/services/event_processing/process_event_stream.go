package event_processing

import (
	"context"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/crdao/services/settings"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"

	"casper-dao-middleware/internal/crdao/di"
)

// ProcessEventStream command start number of concurrent worker to process event from synchronous event stream
type ProcessEventStream struct {
	di.BaseStreamURLAware
	di.CasperClientAware
	di.EntityManagerAware
	di.DAOContractsMetadataAware

	daoContractHashes         map[string]types.Hash
	eventStreamPath           string
	nodeStartFromEventID      uint64
	dictionarySetEventsBuffer uint32
}

func NewProcessEventStream() *ProcessEventStream {
	return &ProcessEventStream{}
}

func (c *ProcessEventStream) SetNodeStartFromEventID(eventID uint64) *ProcessEventStream {
	c.nodeStartFromEventID = eventID
	return c
}

func (c *ProcessEventStream) SetDAOContractHashes(daoContractHashes map[string]types.Hash) *ProcessEventStream {
	c.daoContractHashes = daoContractHashes
	return c
}

func (c *ProcessEventStream) SetDictionarySetEventsBuffer(buffer uint32) *ProcessEventStream {
	c.dictionarySetEventsBuffer = buffer
	return c
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

	daoMetaData := c.GetDAOContractsMetadata()

	syncDaoSetting := settings.NewSyncDAOSettings()
	syncDaoSetting.SetCasperClient(c.GetCasperClient())
	syncDaoSetting.SetVariableRepositoryContractStorageUref(daoMetaData.VariableRepositoryContractStorageUref)
	syncDaoSetting.SetEntityManager(c.GetEntityManager())
	syncDaoSetting.SetSettings(settings.VariableRepoSettings)
	syncDaoSetting.Execute()

	cesParser, err := ces.NewParser(c.GetCasperClient(), daoMetaData.CESContracts())
	if err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to create CES Parser")
		return err
	}

	processRawDeploy := NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(c.GetEntityManager())
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(daoMetaData)

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
