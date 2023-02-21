package event_tracking

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
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
	// TODO: implement tracking
	return nil
}
