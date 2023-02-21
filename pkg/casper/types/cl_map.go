package types

type CLMap struct {
	KeyType   CLType
	ValueType CLType
	// key of the map is encoded CLValue of the KeyType type
	Data map[*CLValue]CLValue
}
