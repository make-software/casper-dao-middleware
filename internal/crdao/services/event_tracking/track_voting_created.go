package event_tracking

import (
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/events"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackVotingCreated struct {
	di.EntityManagerAware

	deployProcessed casper.DeployProcessed
	cesEvent        ces.Event
}

func NewTrackVotingCreated() *TrackVotingCreated {
	return &TrackVotingCreated{}
}

func (s *TrackVotingCreated) SetDeployProcessed(deployProcessed casper.DeployProcessed) {
	s.deployProcessed = deployProcessed
}

func (s *TrackVotingCreated) SetCESEvent(event ces.Event) {
	s.cesEvent = event
}

func (s *TrackVotingCreated) Execute() error {
	var (
		voting entities.Voting
		err    error
	)

	switch s.cesEvent.Name {
	case events.SimpleVotingCreatedEventName:
		voting, err = s.newVotingFromSimpleVotingCreated()
		if err != nil {
			return err
		}
		// place for other type of voting
	}

	return s.GetEntityManager().VotingRepository().Save(&voting)
}

func (s *TrackVotingCreated) newVotingFromSimpleVotingCreated() (entities.Voting, error) {
	simpleVotingCreated, err := events.ParseSimpleVotingCreatedEvent(s.cesEvent)
	if err != nil {
		return entities.Voting{}, err
	}

	var creator types.Hash
	if simpleVotingCreated.Creator.AccountHash != nil {
		creator = *simpleVotingCreated.Creator.AccountHash
	} else {
		creator = *simpleVotingCreated.Creator.ContractPackageHash
	}

	var isFormal bool
	var votingQuorum uint32
	var votingTime uint64

	if simpleVotingCreated.ConfigFormalQuorum != 0 {
		isFormal = true
		votingQuorum = simpleVotingCreated.ConfigFormalQuorum
		votingTime = simpleVotingCreated.ConfigFormalVotingTime
	}

	configTotalOnboarded := simpleVotingCreated.ConfigTotalOnboarded

	return entities.NewVoting(
		creator,
		s.deployProcessed.DeployHash,
		simpleVotingCreated.VotingID,
		votingQuorum,
		votingTime,
		isFormal,
		simpleVotingCreated.ConfigDoubleTimeBetweenVotings,
		simpleVotingCreated.DocumentHash,
		configTotalOnboarded.Into().Uint64(),
		simpleVotingCreated.ConfigVotingClearnessDelta.Into().Uint64(),
		simpleVotingCreated.ConfigTimeBetweenInformalAndFormalVoting,
		s.deployProcessed.Timestamp,
	), nil
}
