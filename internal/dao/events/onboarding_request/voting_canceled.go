package onboarding_request

import (
	"errors"
	"fmt"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const VotingCanceledEventName = "VotingCanceled"

type VotingCanceledEvent struct {
	VotingID   uint32
	VotingType uint8
	Unstakes   map[types.Tuple2]casper_types.U512
}

func ParseVotingCanceledEvent(event ces.Event) (VotingCanceledEvent, error) {
	var (
		votingCanceled VotingCanceledEvent
		err            error
	)

	val, ok := event.Data["voting_id"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU32 {
		return VotingCanceledEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCanceled.VotingID = *val.U32

	val, ok = event.Data["voting_type"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDU8 {
		return VotingCanceledEvent{}, errors.New("invalid voting_type value in event")
	}
	votingCanceled.VotingType = *val.U8

	val, ok = event.Data["unstakes"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDMap {
		return VotingCanceledEvent{}, errors.New("invalid unstakes value in event")
	}

	if val.Map == nil {
		return VotingCanceledEvent{}, errors.New("nil unstakes map")
	}

	unstakes, err := types.ParseTuple2U512MapFromCLValue(val)
	if err != nil {
		return VotingCanceledEvent{}, fmt.Errorf("failed to parse unstakes map")
	}

	votingCanceled.Unstakes = unstakes

	return votingCanceled, nil
}
