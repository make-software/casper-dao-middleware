package slashing_voter

import (
	"time"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/types"
	"casper-dao-middleware/internal/dao/utils"
	casper_types "casper-dao-middleware/pkg/casper/types"
)

type TrackVotingEnded struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewTrackVotingEnded() *TrackVotingEnded {
	return &TrackVotingEnded{}
}

func (s *TrackVotingEnded) Execute() error {
	votingEnded, err := slashing_voter.ParseVotingEndedEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

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
			formalStartsAt = storedVoting.InformalVotingEndsAt.Add(time.Second * time.Duration(storedVoting.ConfigTimeBetweenInformalAndFormalVoting))
			formalEndsAt = formalStartsAt.Add(time.Second * time.Duration(storedVoting.FormalVotingTime))
		} else {
			formalStartsAt = storedVoting.InformalVotingEndsAt.Add(time.Second * time.Duration(storedVoting.ConfigTimeBetweenInformalAndFormalVoting*2))
			formalEndsAt = formalStartsAt.Add(time.Second * time.Duration(storedVoting.FormalVotingTime))
		}

		storedVoting.FormalVotingStartsAt = &formalStartsAt
		storedVoting.FormalVotingEndsAt = &formalEndsAt
	}

	deployProcessedEvent := s.GetDeployProcessedEvent()
	changes := make([]entities.ReputationChange, 0, len(votingEnded.Burns)+len(votingEnded.Mints)+len(votingEnded.Unstakes)*2)

	for key, val := range votingEnded.Unstakes {
		address, _ := casper_types.NewHashFromHexString(key.Element1)
		unstaked := val.Into().Int64()
		changes = append(changes,
			// reverse operation to BallotCast, one positive reputation change to ReputationContractPackageHash
			// and negative from VoterContractPackageHash
			entities.NewReputationChange(
				address,
				s.GetDAOContractsMetadata().ReputationContractPackageHash,
				&votingEnded.VotingID,
				unstaked,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
			entities.NewReputationChange(
				address,
				s.GetDAOContractsMetadata().SlashingVoterContractPackageHash,
				&votingEnded.VotingID,
				-unstaked,
				deployProcessedEvent.DeployProcessed.DeployHash,
				entities.ReputationChangeReasonUnstaked,
				deployProcessedEvent.DeployProcessed.Timestamp,
			),
		)
	}

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

	for key, val := range votingEnded.Burns {
		address, _ := casper_types.NewHashFromHexString(key.Element1)

		changes = append(changes, entities.NewReputationChange(
			address,
			s.GetDAOContractsMetadata().ReputationContractPackageHash,
			nil,
			-val.Into().Int64(),
			deployProcessedEvent.DeployProcessed.DeployHash,
			entities.ReputationChangeReasonVotingLost,
			deployProcessedEvent.DeployProcessed.Timestamp),
		)
	}

	if votingEnded.VotingType == types.VotingTypeInformal {
		storedVoting.InformalVotingResult = &votingEnded.VotingResult
	} else {
		storedVoting.FormalVotingResult = &votingEnded.VotingResult
	}

	if err := s.GetEntityManager().VotingRepository().Update(storedVoting); err != nil {
		return err
	}

	return s.GetEntityManager().ReputationChangeRepository().SaveBatch(changes)
}
