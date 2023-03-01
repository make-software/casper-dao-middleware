package reputation_voter

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const BallotCastEventName = "BallotCast"

type BallotCastEvent struct {
	Voter      types.Address
	VotingType uint8
	Choice     types.Choice
	VotingID   uint32
	Stake      casper_types.U512
}

func ParseBallotCastEvent(event ces.Event) (BallotCastEvent, error) {
	var ballotCast BallotCastEvent
	var err error

	val, ok := event.Data["voter"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDKey {
		return BallotCastEvent{}, errors.New("invalid voter value in event")
	}
	ballotCast.Voter, err = types.NewAddressFromCLValue(val)
	if err != nil {
		return BallotCastEvent{}, err
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return BallotCastEvent{}, errors.New("invalid voting_id value in event")
	}
	ballotCast.VotingID = *val.U32

	val, ok = event.Data["voting_type"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return BallotCastEvent{}, errors.New("invalid voting_type value in event")
	}
	ballotCast.VotingType = *val.U8

	val, ok = event.Data["choice"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return BallotCastEvent{}, errors.New("invalid choice value in event")
	}
	ballotCast.Choice, err = types.NewChoiceFromByte(*val.U8)
	if err != nil {
		return BallotCastEvent{}, err
	}

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU512 {
		return BallotCastEvent{}, errors.New("invalid stake value in event")
	}
	ballotCast.Stake = *val.U512

	return ballotCast, nil
}
