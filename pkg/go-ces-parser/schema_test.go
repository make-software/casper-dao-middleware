package ces

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"casper-middleware/pkg/casper/types"
)

func Test_NewSchemaFromBytesWithReminder(t *testing.T) {
	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("amount", u8::cl_type());
	//expected_schema.with_elem("from", Key::cl_type());

	hexStr := `0200000006000000616d6f756e74030400000066726f6d0b`

	res, _ := hex.DecodeString(hexStr)
	schema, _, err := newSchemaFromBytesWithReminder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "amount")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeU8)
	assert.Equal(t, schema[1].Property, "from")
	assert.Equal(t, schema[1].Value.CLTypeID, types.CLTypeKey)

	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("option", CLType::Option(Box::new(Key::cl_type())));

	hexStr = `01000000060000006f7074696f6e0d0b`

	res, _ = hex.DecodeString(hexStr)
	schema, _, err = newSchemaFromBytesWithReminder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "option")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeOption)
	assert.Equal(t, schema[0].Value.CLType.CLTypeID, types.CLTypeKey)

	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("list", CLType::List(Box::new(Key::cl_type())));
	//expected_schema.with_elem("any", CLType::Any);

	hexStr = `02000000040000006c6973740e0b03000000616e7915`

	res, _ = hex.DecodeString(hexStr)
	schema, _, err = newSchemaFromBytesWithReminder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "list")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeList)
	assert.Equal(t, schema[0].Value.CLType.CLTypeID, types.CLTypeKey)
	assert.Equal(t, schema[1].Property, "any")
	assert.Equal(t, schema[1].Value.CLTypeID, types.CLTypeAny)
}
