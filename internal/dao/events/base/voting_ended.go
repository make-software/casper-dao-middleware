package base

import (
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const VotingEndedEventName = "VotingEnded"

type VotingEndedEvent struct {
	VotingID             uint32
	VotingType           types.VotingType
	VotingResult         uint8
	StakeInFavour        clvalue.UInt512
	StakeAgainst         clvalue.UInt512
	UnboundStakeInFavour clvalue.UInt512
	UnboundStakeAgainst  clvalue.UInt512
	VotesInFavor         uint32
	VotesAgainst         uint32
	Unstakes             map[types.Tuple2]clvalue.UInt512
	Stakes               map[types.Tuple2]clvalue.UInt512
	Burns                map[types.Tuple2]clvalue.UInt512
	Mints                map[types.Tuple2]clvalue.UInt512
}

func ParseVotingEndedEvent(event ces.Event) (VotingEndedEvent, error) {
	var votingEnded VotingEndedEvent
	var err error

	val, ok := event.Data["voting_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingEndedEvent{}, errors.New("invalid voting_id value in event")
	}
	votingEnded.VotingID = val.UI32.Value()

	val, ok = event.Data["voting_type"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingEndedEvent{}, errors.New("invalid voting_type value in event")
	}

	votingEnded.VotingType, err = types.NewVotingTypeFromByte(byte(val.UI32.Value()))
	if err != nil {
		return VotingEndedEvent{}, err
	}

	val, ok = event.Data["voting_result"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingEndedEvent{}, errors.New("invalid voting_result value in event")
	}
	votingEnded.VotingResult = byte(val.UI32.Value())

	val, ok = event.Data["stake_in_favor"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingEndedEvent{}, errors.New("invalid stake_in_favor value in event")
	}
	votingEnded.StakeInFavour = *val.UI512

	val, ok = event.Data["stake_against"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingEndedEvent{}, errors.New("invalid stake_against value in event")
	}
	votingEnded.StakeAgainst = *val.UI512

	val, ok = event.Data["unbound_stake_in_favor"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingEndedEvent{}, errors.New("invalid unbound_stake_in_favor value in event")
	}
	votingEnded.UnboundStakeInFavour = *val.UI512

	val, ok = event.Data["unbound_stake_against"]
	if !ok || val.Type != cltype.UInt512 {
		return VotingEndedEvent{}, errors.New("invalid unbound_stake_against value in event")
	}
	votingEnded.UnboundStakeAgainst = *val.UI512

	val, ok = event.Data["votes_in_favor"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingEndedEvent{}, errors.New("invalid votes_in_favor value in event")
	}
	votingEnded.VotesInFavor = val.UI32.Value()

	val, ok = event.Data["votes_against"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingEndedEvent{}, errors.New("invalid votes_against value in event")
	}
	votingEnded.VotesInFavor = val.UI32.Value()

	val, ok = event.Data["unstakes"]
	if val.Map == nil {
		return VotingEndedEvent{}, errors.New("nil unstakes map")
	}

	unstakes, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingEndedEvent{}, fmt.Errorf("failed to parse unstakes map")
	}

	votingEnded.Unstakes = unstakes

	val, ok = event.Data["stakes"]
	if !ok {
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
	if !ok {
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
	if !ok {
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
