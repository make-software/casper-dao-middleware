package events

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const ReputationVotingCreatedEventName = "ReputationVotingCreated"

type ReputationVotingCreatedEvent struct {
	Account                                  types.Address
	Creator                                  types.Address
	DocumentHash                             string
	Stake                                    casper_types.U512
	Action                                   uint8
	Amount                                   casper_types.U512
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

func ParseReputationVotingCreatedEvent(event ces.Event) (ReputationVotingCreatedEvent, error) {
	var votingCreated ReputationVotingCreatedEvent

	val, ok := event.Data["account"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return ReputationVotingCreatedEvent{}, errors.New("invalid account value in event")
	}
	votingCreated.Account = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["creator"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return ReputationVotingCreatedEvent{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator = types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["document_hash"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDString {
		return ReputationVotingCreatedEvent{}, errors.New("invalid document_hash value in event")
	}
	votingCreated.DocumentHash = *val.String

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid stake value in event")
	}
	votingCreated.Stake = *val.U512

	val, ok = event.Data["amount"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid amount value in event")
	}
	votingCreated.Amount = *val.U512

	val, ok = event.Data["action"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid action value in event")
	}
	votingCreated.Action = *val.U8

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCreated.VotingID = *val.U32

	val, ok = event.Data["config_informal_quorum"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_informal_quorum value in event")
	}
	votingCreated.ConfigInformalQuorum = *val.U32

	val, ok = event.Data["config_informal_voting_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_informal_voting_time value in event")
	}
	votingCreated.ConfigInformalVotingTime = *val.U64

	val, ok = event.Data["config_formal_quorum"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_formal_quorum value in event")
	}
	votingCreated.ConfigFormalQuorum = *val.U32

	val, ok = event.Data["config_formal_voting_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_formal_voting_time value in event")
	}
	votingCreated.ConfigFormalVotingTime = *val.U64

	val, ok = event.Data["config_total_onboarded"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_total_onboarded value in event")
	}
	votingCreated.ConfigTotalOnboarded = *val.U512

	val, ok = event.Data["config_double_time_between_votings"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDBool {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_double_time_between_votings value in event")
	}
	votingCreated.ConfigDoubleTimeBetweenVotings = *val.Bool

	val, ok = event.Data["config_voting_clearness_delta"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_voting_clearness_delta value in event")
	}
	votingCreated.ConfigVotingClearnessDelta = *val.U512

	val, ok = event.Data["config_time_between_informal_and_formal_voting"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU64 {
		return ReputationVotingCreatedEvent{}, errors.New("invalid config_time_between_informal_and_formal_voting value in event")
	}
	votingCreated.ConfigTimeBetweenInformalAndFormalVoting = *val.U64

	return votingCreated, nil
}
