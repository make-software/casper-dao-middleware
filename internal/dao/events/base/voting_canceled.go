package base

import (
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
)

const VotingCanceledEventName = "VotingCanceled"

type VotingCanceledEvent struct {
	VotingID   uint32
	VotingType uint8
	Unstakes   map[types.Tuple2]clvalue.UInt512
}

func ParseVotingCanceledEvent(event ces.Event) (VotingCanceledEvent, error) {
	var (
		votingCanceled VotingCanceledEvent
		err            error
	)

	val, ok := event.Data["voting_id"]
	if !ok || val.Type != cltype.UInt32 {
		return VotingCanceledEvent{}, errors.New("invalid voting_id value in event")
	}
	votingCanceled.VotingID = val.UI32.Value()

	val, ok = event.Data["voting_type"]
	if !ok || val.Type != cltype.UInt8 {
		return VotingCanceledEvent{}, errors.New("invalid voting_type value in event")
	}
	votingCanceled.VotingType = val.UI8.Value()

	val, ok = event.Data["unstakes"]
	if !ok {
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
