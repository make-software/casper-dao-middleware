package events

import (
	"encoding/binary"
	"errors"

	dao_types "casper-dao-middleware/internal/crdao/dao_event_parser/types"
	"casper-dao-middleware/pkg/casper/types"

	"go.uber.org/zap"
)

const BallotCastName = "BallotCast"

type Choice byte

const (
	ChoiceAgainst Choice = 1
	ChoiceInFavor Choice = 2
)

type BallotCast struct {
	Address  dao_types.Address
	Choice   Choice
	VotingID types.U256
	Stake    types.U256
}

func ParseBallotCastEvent(bytes []byte) (BallotCast, error) {
	key, reminder, err := types.ParseKeyFromBytes(bytes)
	if err != nil {
		return BallotCast{}, err
	}

	event := BallotCast{}
	if key.AccountHash == nil && key.Hash == nil {
		return BallotCast{}, errors.New("expected Address in BallotCast event")
	}

	var address dao_types.Address
	if key.AccountHash != nil {
		address.AccountHash = key.AccountHash
	} else {
		address.ContractPackageHash = key.Hash
	}

	event.Address = address

	event.VotingID, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return BallotCast{}, err
	}

	if len(reminder) == 0 {
		return BallotCast{}, errors.New("not full BallotCast event received")
	}

	choice := Choice(binary.LittleEndian.Uint32(reminder))
	if choice != ChoiceAgainst && choice != ChoiceInFavor {
		return BallotCast{}, errors.New("invalid choice value in event: expect Against(1) or InFavor(2)")
	}
	event.Choice = choice

	// extract choice from remainder
	reminder = reminder[4:]

	event.Stake, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return BallotCast{}, err
	}

	if len(reminder) != 0 {
		return BallotCast{}, errors.New("invalid BallotCast event: remainder is left")
	}

	zap.S().Info("Successfully parsed BallotCast event")
	return event, nil
}
