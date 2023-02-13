package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const SimpleVotingCreatedEventName = "SimpleVotingCreated"

type SimpleVotingCreated struct {
	Creator                                  dao_types.Address
	DocumentHash                             string
	Stake                                    types.U512
	VotingID                                 uint32
	ConfigInformalQuorum                     uint32
	ConfigInformalVotingTime                 uint64
	ConfigFormalQuorum                       uint32
	ConfigFormalVotingTime                   uint64
	ConfigTotalOnboarded                     types.U512
	ConfigDoubleTimeBetweenVotings           bool
	ConfigVotingClearnessDelta               types.U512
	ConfigTimeBetweenInformalAndFormalVoting uint64
}

func ParseSimpleVotingCreatedEvent(event ces.Event) (SimpleVotingCreated, error) {
	var votingCreated SimpleVotingCreated

	val, ok := event.Data["creator"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return SimpleVotingCreated{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["document_hash"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDString {
		return SimpleVotingCreated{}, errors.New("invalid document_hash value in event")
	}
	votingCreated.DocumentHash = *val.String

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return SimpleVotingCreated{}, errors.New("invalid stake value in event")
	}
	votingCreated.Stake = *val.U512

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return SimpleVotingCreated{}, errors.New("invalid voting_id value in event")
	}
	votingCreated.VotingID = *val.U32

	val, ok = event.Data["config_informal_quorum"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return SimpleVotingCreated{}, errors.New("invalid config_informal_quorum value in event")
	}
	votingCreated.ConfigInformalQuorum = *val.U32

	val, ok = event.Data["config_informal_voting_time"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return SimpleVotingCreated{}, errors.New("invalid config_informal_voting_time value in event")
	}
	votingCreated.ConfigInformalVotingTime = *val.U64

	val, ok = event.Data["config_formal_quorum"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return SimpleVotingCreated{}, errors.New("invalid config_formal_quorum value in event")
	}
	votingCreated.ConfigFormalQuorum = *val.U32

	val, ok = event.Data["config_formal_voting_time"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return SimpleVotingCreated{}, errors.New("invalid config_formal_voting_time value in event")
	}
	votingCreated.ConfigFormalVotingTime = *val.U64

	val, ok = event.Data["config_total_onboarded"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return SimpleVotingCreated{}, errors.New("invalid config_total_onboarded value in event")
	}
	votingCreated.ConfigTotalOnboarded = *val.U512

	val, ok = event.Data["config_double_time_between_votings"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDBool {
		return SimpleVotingCreated{}, errors.New("invalid config_double_time_between_votings value in event")
	}
	votingCreated.ConfigDoubleTimeBetweenVotings = *val.Bool

	val, ok = event.Data["config_voting_clearness_delta"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return SimpleVotingCreated{}, errors.New("invalid config_voting_clearness_delta value in event")
	}
	votingCreated.ConfigVotingClearnessDelta = *val.U512

	val, ok = event.Data["config_time_between_informal_and_formal_voting"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return SimpleVotingCreated{}, errors.New("invalid config_time_between_informal_and_formal_voting value in event")
	}
	votingCreated.ConfigTimeBetweenInformalAndFormalVoting = *val.U64

	return votingCreated, nil
}
