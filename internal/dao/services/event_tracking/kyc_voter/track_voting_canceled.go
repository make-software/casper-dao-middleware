package kyc_voter

import (
	base_events "casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/services/event_tracking/base"
)

type TrackVotingCanceled struct {
	base.TrackVotingCanceled
}

func NewTrackVotingCanceled() *TrackVotingCanceled {
	return &TrackVotingCanceled{}
}

func (s *TrackVotingCanceled) Execute() error {
	votingCanceled, err := base_events.ParseVotingCanceledEvent(s.GetCESEvent())
	if err != nil {
		return err
	}

	if err := s.UpdateVotingIsCancel(votingCanceled); err != nil {
		return err
	}

	if err := s.CollectReputationChanges(votingCanceled, s.GetDAOContractsMetadata().KycVoterContractPackageHash); err != nil {
		return err
	}

	if err := s.AggregateReputationTotals(votingCanceled); err != nil {
		return err
	}

	return nil
}
