package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const KYCVotingCreated = "KycVotingCreated"

type KYCVotingCreatedEvent struct {
	SubjectAddress                           dao_types.Address
	DocumentHash                             string
	Creator                                  dao_types.Address
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

func ParseKYCVotingCreatedEvent(event ces.Event) (KYCVotingCreatedEvent, error) {
	var votingCreated KYCVotingCreatedEvent

	val, ok := event.Data["subject_address"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return KYCVotingCreatedEvent{}, errors.New("invalid subject_address value in event")
	}
	votingCreated.SubjectAddress = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["creator"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return KYCVotingCreatedEvent{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return KYCVotingCreatedEvent{}, errors.New("invalid stake value in event")
	}
	votingCreated.Stake = *val.U512

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return KYCVotingCreatedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCreated.VotingID = *val.U32

	val, ok = event.Data["config_informal_quorum"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_informal_quorum value in event")
	}
	votingCreated.ConfigInformalQuorum = *val.U32

	val, ok = event.Data["config_informal_voting_time"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_informal_voting_time value in event")
	}
	votingCreated.ConfigInformalVotingTime = *val.U64

	val, ok = event.Data["config_formal_quorum"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_formal_quorum value in event")
	}
	votingCreated.ConfigFormalQuorum = *val.U32

	val, ok = event.Data["config_formal_voting_time"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_formal_voting_time value in event")
	}
	votingCreated.ConfigFormalVotingTime = *val.U64

	val, ok = event.Data["config_total_onboarded"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_total_onboarded value in event")
	}
	votingCreated.ConfigTotalOnboarded = *val.U512

	val, ok = event.Data["config_double_time_between_votings"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDBool {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_double_time_between_votings value in event")
	}
	votingCreated.ConfigDoubleTimeBetweenVotings = *val.Bool

	val, ok = event.Data["config_voting_clearness_delta"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_voting_clearness_delta value in event")
	}
	votingCreated.ConfigVotingClearnessDelta = *val.U512

	val, ok = event.Data["config_time_between_informal_and_formal_voting"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU64 {
		return KYCVotingCreatedEvent{}, errors.New("invalid config_time_between_informal_and_formal_voting value in event")
	}
	votingCreated.ConfigTimeBetweenInformalAndFormalVoting = *val.U64

	return votingCreated, nil
}
