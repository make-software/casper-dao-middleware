package varaible_repository

import (
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/variable_repository"
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
	case variable_repository.ValueUpdatedEventName:
		trackValueUpdated := NewTrackValueUpdated()
		trackValueUpdated.SetCESEvent(cesEvent)
		trackValueUpdated.SetEntityManager(s.GetEntityManager())
		if err := trackValueUpdated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.VariableRepositoryContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
