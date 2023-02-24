package va_nft

import (
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/va_nft"
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
	case va_nft.TransferEventName:
		trackTransfer := NewTrackTransfer()
		trackTransfer.SetCESEvent(cesEvent)
		trackTransfer.SetEntityManager(s.GetEntityManager())
		if err := trackTransfer.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", doaContractMetadata.VANFTContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
