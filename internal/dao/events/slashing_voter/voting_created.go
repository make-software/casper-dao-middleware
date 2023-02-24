package slashing_voter

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const VotingCreatedEventName = "SlashingVotingCreated"

type VotingCreatedEvent struct {
	AddressToSlash                           types.Address
	SlashRation                              uint32
	Creator                                  types.Address
	Stake                                    *casper_types.U512
	VotingID                                 uint32
	ConfigInformalQuorum                     uint32
	ConfigInformalVotingTime                 uint64
	ConfigFormalQuorum                       uint32
	ConfigFormalVotingTime                   uint64
	ConfigTotalOnboarded                     casper_types.U512
	ConfigDoubleTimeBetweenVotings           bool
	ConfigVotingClearnessDelta               casper_types.U512
	ConfigTimeBetweenInformalAndFormalVoting uint64
}

func ParseVotingCreatedEvent(event ces.Event) (VotingCreatedEvent, error) {
	var votingCreated VotingCreatedEvent
	var err error

	val, ok := event.Data["address_to_slash"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return VotingCreatedEvent{}, errors.New("invalid address_to_slash value in event")
	}
	votingCreated.AddressToSlash, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return VotingCreatedEvent{}, err
	}

	val, ok = event.Data["slash_ratio"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingCreatedEvent{}, errors.New("invalid slash_ratio value in event")
	}
	votingCreated.SlashRation = *val.U32

	val, ok = event.Data["creator"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return VotingCreatedEvent{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return VotingCreatedEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return VotingCreatedEvent{}, errors.New("invalid stake value in event")
	}

	if val.Option != nil {
		if val.Option.Type.CLTypeID != casper_types.CLTypeIDU512 {
			return VotingCreatedEvent{}, errors.New("invalid value inside option of `stake` value")
		}

		votingCreated.Stake = val.U512
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingCreatedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCreated.VotingID = *val.U32

	val, ok = event.Data["config_informal_quorum"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingCreatedEvent{}, errors.New("invalid config_informal_quorum value in event")
	}
	votingCreated.ConfigInformalQuorum = *val.U32

	val, ok = event.Data["config_informal_voting_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return VotingCreatedEvent{}, errors.New("invalid config_informal_voting_time value in event")
	}
	votingCreated.ConfigInformalVotingTime = *val.U64

	val, ok = event.Data["config_formal_quorum"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingCreatedEvent{}, errors.New("invalid config_formal_quorum value in event")
	}
	votingCreated.ConfigFormalQuorum = *val.U32

	val, ok = event.Data["config_formal_voting_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return VotingCreatedEvent{}, errors.New("invalid config_formal_voting_time value in event")
	}
	votingCreated.ConfigFormalVotingTime = *val.U64

	val, ok = event.Data["config_total_onboarded"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingCreatedEvent{}, errors.New("invalid config_total_onboarded value in event")
	}
	votingCreated.ConfigTotalOnboarded = *val.U512

	val, ok = event.Data["config_double_time_between_votings"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDBool {
		return VotingCreatedEvent{}, errors.New("invalid config_double_time_between_votings value in event")
	}
	votingCreated.ConfigDoubleTimeBetweenVotings = *val.Bool

	val, ok = event.Data["config_voting_clearness_delta"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingCreatedEvent{}, errors.New("invalid config_voting_clearness_delta value in event")
	}
	votingCreated.ConfigVotingClearnessDelta = *val.U512

	val, ok = event.Data["config_time_between_informal_and_formal_voting"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return VotingCreatedEvent{}, errors.New("invalid config_time_between_informal_and_formal_voting value in event")
	}
	votingCreated.ConfigTimeBetweenInformalAndFormalVoting = *val.U64

	return votingCreated, nil
}
