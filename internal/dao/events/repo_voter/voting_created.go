package repo_voter

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const VotingCreatedEventName = "RepoVotingCreated"

type VotingCreatedEvent struct {
	VariableRepoToEdit                       types.Address
	Key                                      string
	Value                                    []byte
	ActivationTime                           *uint64
	Creator                                  types.Address
	Stake                                    *clvalue.UInt512
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

	val, ok := event.Data["variable_repo_to_edit"]
	if !ok || val.Type != cltype.Key {
		return VotingCreatedEvent{}, errors.New("invalid variable_repo_to_edit value in event")
	}
	votingCreated.VariableRepoToEdit, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return VotingCreatedEvent{}, err
	}

	val, ok = event.Data["key"]
	if !ok || val.Type != cltype.String {
		return VotingCreatedEvent{}, errors.New("invalid key key in event")
	}
	votingCreated.Key = val.StringVal.String()

	val, ok = event.Data["value"]
	if !ok {
		return VotingCreatedEvent{}, errors.New("invalid value value in event")
	}

	listClValue := val.List.Elements
	if len(listClValue) == 0 || listClValue[0].Type != cltype.UInt8 {
		return VotingCreatedEvent{}, errors.New("expected List(u8) for value field")
	}

	valueBytes := make([]byte, 0, len(listClValue))
	for _, clValue := range listClValue {
		valueBytes = append(valueBytes, clValue.UI8.Value())
	}

	votingCreated.Value = valueBytes

	val, ok = event.Data["activation_time"]
	if !ok {
		return VotingCreatedEvent{}, errors.New("invalid activation_time value in event")
	}

	if val.Option != nil {
		if val.Option.Inner.Type != cltype.UInt64 {
			return VotingCreatedEvent{}, errors.New("invalid value inside option of `activation_time` value")
		}

		value := val.Option.Inner.UI64.Value()
		votingCreated.ActivationTime = &value
	}

	val, ok = event.Data["creator"]
	if !ok || val.Type != cltype.Key {
		return VotingCreatedEvent{}, errors.New("invalid creator value in event")
	}
	votingCreated.Creator, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return VotingCreatedEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok {
		return VotingCreatedEvent{}, errors.New("invalid stake value in event")
	}

	if val.Option != nil {
		if val.Option.Inner.Type != cltype.UInt512 {
			return VotingCreatedEvent{}, errors.New("invalid value inside option of `stake` value")
		}

		votingCreated.Stake = val.Option.Inner.UI512
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
