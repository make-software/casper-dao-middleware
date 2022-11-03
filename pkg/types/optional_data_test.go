package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseOptionalFieldsTree(t *testing.T) {
	t.Run("Success Schema Parsing", func(t *testing.T) {
		t.Run("simple schema", func(t *testing.T) {
			schema := "account_info{owner}"
			accountInfoFields, err := ParseOptionalData(schema)
			assert.True(t, len(accountInfoFields.GetNested()) == 1)

			assert.NoError(t, err)
			assert.Equal(t, "root", accountInfoFields.GetName())
			assert.Equal(t, "owner", accountInfoFields.GetNestedByName("account_info").GetNestedByName("owner").GetName())
		})

		t.Run("simple function schema", func(t *testing.T) {
			schema := "rate(1)"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			rateFuncArgs, ok := optionalData.ContainsFunc("rate")
			assert.True(t, ok)

			assert.Equal(t, 1, len(rateFuncArgs))
			assert.Equal(t, "1", rateFuncArgs[0])
		})

		t.Run("invalid function schema", func(t *testing.T) {
			schema := "rate(1))"
			_, err := ParseOptionalData(schema)
			assert.Error(t, err)
		})

		t.Run("simple function schema with several params", func(t *testing.T) {
			schema := "rate(1, data)"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			rateFuncArgs, ok := optionalData.ContainsFunc("rate")
			assert.True(t, ok)

			assert.Equal(t, 2, len(rateFuncArgs))
			assert.Equal(t, "1", rateFuncArgs[0])
			assert.Equal(t, "data", rateFuncArgs[1])
		})

		t.Run("schema with several functions", func(t *testing.T) {
			schema := "rate(1), block(data)"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			rateFuncArgs, ok := optionalData.ContainsFunc("rate")
			assert.True(t, ok)

			blockFuncArgs, ok := optionalData.ContainsFunc("block")
			assert.True(t, ok)

			assert.Equal(t, 1, len(rateFuncArgs))
			assert.Equal(t, "1", rateFuncArgs[0])

			assert.Equal(t, 1, len(blockFuncArgs))
			assert.Equal(t, "data", blockFuncArgs[0])
		})

		t.Run("schema with function and field", func(t *testing.T) {
			schema := "account_info{owner, name}, rate(1)"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			assert.True(t, len(optionalData.GetNested()) == 2)

			rateFuncArgs, ok := optionalData.ContainsFunc("rate")
			assert.True(t, ok)

			assert.Equal(t, 1, len(rateFuncArgs))
			assert.Equal(t, "1", rateFuncArgs[0])

			accountInfo := optionalData.GetNestedByName("account_info")
			assert.Equal(t, "account_info", accountInfo.GetName())
			assert.Equal(t, "owner", accountInfo.GetNestedByName("owner").GetName())
			assert.Equal(t, "name", accountInfo.GetNestedByName("name").GetName())
		})

		t.Run("simple schema with multiple fields", func(t *testing.T) {
			schema := "account_info{owner, name} , data{inner}"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			assert.True(t, len(optionalData.GetNested()) == 2)

			accountInfo := optionalData.GetNestedByName("account_info")
			assert.Equal(t, "account_info", accountInfo.GetName())
			assert.Equal(t, "owner", accountInfo.GetNestedByName("owner").GetName())
			assert.Equal(t, "name", accountInfo.GetNestedByName("name").GetName())

			dataField := optionalData.GetNestedByName("data")
			assert.Equal(t, "data", dataField.GetName())
			assert.Equal(t, "inner", dataField.GetNestedByName("inner").GetName())
		})

		t.Run("simple schema", func(t *testing.T) {
			schema := "account_info{owner{name,branding}}"
			optionalData, err := ParseOptionalData(schema)
			assert.True(t, len(optionalData.GetNested()) == 1)
			assert.NoError(t, err)

			accountInfo := optionalData.GetNestedByName("account_info")

			assert.Equal(t, "owner", accountInfo.GetNestedByName("owner").GetName())
			assert.Equal(t, "name", accountInfo.GetNestedByName("owner").GetNestedByName("name").GetName())
			assert.Equal(t, "branding", accountInfo.GetNestedByName("owner").GetNestedByName("branding").GetName())
		})

		t.Run("with indentation", func(t *testing.T) {
			schema := "account_info   {           owner{             name,     branding       }}"
			optionalData, err := ParseOptionalData(schema)
			assert.True(t, len(optionalData.GetNested()) == 1)
			assert.NoError(t, err)

			accountInfo := optionalData.GetNestedByName("account_info")

			assert.Equal(t, "account_info", accountInfo.GetName())
			assert.Equal(t, "owner", accountInfo.GetNestedByName("owner").GetName())
			assert.Equal(t, "name", accountInfo.GetNestedByName("owner").GetNestedByName("name").GetName())
			assert.Equal(t, "branding", accountInfo.GetNestedByName("owner").GetNestedByName("branding").GetName())
		})

		t.Run("func with indentation", func(t *testing.T) {
			schema := " rate (     5   )"
			optionalData, err := ParseOptionalData(schema)
			assert.NoError(t, err)

			rateFuncArgs, ok := optionalData.ContainsFunc("rate")
			assert.True(t, ok)

			assert.Equal(t, 1, len(rateFuncArgs))
			assert.Equal(t, "5", rateFuncArgs[0])
		})

		t.Run("complex schema", func(t *testing.T) {
			schema := "account_info{owner{name,branding,value{field1, filed2}}, cert{data}, role}"
			optionalData, err := ParseOptionalData(schema)
			assert.True(t, len(optionalData.GetNested()) == 1)
			assert.NoError(t, err)

			accountInfo := optionalData.GetNestedByName("account_info")

			assert.Equal(t, "account_info", accountInfo.GetName())
			assert.Equal(t, "owner", accountInfo.GetNestedByName("owner").GetName())
			assert.Equal(t, "cert", accountInfo.GetNestedByName("cert").GetName())
			assert.Equal(t, "role", accountInfo.GetNestedByName("role").GetName())
			assert.Equal(t, "data", accountInfo.GetNestedByName("cert").GetNestedByName("data").GetName())
			assert.Equal(t, "name", accountInfo.GetNestedByName("owner").GetNestedByName("name").GetName())
			assert.Equal(t, "branding", accountInfo.GetNestedByName("owner").GetNestedByName("branding").GetName())
			assert.Equal(t, "value", accountInfo.GetNestedByName("owner").GetNestedByName("value").GetName())
			assert.Equal(t, "field1", accountInfo.GetNestedByName("owner").GetNestedByName("value").GetNestedByName("field1").GetName())
			assert.Equal(t, "filed2", accountInfo.GetNestedByName("owner").GetNestedByName("value").GetNestedByName("filed2").GetName())
		})
	})

	t.Run("Invalid Schema Parsing", func(t *testing.T) {
		t.Run("inconsistent curly brackets number", func(t *testing.T) {
			schema := "account_info{owner{name,branding}}}"
			_, err := ParseOptionalData(schema)

			assert.Error(t, err)
		})

		t.Run("redundant commas", func(t *testing.T) {
			schema := "account_info{owner{name, ,branding}}}"
			_, err := ParseOptionalData(schema)

			assert.Error(t, err)
		})
	})
}
