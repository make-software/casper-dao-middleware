package events

import (
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"

	"go.uber.org/zap"
)

const BallotCastName = "BallotCast"

type Choice byte

const (
	ChoiceAgainst Choice = 1
	ChoiceInFavor Choice = 2
)

type BallotCast struct {
	Voter      dao_types.Address
	VotingType uint8
	Choice     Choice
	VotingID   uint32
	Stake      types.U512
}

func ParseBallotCastEvent(event ces.Event) (BallotCast, error) {
	var ballotCast BallotCast

	val, ok := event.Data["voter"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDKey {
		return BallotCast{}, errors.New("invalid voter value in event")
	}
	ballotCast.Voter = dao_types.Address{
		AccountHash:         val.Key.AccountHash,
		ContractPackageHash: val.Key.Hash,
	}

	val, ok = event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU32 {
		return BallotCast{}, errors.New("invalid voting_id value in event")
	}
	ballotCast.VotingID = *val.U32

	val, ok = event.Data["voting_type"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU8 {
		return BallotCast{}, errors.New("invalid voting_type value in event")
	}
	ballotCast.VotingType = *val.U8

	val, ok = event.Data["choice"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU8 {
		return BallotCast{}, errors.New("invalid choice value in event")
	}
	ballotCast.Choice = Choice(*val.U8)

	val, ok = event.Data["stake"]
	if !ok || val.Type.CLTypeID != types.CLTypeIDU512 {
		return BallotCast{}, errors.New("invalid stake value in event")
	}
	ballotCast.Stake = *val.U512

	zap.S().Info("Successfully parsed BallotCast event")
	return ballotCast, nil
}
