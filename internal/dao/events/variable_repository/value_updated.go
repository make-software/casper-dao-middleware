package variable_repository

import (
	"errors"

	"casper-dao-middleware/internal/dao/types"
	casper_types "casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

const ValueUpdatedEventName = "ValueUpdated"

type ValueUpdatedEvent struct {
	Key            string
	Value          types.RecordValue
	ActivationTime *uint64
}

func ParseValueUpdatedEvent(event ces.Event) (ValueUpdatedEvent, error) {
	var valueUpdated ValueUpdatedEvent

	val, ok := event.Data["key"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDString {
		return ValueUpdatedEvent{}, errors.New("invalid key value in event")
	}
	valueUpdated.Key = *val.String

	val, ok = event.Data["value"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDList {
		return ValueUpdatedEvent{}, errors.New("invalid key value in event")
	}

	listClValue := *val.List
	if len(listClValue) == 0 || listClValue[0].Type.CLTypeID != casper_types.CLTypeIDU8 {
		return ValueUpdatedEvent{}, errors.New("expected List(u8) for value field")
	}

	recordValueBytes := make([]byte, 0, len(listClValue))
	for _, clValue := range listClValue {
		recordValueBytes = append(recordValueBytes, *clValue.U8)
	}

	var err error
	valueUpdated.Value, err = types.NewRecordValueFromBytes(recordValueBytes)
	if err != nil {
		return ValueUpdatedEvent{}, err
	}

	val, ok = event.Data["activation_time"]
	if !ok || val.Type.CLTypeID != casper_types.CLTypeIDOption {
		return ValueUpdatedEvent{}, errors.New("invalid activation_time value in event")
	}

	if val.Option != nil {
		if val.Option.U64 == nil {
			return ValueUpdatedEvent{}, errors.New("expected U64 in Option(Some())")
		}

		valueUpdated.ActivationTime = val.Option.U64
	}

	return valueUpdated, nil
}
