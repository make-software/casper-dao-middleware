package variable_repository

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/types"
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
	if !ok || val.Type != cltype.String {
		return ValueUpdatedEvent{}, errors.New("invalid key value in event")
	}
	valueUpdated.Key = val.StringVal.String()

	val, ok = event.Data["value"]
	if !ok {
		return ValueUpdatedEvent{}, errors.New("invalid key value in event")
	}

	if val.List == nil {
		return ValueUpdatedEvent{}, errors.New("expected `value` key as list")
	}

	listClValue := val.List.Elements
	if len(listClValue) == 0 || listClValue[0].Type != cltype.UInt8 {
		return ValueUpdatedEvent{}, errors.New("expected List(u8) for value field")
	}

	recordValueBytes := make([]byte, 0, len(listClValue))
	for _, clValue := range listClValue {
		recordValueBytes = append(recordValueBytes, clValue.UI8.Value())
	}

	//var err error
	//valueUpdated.Value, err = types.NewRecordValueFromBytes(recordValueBytes)
	//if err != nil {
	//	return ValueUpdatedEvent{}, err
	//}

	val, ok = event.Data["activation_time"]
	if !ok {
		return ValueUpdatedEvent{}, errors.New("invalid activation_time value in event")
	}

	if val.Option != nil {
		if val.Option.Inner.UI64 == nil {
			return ValueUpdatedEvent{}, errors.New("expected U64 in Option(Some())")
		}

		value := val.Option.Inner.UI64.Value()
		valueUpdated.ActivationTime = &value
	}

	return valueUpdated, nil
}
