package ces

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"casper-dao-middleware/pkg/casper/types"
)

func Test_NewSchemaFromBytesWithRemainder(t *testing.T) {
	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("amount", u8::cl_type());
	//expected_schema.with_elem("from", Key::cl_type());

	hexStr := `0200000006000000616d6f756e74030400000066726f6d0b`

	res, _ := hex.DecodeString(hexStr)
	schema, _, err := newSchemaFromBytesWithRemainder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "amount")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeIDU8)
	assert.Equal(t, schema[1].Property, "from")
	assert.Equal(t, schema[1].Value.CLTypeID, types.CLTypeIDKey)

	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("option", CLType::Option(Box::new(Key::cl_type())));

	hexStr = `01000000060000006f7074696f6e0d0b`

	res, _ = hex.DecodeString(hexStr)
	schema, _, err = newSchemaFromBytesWithRemainder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "option")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeIDOption)
	assert.Equal(t, schema[0].Value.CLTypeOption.CLTypeInner.CLTypeID, types.CLTypeIDKey)

	//let mut expected_schema = Schema::new();
	//expected_schema.with_elem("list", CLType::List(Box::new(Key::cl_type())));
	//expected_schema.with_elem("any", CLType::Any);

	hexStr = `02000000040000006c6973740e0b03000000616e7915`

	res, _ = hex.DecodeString(hexStr)
	schema, _, err = newSchemaFromBytesWithRemainder(res)
	assert.NoError(t, err)

	assert.Equal(t, schema[0].Property, "list")
	assert.Equal(t, schema[0].Value.CLTypeID, types.CLTypeIDList)
	assert.NotNil(t, schema[0].Value.CLTypeList)
	assert.Equal(t, schema[1].Property, "any")
	assert.Equal(t, schema[1].Value.CLTypeID, types.CLTypeIDAny)
}
