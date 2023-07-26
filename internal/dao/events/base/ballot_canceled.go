package base

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const BallotCanceledEventName = "BallotCanceled"

type BallotCanceledEvent struct {
	Voter      types.Address
	VotingType types.VotingType
	Choice     types.Choice
	VotingID   uint32
	Stake      clvalue.UInt512
}

func ParseBallotCanceledEvent(event ces.Event) (BallotCanceledEvent, error) {
	var ballotCanceled BallotCanceledEvent
	var err error

	val, ok := event.Data["voter"]
	if !ok {
		return BallotCanceledEvent{}, errors.New("invalid voter value in event")
	}
	ballotCanceled.Voter, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return BallotCanceledEvent{}, err
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCanceledEvent{}, errors.New("invalid voting_id value in event")
	}
	ballotCanceled.VotingID = val.UI32.Value()

	val, ok = event.Data["voting_type"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCanceledEvent{}, errors.New("invalid voting_type value in event")
	}

	ballotCanceled.VotingType, err = types.NewVotingTypeFromByte(byte(val.UI32.Value()))
	if err != nil {
		return BallotCanceledEvent{}, err
	}

	val, ok = event.Data["choice"]
	if !ok || val.Type != cltype.UInt32 {
		return BallotCanceledEvent{}, errors.New("invalid choice value in event")
	}
	ballotCanceled.Choice, err = types.NewChoiceFromByte(byte(val.UI32.Value()))
	if err != nil {
		return BallotCanceledEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type != cltype.UInt512 {
		return BallotCanceledEvent{}, errors.New("invalid stake value in event")
	}
	ballotCanceled.Stake = *val.UI512

	return ballotCanceled, nil
}
