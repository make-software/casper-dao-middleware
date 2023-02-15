package types

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type RawCLValue struct {
	Type  CLType
	Bytes []byte
}

// TODO: extend CLValue type
type CLValue struct {
	Type   CLType
	Bool   *bool
	I32    *int32
	I64    *int64
	U8     *uint8
	U32    *uint32
	U64    *uint64
	U128   *U128
	U256   *U256
	U512   *U512
	String *string
	Key    *Key
	Option *CLValue
	List   *[]CLValue
	//ByteArray *FixedByteArray
	//Result    *CLValueResult
	//Map       *CLMap
	Tuple1    *[1]CLValue
	Tuple2    *[2]CLValue
	Tuple3    *[3]CLValue
	PublicKey *PublicKey
	Any       []byte
}

func NewCLValueFromBytesWithReminder(clType CLType, data []byte) (CLValue, []byte, error) {
	reminder := data
	switch clType.CLTypeID {
	case CLTypeIDU8:
		if len(data) < 1 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}
		return CLValue{
			Type: clType,
			U8:   &data[0],
		}, reminder[1:], nil
	case CLTypeIDBool:
		if len(data) < 1 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}
		var res bool
		if data[0] == 1 {
			res = false
		}
		return CLValue{
			Type: clType,
			Bool: &res,
		}, reminder[1:], nil
	case CLTypeIDU32:
		if len(data) < 4 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		value := binary.LittleEndian.Uint32(data)

		return CLValue{
			Type: clType,
			U32:  &value,
		}, reminder[4:], nil

	case CLTypeIDU64:
		if len(data) < 8 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		value := binary.LittleEndian.Uint64(data)
		return CLValue{
			Type: clType,
			U64:  &value,
		}, reminder[8:], nil
	case CLTypeIDU128:
		val, reminder, err := ParseUTypeFromBytes[U128](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U128: &val,
		}, reminder, nil
	case CLTypeIDU256:
		val, reminder, err := ParseUTypeFromBytes[U256](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U256: &val,
		}, reminder, nil
	case CLTypeIDU512:
		val, reminder, err := ParseUTypeFromBytes[U512](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U512: &val,
		}, reminder, nil

	case CLTypeIDKey:
		key, reminder, err := ParseKeyFromBytes(reminder)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			Key:  &key,
		}, reminder, nil
	case CLTypeIDAny:
		rawParsed, reminder, err := ParseBytesWithReminder(data)
		if err != nil {
			return CLValue{}, nil, err
		}

		return CLValue{
			Type: clType,
			Any:  rawParsed,
		}, reminder, nil
	case CLTypeIDString:
		rawParsed, reminder, err := ParseBytesWithReminder(data)
		if err != nil {
			return CLValue{}, nil, err
		}

		parsed := string(rawParsed)
		return CLValue{
			Type:   clType,
			String: &parsed,
		}, reminder, nil
	case CLTypeIDOption:
		if reminder[0] != 0 && clType.CLTypeOption != nil {
			reminder = reminder[1:]
			clValue, reminder, err := NewCLValueFromBytesWithReminder(clType.CLTypeOption.CLTypeInner, reminder)
			if err != nil {
				return CLValue{}, nil, err
			}
			return clValue, reminder, nil
		}
		reminder = reminder[1:]

		return CLValue{
			Type: clType,
		}, reminder, nil
	}

	return CLValue{}, reminder, errors.New("unknown CLType provided")
}

func ParseCLValueFromBytesWithReminder(data string) (RawCLValue, []byte, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	bytes, reminder, err := ParseBytesWithReminder(decoded)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	clType, reminder, err := ClTypeFromBytes(0, reminder)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	return RawCLValue{
		Type:  clType,
		Bytes: bytes,
	}, reminder, nil
}

func newInvalidLengthErr(clTypeID CLTypeID) error {
	return fmt.Errorf("invalid bytes length value for type - %s", clTypeID.ToString())
}
