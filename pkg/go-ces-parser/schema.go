package ces

import (
	"encoding/binary"
	"errors"
	"log"

	"casper-dao-middleware/pkg/casper/types"
)

type (
	Schemas map[string]Schema

	Schema             []PropertyDefinition
	PropertyDefinition struct {
		Property string
		Value    types.CLType
	}
)

func NewSchemasFromBytes(rawSchemas []byte) (Schemas, error) {
	schemasNumber := binary.LittleEndian.Uint32(rawSchemas)
	if schemasNumber == 0 || int(schemasNumber) > len(rawSchemas) {
		return nil, errors.New("invalid schemasNumber value")
	}

	// without uint32 (schemasNumber)
	var reminder = rawSchemas[4:]
	var err error

	schemas := make(map[string]Schema, schemasNumber)
	for i := 0; i < int(schemasNumber); i++ {
		var schemaName []byte
		var schema Schema

		schemaName, reminder, err = types.ParseBytesWithReminder(reminder)
		if err != nil {
			return nil, err
		}

		schema, reminder, err = newSchemaFromBytesWithReminder(reminder)
		if err != nil {
			return nil, err
		}

		schemas[string(schemaName)] = schema
	}
	return schemas, nil
}

func newSchemaFromBytesWithReminder(bytes []byte) (Schema, []byte, error) {
	itemNumber := binary.LittleEndian.Uint32(bytes)
	if int(itemNumber) > len(bytes) {
		return nil, nil, errors.New("invalid itemNumber value")
	}

	reminder := make([]byte, len(bytes)-4)
	copy(reminder, bytes[4:])

	var err error
	schema := make([]PropertyDefinition, 0, int(itemNumber))
	for i := 0; i < int(itemNumber); i++ {
		var item []byte
		item, reminder, err = types.ParseBytesWithReminder(reminder)
		if err != nil {
			log.Println("failed to parse Schema item")
			return nil, nil, err
		}

		var clType types.CLType
		clType, reminder, err = types.ClTypeFromBytes(0, reminder)
		if err != nil {
			return nil, nil, err
		}

		schema = append(schema, newPropertyDefinition(string(item), clType))
	}
	return schema, reminder, nil
}

func newPropertyDefinition(name string, value types.CLType) PropertyDefinition {
	return PropertyDefinition{
		Property: name,
		Value:    value,
	}
}
