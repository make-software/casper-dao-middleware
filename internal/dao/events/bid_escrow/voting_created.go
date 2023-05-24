package bid_escrow

import (
	"errors"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const VotingCreatedEventName = "BidEscrowVotingCreated"

type VotingCreatedEvent struct {
	JobOfferID                               uint32
	BidID                                    uint32
	JobID                                    uint32
	JobPoster                                casper.Hash
	Worker                                   casper.Hash
	Creator                                  types.Address
	VotingID                                 uint32
	ConfigInformalQuorum                     uint32
	ConfigInformalVotingTime                 uint64
	ConfigFormalQuorum                       uint32
	ConfigFormalVotingTime                   uint64
	ConfigTotalOnboarded                     clvalue.UInt512
	ConfigDoubleTimeBetweenVotings           bool
	ConfigVotingClearnessDelta               clvalue.UInt512
	ConfigTimeBetweenInformalAndFormalVoting uint64
}

func ParseVotingCreatedEvent(event ces.Event) (VotingCreatedEvent, error) {
	var votingCreated VotingCreatedEvent
	var err error

	val, ok := event.Data["bid_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid bid_id value in event")
	}
	votingCreated.BidID = val.UI32.Value()

	val, ok = event.Data["job_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid job_id value in event")
	}
	votingCreated.JobID = val.UI32.Value()

	val, ok = event.Data["job_offer_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid job_offer_id value in event")
	}
	votingCreated.JobID = val.UI32.Value()

	val, ok = event.Data["worker"]
	if !ok || val.Type != cltype.Key {
		return VotingCreatedEvent{}, errors.New("invalid worker value in event")
	}

	if val.Key.Account != nil {
		votingCreated.Worker = val.Key.Account.Hash
	} else {
		votingCreated.Worker = *val.Key.Hash
	}

	val, ok = event.Data["job_poster"]
	if !ok || val.Type != cltype.Key {
		return VotingCreatedEvent{}, errors.New("invalid job_poster value in event")
	}
	if val.Key.Account != nil {
		votingCreated.JobPoster = val.Key.Account.Hash
	} else {
		votingCreated.JobPoster = *val.Key.Hash
	}

	val, ok = event.Data["creator"]
	if !ok || val.Type != cltype.Key {
		return VotingCreatedEvent{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return VotingCreatedEvent{}, err
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCreated.VotingID = val.UI32.Value()

	val, ok = event.Data["config_informal_quorum"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid config_informal_quorum value in event")
	}
	votingCreated.ConfigInformalQuorum = val.UI32.Value()

	val, ok = event.Data["config_informal_voting_time"]
	if !ok || val.Type != cltype.UInt64 {
		return VotingCreatedEvent{}, errors.New("invalid config_informal_voting_time value in event")
	}
	votingCreated.ConfigInformalVotingTime = val.UI64.Value()

	val, ok = event.Data["config_formal_quorum"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCreatedEvent{}, errors.New("invalid config_formal_quorum value in event")
	}
	votingCreated.ConfigFormalQuorum = val.UI32.Value()

	val, ok = event.Data["config_formal_voting_time"]
	if !ok || val.Type != cltype.UInt64 {
		return VotingCreatedEvent{}, errors.New("invalid config_formal_voting_time value in event")
	}
	votingCreated.ConfigFormalVotingTime = val.UI64.Value()

	val, ok = event.Data["config_total_onboarded"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingCreatedEvent{}, errors.New("invalid config_total_onboarded value in event")
	}
	votingCreated.ConfigTotalOnboarded = *val.UI512

	val, ok = event.Data["config_double_time_between_votings"]
	if !ok || val.Type != cltype.Bool {
		return VotingCreatedEvent{}, errors.New("invalid config_double_time_between_votings value in event")
	}
	votingCreated.ConfigDoubleTimeBetweenVotings = val.Bool.Value()

	val, ok = event.Data["config_voting_clearness_delta"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingCreatedEvent{}, errors.New("invalid config_voting_clearness_delta value in event")
	}
	votingCreated.ConfigVotingClearnessDelta = *val.UI512

	val, ok = event.Data["config_time_between_informal_and_formal_voting"]
	if !ok || val.Type != cltype.UInt64 {
		return VotingCreatedEvent{}, errors.New("invalid config_time_between_informal_and_formal_voting value in event")
	}
	votingCreated.ConfigTimeBetweenInformalAndFormalVoting = val.UI64.Value()

	return votingCreated, nil
}
