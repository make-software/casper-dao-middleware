package reputation

import (
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/reputation"
)

type TrackContract struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackContract() *TrackContract {
	return &TrackContract{}
}

func (s *TrackContract) Execute() error {
	cesEvent := s.GetCESEvent()
	doaContractMetadata := s.GetDAOContractsMetadata()

	switch cesEvent.Name {
	case reputation.MintEventName:
		trackMint := NewTrackMint()
		trackMint.SetCESEvent(cesEvent)
		trackMint.SetEntityManager(s.GetEntityManager())
		trackMint.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackMint.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackMint.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationContractHash.String())).Info("failed to track event")
			return err
		}
	case reputation.BurnEventName:
		trackBurn := NewTrackBurn()
		trackBurn.SetCESEvent(cesEvent)
		trackBurn.SetEntityManager(s.GetEntityManager())
		trackBurn.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBurn.SetDAOContractsMetadata(doaContractMetadata)
		if err := trackBurn.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.ReputationContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
