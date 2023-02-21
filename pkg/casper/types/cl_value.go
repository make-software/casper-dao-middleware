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
	Map       *CLMap
	Tuple1    *[1]CLValue
	Tuple2    *[2]CLValue
	Tuple3    *[3]CLValue
	PublicKey *PublicKey
	Any       []byte
}

func NewCLValueFromBytesWithRemainder(clType CLType, data []byte) (CLValue, []byte, error) {
	var (
		remainder = data
		err       error
	)
	switch clType.CLTypeID {
	case CLTypeIDU8:
		if len(data) < 1 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}
		return CLValue{
			Type: clType,
			U8:   &data[0],
		}, remainder[1:], nil
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
		}, remainder[1:], nil
	case CLTypeIDU32:
		if len(data) < 4 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		value := binary.LittleEndian.Uint32(data)

		return CLValue{
			Type: clType,
			U32:  &value,
		}, remainder[4:], nil

	case CLTypeIDU64:
		if len(data) < 8 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		value := binary.LittleEndian.Uint64(data)
		return CLValue{
			Type: clType,
			U64:  &value,
		}, remainder[8:], nil
	case CLTypeIDU128:
		var val U128
		val, remainder, err = ParseUTypeFromBytes[U128](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U128: &val,
		}, remainder, nil
	case CLTypeIDU256:
		var val U256
		val, remainder, err = ParseUTypeFromBytes[U256](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U256: &val,
		}, remainder, nil
	case CLTypeIDU512:
		var val U512
		val, remainder, err = ParseUTypeFromBytes[U512](data)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			U512: &val,
		}, remainder, nil

	case CLTypeIDKey:
		var key Key
		key, remainder, err = ParseKeyFromBytes(remainder)
		if err != nil {
			return CLValue{}, nil, err
		}
		return CLValue{
			Type: clType,
			Key:  &key,
		}, remainder, nil
	case CLTypeIDAny:
		var rawParsed []byte
		rawParsed, remainder, err = ParseBytesWithRemainder(data)
		if err != nil {
			return CLValue{}, nil, err
		}

		return CLValue{
			Type: clType,
			Any:  rawParsed,
		}, remainder, nil
	case CLTypeIDString:
		var rawParsed []byte
		rawParsed, remainder, err = ParseBytesWithRemainder(data)
		if err != nil {
			return CLValue{}, nil, err
		}

		parsed := string(rawParsed)
		return CLValue{
			Type:   clType,
			String: &parsed,
		}, remainder, nil
	case CLTypeIDMap:
		if len(data) < 4 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		if clType.CLTypeMap == nil {
			return CLValue{}, nil, errors.New("expected not nil ClMap")
		}

		itemNumber := binary.LittleEndian.Uint32(data)
		var (
			//without uint32
			rem                = data[4:]
			clValueK, clValueV CLValue
			err                error
		)

		mapData := make(map[*CLValue]CLValue, itemNumber)
		for i := 0; i < int(itemNumber); i++ {
			clValueK, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeMap.CLTypeKey, rem)
			if err != nil {
				return CLValue{}, nil, err
			}

			clValueV, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeMap.CLTypeValue, rem)
			if err != nil {
				return CLValue{}, nil, err
			}
			mapData[&clValueK] = clValueV
		}
	case CLTypeIDList:
		if len(data) < 4 {
			return CLValue{}, nil, newInvalidLengthErr(clType.CLTypeID)
		}

		if clType.CLTypeList == nil {
			return CLValue{}, nil, errors.New("expected not nil ClMap")
		}

		itemNumber := binary.LittleEndian.Uint32(data)
		var (
			//without uint32
			rem          = data[4:]
			clValueInner CLValue
			err          error
		)

		listData := make([]CLValue, 0, itemNumber)
		for i := 0; i < int(itemNumber); i++ {
			clValueInner, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeList.CLTypeInner, rem)
			if err != nil {
				return CLValue{}, nil, err
			}

			listData = append(listData, clValueInner)
		}
		return CLValue{
			Type: clType,
			List: &listData,
		}, rem, nil

	case CLTypeIDTuple2:
		if clType.CLTypeTuple2 == nil {
			return CLValue{}, nil, errors.New("expected not nil CLTypeTuple2")
		}

		var (
			rem                    = data
			clValueEl1, clValueEl2 CLValue
			err                    error
		)

		clValueEl1, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeTuple2.CLTypeElement1, rem)
		if err != nil {
			return CLValue{}, nil, err
		}

		clValueEl2, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeTuple2.CLTypeElement2, rem)
		if err != nil {
			return CLValue{}, nil, err
		}

		return CLValue{
			Type:   clType,
			Tuple2: &[2]CLValue{clValueEl1, clValueEl2},
		}, rem, nil

	case CLTypeIDTuple3:
		if clType.CLTypeTuple3 == nil {
			return CLValue{}, nil, errors.New("expected not nil CLTypeTuple3")
		}

		var (
			rem                                = data
			clValueEl1, clValueEl2, clValueEl3 CLValue
			err                                error
		)

		clValueEl1, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeTuple3.CLTypeElement1, rem)
		if err != nil {
			return CLValue{}, nil, err
		}

		clValueEl2, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeTuple3.CLTypeElement2, rem)
		if err != nil {
			return CLValue{}, nil, err
		}

		clValueEl3, rem, err = NewCLValueFromBytesWithRemainder(clType.CLTypeTuple3.CLTypeElement3, rem)
		if err != nil {
			return CLValue{}, nil, err
		}

		return CLValue{
			Type:   clType,
			Tuple3: &[3]CLValue{clValueEl1, clValueEl2, clValueEl3},
		}, rem, nil

	case CLTypeIDOption:
		if remainder[0] != 0 && clType.CLTypeOption != nil {
			var clValue CLValue
			clValue, remainder, err = NewCLValueFromBytesWithRemainder(clType.CLTypeOption.CLTypeInner, remainder[1:])
			if err != nil {
				return CLValue{}, nil, err
			}
			return CLValue{
				Type:   clType,
				Option: &clValue,
			}, remainder, nil
		}
		remainder = remainder[1:]

		return CLValue{
			Type: clType,
		}, remainder, nil
	}

	return CLValue{}, remainder, errors.New("unknown CLType provided")
}

func ParseCLValueFromBytesWithRemainder(data string) (RawCLValue, []byte, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	bytes, remainder, err := ParseBytesWithRemainder(decoded)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	clType, remainder, err := ClTypeFromBytes(0, remainder)
	if err != nil {
		return RawCLValue{}, nil, err
	}

	return RawCLValue{
		Type:  clType,
		Bytes: bytes,
	}, remainder, nil
}

func newInvalidLengthErr(clTypeID CLTypeID) error {
	return fmt.Errorf("invalid bytes length value for type - %s", clTypeID.ToString())
}
