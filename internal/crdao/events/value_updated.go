package events

import (
	"encoding/binary"
	"errors"

	"casper-dao-middleware/internal/crdao/types"
	pkgTypes "casper-dao-middleware/pkg/casper/types"
)

const ValueUpdatedEventName = "ValueUpdated"

type ValueUpdated struct {
	Key            string
	Value          types.RecordValue
	ActivationTime *uint64
}

func ParseValueUpdatedEvent(bytes []byte) (ValueUpdated, error) {
	key, reminder, err := pkgTypes.ParseStringFromBytes(bytes)
	if err != nil {
		return ValueUpdated{}, err
	}

	recordValue, reminder, err := types.NewRecordValueFromBytesWithReminder(reminder)
	if err != nil {
		return ValueUpdated{}, err
	}

	if len(reminder) == 0 {
		return ValueUpdated{}, errors.New("invalid payload format expect ")
	}

	var activationTime *uint64
	// if reminder[0] equals to 1, parse activationTime
	if reminder[0] == 1 {
		value := binary.LittleEndian.Uint64(reminder[1:])
		activationTime = &value
	}

	return ValueUpdated{
		Key:            key,
		Value:          recordValue,
		ActivationTime: activationTime,
	}, nil
}
