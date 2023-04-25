package voting

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/types"
	"casper-dao-middleware/internal/dao/utils"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackVotingEnded struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware

	voterContractPackageHash casper_types.Hash
}

func NewTrackVotingEnded() *TrackVotingEnded {
	return &TrackVotingEnded{}
}

func (s *TrackVotingEnded) SetVoterContractPackageHash(hash casper_types.Hash) {
	s.voterContractPackageHash = hash
}

func (s *TrackVotingEnded) Execute() error {
	votingEnded, err := base.ParseVotingEndedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.updateVotingState(votingEnded); err != nil {
		return err
	}

	if err := s.collectReputationChanges(votingEnded, s.voterContractPackageHash); err != nil {
		return err
	}

	if err := s.aggregateReputationTotals(votingEnded); err != nil {
		return err
	}

	return nil
}

func (s *TrackVotingEnded) collectReputationChanges(votingEnded base.VotingEndedEvent, voterContractPackageHash casper_types.Hash) error {
	changes := make([]entities.ReputationChange, 0, len(votingEnded.Burns)+len(votingEnded.Mints)+len(votingEnded.Unstakes)*2)
	deployProcessedEvent := s.GetDeployProcessedEvent()

	// if we have unstakes it means that address will also have one record in mints
	// so here, we need only subtract unstake amount from voter contract and add unstake amount to reputation contract
	for key, val := range votingEnded.Unstakes {
		address, _ := casper_types.NewHashFromHexString(key.Element1)
		changes = append(changes,
			entities.NewReputationChange(
				address,
				voterContractPackageHash,
				&votingEnded.VotingID,
				-val.Into().Int64(),
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),

			entities.NewReputationChange(
				address,
				s.GetDAOContractsMetadata().ReputationContractPackageHash,
				&votingEnded.VotingID,
				val.Into().Int64(),
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
		)
	}

	// in case of mints, just add mint amount to reputation contract
	for key, val := range votingEnded.Mints {
		address, _ := casper_types.NewHashFromHexString(key.Element1)

		changes = append(changes, entities.NewReputationChange(
			address,
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			val.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingGained,
			deployProcessedEvent.DeployProcessed.Timestamp),
		)
	}

	// in case of burns, subtract burn amount from voter contract
	for key, val := range votingEnded.Burns {
		address, _ := casper_types.NewHashFromHexString(key.Element1)

		changes = append(changes, entities.NewReputationChange(
			address,
			voterContractPackageHash,
			nil,
			-val.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingLost,
			deployProcessedEvent.DeployProcessed.Timestamp),
		)
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}

func (s *TrackVotingEnded) aggregateReputationTotals(votingEnded base.VotingEndedEvent) error {
	deployProcessedEvent := s.GetDeployProcessedEvent()

	addresses := make([]casper_types.Hash, 0, len(votingEnded.Mints)+len(votingEnded.Burns))

	if len(votingEnded.Mints) == 0 && len(votingEnded.Burns) == 0 {
		for key := range votingEnded.Unstakes {
			address, _ := casper_types.NewHashFromHexString(key.Element1)
			addresses = append(addresses, address)
		}
	}

	for key := range votingEnded.Mints {
		address, _ := casper_types.NewHashFromHexString(key.Element1)
		addresses = append(addresses, address)
	}

	for key := range votingEnded.Burns {
		address, _ := casper_types.NewHashFromHexString(key.Element1)
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

	totals := make([]entities.TotalReputationSnapshot, 0, len(votingEnded.Mints)+len(votingEnded.Burns))

	// if Mints and Burns are empty, just iterate over Unstakes to create record with released reputation
	if len(votingEnded.Mints) == 0 && len(votingEnded.Burns) == 0 {
		for key := range votingEnded.Unstakes {
			address, _ := casper_types.NewHashFromHexString(key.Element1)

			liquidStakeReputation, ok := addressToLiquidStakeReputation[address.ToHex()]
			if !ok {
				continue
			}

			totals = append(totals, entities.NewTotalReputationSnapshot(
				address,
				&votingEnded.VotingID,
				*liquidStakeReputation.LiquidAmount,
				*liquidStakeReputation.StakedAmount,
				0,
				0,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp))
		}
	}

	for key, val := range votingEnded.Mints {
		address, _ := casper_types.NewHashFromHexString(key.Element1)

		liquidStakeReputation, ok := addressToLiquidStakeReputation[address.ToHex()]
		if !ok {
			continue
		}

		totals = append(totals, entities.NewTotalReputationSnapshot(
			address,
			&votingEnded.VotingID,
			*liquidStakeReputation.LiquidAmount,
			*liquidStakeReputation.StakedAmount,
			0,
			val.Into().Uint64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingGained,
			deployProcessedEvent.DeployProcessed.Timestamp))
	}

	for key, val := range votingEnded.Burns {
		address, _ := casper_types.NewHashFromHexString(key.Element1)

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
			&votingEnded.VotingID,
			liquidReputation,
			stakedReputation,
			val.Into().Uint64(),
			0,
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingLost,
			deployProcessedEvent.DeployProcessed.Timestamp))
	}

	return s.GetEntityManager().TotalReputationSnapshotRepository().SaveBatch(totals)
}

func (s *TrackVotingEnded) updateVotingState(votingEnded base.VotingEndedEvent) error {
	storedVoting, err := s.GetEntityManager().VotingRepository().GetByVotingID(votingEnded.VotingID)
	if err != nil {
		return err
	}

	// we need to calculate FormalVotingStarts/FormalVotingEnds based on the VotingEnded result
	if storedVoting.FormalVotingStartsAt == nil {
		var formalStartsAt, formalEndsAt time.Time
		var allVotesCount = votingEnded.VotesInFavor + votingEnded.VotesAgainst
		var inFavourPercent = utils.PercentOf(votingEnded.VotesInFavor, allVotesCount)

		//This behavior is configured using VotingClearnessDelta Governance Variable.
		//It is a numeric value which tells how far from 50/50 result can be in percent points, before the time will be doubled.
		//For example, when VotingClearnessDelta is set to 8 and the result of the Informal Voting is 42 percent "for" and 58 "against" then the time between votings should be doubled.
		//When the result is 41/59, the default value of time will be used.
		if 50-inFavourPercent > float64(storedVoting.ConfigVotingClearnessDelta) {
			formalStartsAt = storedVoting.InformalVotingEndsAt.Add(time.Millisecond * time.Duration(storedVoting.ConfigTimeBetweenInformalAndFormalVoting))
			formalEndsAt = formalStartsAt.Add(time.Millisecond * time.Duration(storedVoting.FormalVotingTime))
		} else {
			formalStartsAt = storedVoting.InformalVotingEndsAt.Add(time.Millisecond * time.Duration(storedVoting.ConfigTimeBetweenInformalAndFormalVoting*2))
			formalEndsAt = formalStartsAt.Add(time.Millisecond * time.Duration(storedVoting.FormalVotingTime))
		}

		storedVoting.FormalVotingStartsAt = &formalStartsAt
		storedVoting.FormalVotingEndsAt = &formalEndsAt
	}

	if votingEnded.VotingType == types.VotingTypeInformal {
		storedVoting.InformalVotingResult = &votingEnded.VotingResult
	} else {
		storedVoting.FormalVotingResult = &votingEnded.VotingResult
	}

	return s.GetEntityManager().VotingRepository().Update(storedVoting)
}
