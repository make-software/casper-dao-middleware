package admin

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const BallotCanceledEventName = "BallotCanceled"

type BallotCanceledEvent struct {
	Voter      types.Address
	VotingType uint8
	Choice     types.Choice
	VotingID   uint32
	Stake      casper_types.U512
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
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return BallotCanceledEvent{}, errors.New("invalid voting_id value in event")
	}
	ballotCanceled.VotingID = *val.U32

	val, ok = event.Data["voting_type"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return BallotCanceledEvent{}, errors.New("invalid voting_type value in event")
	}
	ballotCanceled.VotingType = *val.U8

	val, ok = event.Data["choice"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return BallotCanceledEvent{}, errors.New("invalid choice value in event")
	}
	ballotCanceled.Choice, err = types.NewChoiceFromByte(*val.U8)
	if err != nil {
		return BallotCanceledEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return BallotCanceledEvent{}, errors.New("invalid stake value in event")
	}
	ballotCanceled.Stake = *val.U512

	return ballotCanceled, nil
}
