package voting

import (
	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/base"
)

type TrackVotingCanceled struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware

	voterContractPackageHash casper.ContractPackageHash
}

func NewTrackVotingCanceled() *TrackVotingCanceled {
	return &TrackVotingCanceled{}
}

func (s *TrackVotingCanceled) SetVoterContractPackageHash(hash casper.ContractPackageHash) {
	s.voterContractPackageHash = hash
}

func (s *TrackVotingCanceled) Execute() error {
	votingCanceled, err := base.ParseVotingCanceledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.updateVotingIsCancel(votingCanceled); err != nil {
		return err
	}

	if err := s.collectReputationChanges(votingCanceled, s.voterContractPackageHash); err != nil {
		return err
	}

	return nil
}

func (s *TrackVotingCanceled) collectReputationChanges(votingCanceled base.VotingCanceledEvent, voterContractPackageHash casper.ContractPackageHash) error {
	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := make([]entities.ReputationChange, 0, len(votingCanceled.Unstakes)*2)

	for key, val := range votingCanceled.Unstakes {
		address, _ := casper.NewHash(key.Element1)
		unstaked := val.Value().Int64()
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
				voterContractPackageHash,
				&votingCanceled.VotingID,
				-unstaked,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
		)
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}

func (s *TrackVotingCanceled) updateVotingIsCancel(votingCanceled base.VotingCanceledEvent) error {
	return s.GetEntityManager().VotingRepository().UpdateIsCanceled(votingCanceled.VotingID, true)
}

func (s *TrackVotingCanceled) aggregateReputationTotals(votingCanceled base.VotingCanceledEvent) error {
	deployProcessedEvent := s.GetDeployProcessedEvent()

	addresses := make([]casper.Hash, 0, len(votingCanceled.Unstakes))
	for key := range votingCanceled.Unstakes {
		address, _ := casper.NewHash(key.Element1)
		addresses = append(addresses, address)
	}

	aggregatedLiquidStakeReputation, err :=
		s.GetEntityManager().
			ReputationChangeRepository().
			CalculateAggregatedLiquidStakeReputationForAddresses(addresses)

	if err != nil {
		return err
	}

	addressToLiquidStakeReputation := make(map[string]entities.LiquidStakeReputation)
	for _, entry := range aggregatedLiquidStakeReputation {
		addressToLiquidStakeReputation[entry.Address.String()] = entry
	}

	totals := make([]entities.TotalReputationSnapshot, 0, len(votingCanceled.Unstakes))

	for key, val := range votingCanceled.Unstakes {
		address, _ := casper.NewHash(key.Element1)

		liquidStakeReputation, ok := addressToLiquidStakeReputation[address.ToHex()]
		if !ok {
			continue
		}

		var liquidReputation uint64
		if liquidStakeReputation.LiquidAmount != nil {
			liquidReputation = *liquidStakeReputation.LiquidAmount
		}

		var stakedReputation uint64
		if liquidStakeReputation.StakedAmount != nil {
			stakedReputation = *liquidStakeReputation.StakedAmount
		}

		totals = append(totals, entities.NewTotalReputationSnapshot(
			address,
			&votingCanceled.VotingID,
			liquidReputation,
			stakedReputation,
			0,
			val.Value().Uint64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonUnstaked,
			deployProcessedEvent.DeployProcessed.Timestamp))
	}

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch(totals)
}
