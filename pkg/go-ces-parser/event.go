package ces

import (
	"bytes"
	"strings"

	"casper-dao-middleware/pkg/casper/types"
)

type ParseResult struct {
	Error error
	Event Event
}

type Event struct {
	ContractHash        types.Hash
	ContractPackageHash types.Hash
	Data                map[string]types.CLValue
	Name                string
}

// ParseEventNameAndData parse provided rawEvent according to event schema, return EventName and EventData
func ParseEventNameAndData(eventHex string, schemas Schemas) (string, map[string]types.CLValue, error) {
	clValue, _, err := types.ParseCLValueFromBytesWithRemainder(eventHex)
	if err != nil {
		return "", nil, ErrInvalidEventBytes
	}

	if len(clValue.Bytes) < 4 {
		return "", nil, ErrInvalidEventBytes
	}

	eventNameWithPrefix, eventBody, err := types.ParseBytesWithRemainder(clValue.Bytes[4:])
	if err != nil {
		return "", nil, err
	}

	if !bytes.HasPrefix(eventNameWithPrefix, []byte(eventPrefix)) {
		return "", nil, ErrNoEventPrefixInEvent
	}

	eventName := strings.TrimPrefix(string(eventNameWithPrefix), eventPrefix)
	schema, ok := schemas[eventName]
	if !ok {
		return "", nil, ErrEventNameNotInSchema
	}

	eventData, err := parseEventDataFromSchemaBytes(schema, eventBody)
	if err != nil {
		return "", nil, err
	}

	return eventName, eventData, nil
}

func parseEventDataFromSchemaBytes(schema Schema, data []byte) (map[string]types.CLValue, error) {
	result := make(map[string]types.CLValue, len(schema))

	var (
		err       error
		remainder = data
		clValue   types.CLValue
	)

	for _, item := range schema {
		clValue, remainder, err = types.NewCLValueFromBytesWithRemainder(item.Value, remainder)
		if err != nil {
			return nil, err
		}
		result[item.Property] = clValue
	}
	return result, nil
}
