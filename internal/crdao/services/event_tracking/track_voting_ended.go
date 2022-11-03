package event_tracking

import (
	"casper-dao-middleware/internal/crdao/dao_event_parser/events"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"go.uber.org/zap"
)

type TrackVotingEnded struct {
	di.EntityManagerAware

	contractPackage types.Hash
	deployProcessed casper.DeployProcessed
	eventBody       []byte
}

func NewTrackVotingEnded() *TrackVotingEnded {
	return &TrackVotingEnded{}
}

func (s *TrackVotingEnded) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackVotingEnded) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackVotingEnded) SetEventContractPackage(contractPackage types.Hash) {
	s.contractPackage = contractPackage
}

func (s *TrackVotingEnded) Execute() error {
	votingEnded, err := events.ParseVotingEndedEvent(s.eventBody)
	if err != nil {
		return err
	}

	changes := make([]entities.ReputationChange, 0, len(votingEnded.Burns)+len(votingEnded.Transfers))

	for addressHex, amount := range votingEnded.Burns {
		address, _ := types.NewHashFromHexString(addressHex)

		changes = append(changes, entities.NewReputationChange(
			address,
			s.contractPackage,
			nil,
			-(*amount).Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingBurn,
			s.deployProcessed.Timestamp),
		)
	}

	votingID := uint32((*votingEnded.VotingID).Uint64())
	for addressHex, amount := range votingEnded.Transfers {
		address, _ := types.NewHashFromHexString(addressHex)

		changes = append(changes, entities.NewReputationChange(
			address,
			s.contractPackage,
			&votingID,
			(*amount).Int64(),
			s.deployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingDistribution,
			s.deployProcessed.Timestamp),
		)
	}

	if len(changes) == 0 {
		zap.S().Info("No changes in `voting_created` event")
		return nil
	}

	if err := s.GetEntityManager().VotingRepository().UpdateHasEnded(votingID, true); err != nil {
		return err
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
