package event_tracking

import (
	"casper-dao-middleware/internal/dao_event_parser/events"
	"casper-dao-middleware/internal/di"
	"casper-dao-middleware/internal/entities"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
)

type TrackVotingCreated struct {
	di.EntityManagerAware

	deployProcessed casper.DeployProcessed
	eventBody       []byte
}

func NewTrackVotingCreated() *TrackVotingCreated {
	return &TrackVotingCreated{}
}

func (s *TrackVotingCreated) SetEventBody(eventBody []byte) {
	s.eventBody = eventBody
}

func (s *TrackVotingCreated) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackVotingCreated) Execute() error {
	votingCreated, err := events.ParseVotingCreatedEvent(s.eventBody)
	if err != nil {
		return err
	}

	var creator types.Hash
	if votingCreated.Creator.AccountHash != nil {
		creator = *votingCreated.Creator.AccountHash
	} else {
		creator = *votingCreated.Creator.ContractPackageHash
	}

	var isFormal bool
	var informalVotingID *uint32
	var votingQuorum = uint64(votingCreated.ConfigInformalVotingQuorum.Int64())
	var votingTime = votingCreated.ConfigInformalVotingTime

	if votingCreated.FormalVotingID != nil {
		isFormal = true
		votingQuorum = uint64(votingCreated.ConfigFormalVotingQuorum.Int64())
		votingTime = votingCreated.ConfigFormalVotingTime

		informalID := uint32((*votingCreated.VotingID).Uint64())
		informalVotingID = &informalID
	}

	votingID := uint32((*votingCreated.VotingID).Uint64())
	voting := entities.NewVoting(
		creator,
		s.deployProcessed.DeployHash,
		votingID,
		informalVotingID,
		votingTime,
		votingQuorum,
		isFormal,
		s.deployProcessed.Timestamp,
	)

	return s.GetEntityManager().VotingRepository().Save(voting)
}
