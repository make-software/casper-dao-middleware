package events

import (
	"encoding/binary"
	"errors"
	"math/big"

	dao_types "casper-dao-middleware/internal/crdao/dao_event_parser/types"
	"casper-dao-middleware/pkg/casper/types"
)

const VotingCreatedEventName = "VotingCreated"

type VotingCreated struct {
	Creator                       dao_types.Address
	VotingID                      types.U256
	InformalVotingID              types.U256
	FormalVotingID                *types.U256
	ConfigFormalVotingQuorum      *big.Int
	ConfigFormalVotingTime        uint64
	ConfigInformalVotingQuorum    *big.Int
	ConfigInformalVotingTime      uint64
	ConfigCreateMinimumReputation *big.Int
}

func ParseVotingCreatedEvent(bytes []byte) (VotingCreated, error) {
	key, reminder, err := types.ParseKeyFromBytes(bytes)
	if err != nil {
		return VotingCreated{}, err
	}

	event := VotingCreated{}

	if key.AccountHash == nil && key.Hash == nil {
		return VotingCreated{}, errors.New("expected Creator in VotingCreated event")
	}
	var creator dao_types.Address
	if key.AccountHash != nil {
		creator.AccountHash = key.AccountHash
	} else {
		creator.ContractPackageHash = key.Hash
	}

	event.Creator = creator

	event.VotingID, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingCreated{}, err
	}

	event.InformalVotingID, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingCreated{}, err
	}

	if reminder[0] == 0 {
		reminder = reminder[1:]
	}

	event.ConfigFormalVotingQuorum, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingCreated{}, err
	}

	event.ConfigFormalVotingTime = binary.LittleEndian.Uint64(reminder)
	// skip 8 bytes
	reminder = reminder[8:]

	event.ConfigInformalVotingQuorum, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingCreated{}, err
	}

	event.ConfigInformalVotingTime = binary.LittleEndian.Uint64(reminder)
	// skip 8 bytes
	reminder = reminder[8:]

	event.ConfigCreateMinimumReputation, reminder, err = types.ParseU256FromBytes(reminder)
	if err != nil {
		return VotingCreated{}, err
	}

	return event, nil
}
