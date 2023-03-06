package onboarding_request

import (
	"errors"
	"fmt"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const VotingEndedEventName = "VotingEnded"

type VotingEndedEvent struct {
	VotingID             uint32
	VotingType           types.VotingType
	VotingResult         uint8
	StakeInFavour        casper_types.U512
	StakeAgainst         casper_types.U512
	UnboundStakeInFavour casper_types.U512
	UnboundStakeAgainst  casper_types.U512
	VotesInFavor         uint32
	VotesAgainst         uint32
	Unstakes             map[types.Tuple2]casper_types.U512
	Stakes               map[types.Tuple2]casper_types.U512
	Burns                map[types.Tuple2]casper_types.U512
	Mints                map[types.Tuple2]casper_types.U512
}

func ParseVotingEndedEvent(event ces.Event) (VotingEndedEvent, error) {
	var votingEnded VotingEndedEvent
	var err error

	val, ok := event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingEndedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingEnded.VotingID = *val.U32

	val, ok = event.Data["voting_type"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return VotingEndedEvent{}, errors.New("invalid voting_type value in event")
	}

	votingEnded.VotingType, err = types.NewVotingTypeFromByte(*val.U8)
	if err != nil {
		return VotingEndedEvent{}, err
	}

	val, ok = event.Data["voting_result"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return VotingEndedEvent{}, errors.New("invalid voting_result value in event")
	}
	votingEnded.VotingResult = *val.U8

	val, ok = event.Data["stake_in_favor"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingEndedEvent{}, errors.New("invalid stake_in_favor value in event")
	}
	votingEnded.StakeInFavour = *val.U512

	val, ok = event.Data["stake_against"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingEndedEvent{}, errors.New("invalid stake_against value in event")
	}
	votingEnded.StakeAgainst = *val.U512

	val, ok = event.Data["unbound_stake_in_favor"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingEndedEvent{}, errors.New("invalid unbound_stake_in_favor value in event")
	}
	votingEnded.UnboundStakeInFavour = *val.U512

	val, ok = event.Data["unbound_stake_against"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return VotingEndedEvent{}, errors.New("invalid unbound_stake_against value in event")
	}
	votingEnded.UnboundStakeAgainst = *val.U512

	val, ok = event.Data["votes_in_favor"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingEndedEvent{}, errors.New("invalid votes_in_favor value in event")
	}
	votingEnded.VotesInFavor = *val.U32

	val, ok = event.Data["votes_against"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingEndedEvent{}, errors.New("invalid votes_against value in event")
	}
	votingEnded.VotesInFavor = *val.U32

	val, ok = event.Data["unstakes"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDMap {
		return VotingEndedEvent{}, errors.New("invalid unstakes value in event")
	}

	if val.Map == nil {
		return VotingEndedEvent{}, errors.New("nil unstakes map")
	}

	unstakes, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingEndedEvent{}, fmt.Errorf("failed to parse unstakes map")
	}

	votingEnded.Unstakes = unstakes

	val, ok = event.Data["stakes"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDMap {
		return VotingEndedEvent{}, errors.New("invalid stakes value in event")
	}

	if val.Map == nil {
		return VotingEndedEvent{}, errors.New("nil stakes map")
	}

	stakes, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingEndedEvent{}, fmt.Errorf("failed to parse stakes map")
	}

	votingEnded.Stakes = stakes

	val, ok = event.Data["burns"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDMap {
		return VotingEndedEvent{}, errors.New("invalid burns value in event")
	}

	if val.Map == nil {
		return VotingEndedEvent{}, errors.New("nil burns map")
	}

	burns, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingEndedEvent{}, fmt.Errorf("failed to parse burns map")
	}

	votingEnded.Burns = burns

	val, ok = event.Data["mints"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDMap {
		return VotingEndedEvent{}, errors.New("invalid mints value in event")
	}

	if val.Map == nil {
		return VotingEndedEvent{}, errors.New("nil mints map")
	}

	mints, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingEndedEvent{}, fmt.Errorf("failed to parse unstakes map")
	}

	votingEnded.Mints = mints

	return votingEnded, nil
}
