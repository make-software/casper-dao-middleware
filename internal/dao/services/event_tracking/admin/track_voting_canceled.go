package admin

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/admin"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackVotingCanceled struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackVotingCanceled() *TrackVotingCanceled {
	return &TrackVotingCanceled{}
}

func (s *TrackVotingCanceled) Execute() error {
	votingCanceled, err := admin.ParseVotingCanceledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := make([]entities.ReputationChange, 0, len(votingCanceled.Unstakes)*2)
	for key, val := range votingCanceled.Unstakes {
		address, _ := casper_types.NewHashFromHexString(key.Element1)
		unstaked := val.Into().Int64()
		changes = append(changes,
			// reverse operation to BallotCast, one positive reputation change to ReputationContractPackageHash
			// and negative from VoterContractPackageHash
			entities.NewReputationChange(
				address,
				s.GetDAOContractsMetadata().ReputationContractPackageHash,
				&votingCanceled.VotingID,
				unstaked,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
			entities.NewReputationChange(
				address,
				s.GetDAOContractsMetadata().SlashingVoterContractPackageHash,
				&votingCanceled.VotingID,
				-unstaked,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
		)
	}

	if err := s.GetEntityManager().VotingRepository().UpdateIsCanceled(votingCanceled.VotingID, true); err != nil {
		return err
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
