package event_processing

import (
	"context"

	"github.com/make-software/casper-go-sdk/sse"
	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
)

type SyncInstallDAOContracts struct {
	di.EntityManagerAware
	di.CasperClientAware
	di.DAOContractsMetadataAware
	di.CESParserAware

	daoContractsInstallBlocks []uint64
}

func NewSyncInstallDAOContracts() SyncInstallDAOContracts {
	return SyncInstallDAOContracts{}
}

func (c *SyncInstallDAOContracts) SetDAOContractsInstallBlocks(heights []uint64) {
	c.daoContractsInstallBlocks = heights
}

func (c *SyncInstallDAOContracts) Execute() error {
	for _, blockHeight := range c.daoContractsInstallBlocks {
		blockResult, err := c.GetCasperClient().GetBlockByHeight(context.Background(), blockHeight)
		if err != nil {
			return err
		}

		processRawDeploy := NewProcessRawDeploy()
		processRawDeploy.SetEntityManager(c.GetEntityManager())
		processRawDeploy.SetCESParser(c.GetCESParser())
		processRawDeploy.SetDAOContractsMetadata(c.GetDAOContractsMetadata())

		for _, deployHash := range blockResult.Block.Body.DeployHashes {
			deployResult, err := c.GetCasperClient().GetDeploy(context.Background(), deployHash.ToHex())
			if err != nil {
				return err
			}

			if len(deployResult.ExecutionResults) == 0 {
				continue
			}

			processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
				DeployProcessed: sse.DeployProcessedPayload{
					DeployHash:      deployResult.Deploy.Hash,
					Account:         deployResult.Deploy.Header.Account.String(),
					Timestamp:       deployResult.Deploy.Header.Timestamp.ToTime(),
					ExecutionResult: deployResult.ExecutionResults[0].Result,
				},
			})
			if err = processRawDeploy.Execute(); err != nil {
				zap.S().With(zap.Error(err)).Error("Failed to process RawDeploy")
				return err
			}
		}
	}

	return nil
}
