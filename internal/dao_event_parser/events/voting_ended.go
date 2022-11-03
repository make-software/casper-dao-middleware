package events

import (
	"encoding/binary"

	"casper-dao-middleware/pkg/casper/types"
)

const VotingEndedEventName = "VotingEnded"

type VotingEnded struct {
	VotingID         types.U256
	InformalVotingID types.U256
	FormalVotingID   *types.U256
	Result           string
	VotesCount       types.U256
	StakeInFavour    types.U256
	StakeAgainst     types.U256
	Transfers        map[string]types.U256
	Burns            map[string]types.U256
	Mints            map[string]types.U256
}

func ParseVotingEndedEvent(rawBytes []byte) (VotingEnded, error) {
	event := VotingEnded{}

	votingID, reminder, err := types.ParseU256FromBytes(rawBytes)
	if err != nil {
		return VotingEnded{}, err
	}
	event.VotingID = votingID

	event.InformalVotingID, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	// check if formal FormalVotingID is set
	if reminder[0] != 0 {
		reminder = reminder[1:]
		var formalVotingID types.U256
		formalVotingID, reminder, err = types.ParseU256FromBytes(reminder)
		if err != nil {
			return VotingEnded{}, err
		}
		event.FormalVotingID = &formalVotingID
	}

	resultBytes, reminder, err := types.ParseBytesWithReminder(reminder)
	if err != nil {
		return VotingEnded{}, err
	}
	event.Result = string(resultBytes)

	event.VotesCount, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	event.StakeInFavour, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	event.StakeAgainst, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	event.Transfers, reminder, err = parseHashU256MapFromRemainder(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	event.Burns, reminder, err = parseHashU256MapFromRemainder(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	event.Mints, reminder, err = parseHashU256MapFromRemainder(reminder)
	if err != nil {
		return VotingEnded{}, err
	}

	return event, nil
}

func parseHashU256MapFromRemainder(reminder []byte) (map[string]types.U256, []byte, error) {
	itemsCount := binary.LittleEndian.Uint32(reminder)
	reminder = reminder[4:]

	if itemsCount == 0 {
		return nil, reminder, nil
	}

	itemsMap := make(map[string]types.U256, itemsCount)
	for i := 0; i < int(itemsCount); i++ {
		// shift separator to new transfer row
		reminder = reminder[1:]
		address, err := types.NewHashFromRawBytes(reminder[:32])
		if err != nil {
			return nil, nil, err
		}

		reminder = reminder[32:]
		var addressTransfer types.U256
		addressTransfer, reminder, err = types.ParseU256FromBytes(reminder)
		if err != nil {
			return nil, nil, err
		}
		itemsMap[address.ToHex()] = addressTransfer
	}
	return itemsMap, reminder, nil
}
