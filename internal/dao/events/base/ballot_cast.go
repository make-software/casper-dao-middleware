package base

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const BallotCastEventName = "BallotCast"

type BallotCastEvent struct {
	Voter      types.Address
	VotingType types.VotingType
	Choice     types.Choice
	VotingID   uint32
	Stake      clvalue.UInt512
}

func ParseBallotCastEvent(event ces.Event) (BallotCastEvent, error) {
	var ballotCast BallotCastEvent
	var err error

	val, ok := event.Data["voter"]
	if !ok || val.Type != cltype.Key {
		return BallotCastEvent{}, errors.New("invalid voter value in event")
	}
	ballotCast.Voter, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return BallotCastEvent{}, err
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCastEvent{}, errors.New("invalid voting_id value in event")
	}
	ballotCast.VotingID = val.UI32.Value()

	val, ok = event.Data["voting_type"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCastEvent{}, errors.New("invalid voting_type value in event")
	}

	ballotCast.VotingType, err = types.NewVotingTypeFromByte(byte(val.UI32.Value()))
	if err != nil {
		return BallotCastEvent{}, err
	}

	val, ok = event.Data["choice"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCastEvent{}, errors.New("invalid choice value in event")
	}

	ballotCast.Choice, err = types.NewChoiceFromByte(byte(val.UI32.Value()))
	if err != nil {
		return BallotCastEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type != cltype.UInt512 {
		return BallotCastEvent{}, errors.New("invalid stake value in event")
	}
	ballotCast.Stake = *val.UI512

	return ballotCast, nil
}
