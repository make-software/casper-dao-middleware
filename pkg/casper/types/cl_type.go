package types

import (
	"encoding/binary"
	"errors"
)

const CLTypeRecursionDepth = 50

const (
	CLTypeIDBool CLTypeID = iota
	CLTypeIDI32
	CLTypeIDI64
	CLTypeIDU8
	CLTypeIDU32
	CLTypeIDU64
	CLTypeIDU128
	CLTypeIDU256
	CLTypeIDU512
	CLTypeIDUnit
	CLTypeIDString
	CLTypeIDKey
	CLTypeIDURef
	CLTypeIDOption
	CLTypeIDList
	CLTypeIDByteArray
	CLTypeIDResult
	CLTypeIDMap
	CLTypeIDTuple1
	CLTypeIDTuple2
	CLTypeIDTuple3
	CLTypeIDAny
	CLTypeIDPublicKey
)

var (
	ErrExceededRecursionDepth = errors.New("recursion depth exceeded during CLType parsing")
	ErrEmptyCLTypeBytes       = errors.New("empty CLType bytes provided")
)

type (
	СLTypeMap struct {
		CLTypeKey   CLType
		CLTypeValue CLType
	}

	CLTypeTuple1 struct {
		CLTypeElement CLType
	}
	CLTypeTuple2 struct {
		CLTypeElement1 CLType
		CLTypeElement2 CLType
	}
	CLTypeTuple3 struct {
		CLTypeElement1 CLType
		CLTypeElement2 CLType
		CLTypeElement3 CLType
	}
	CLTypeList struct {
		CLTypeInner CLType
	}
	CLTypeOption struct {
		CLTypeInner CLType
	}
	CLTypeResult struct {
		CLTypeOk  CLType
		CLTypeErr CLType
	}
)

type (
	CLTypeID byte

	CLType struct {
		CLTypeID     CLTypeID
		CLTypeMap    *СLTypeMap
		CLTypeTuple1 *CLTypeTuple1
		CLTypeTuple2 *CLTypeTuple2
		CLTypeTuple3 *CLTypeTuple3
		CLTypeList   *CLTypeList
		CLTypeOption *CLTypeOption
		CLTypeResult *CLTypeResult
	}
)

func NewCLType(clTypeID CLTypeID) CLType {
	return CLType{
		CLTypeID: clTypeID,
	}
}

func (t CLTypeID) ToString() string {
	switch t {
	case CLTypeIDBool:
		return "Bool"
	case CLTypeIDI32:
		return "I32"
	case CLTypeIDI64:
		return "I64"
	case CLTypeIDU8:
		return "U8"
	case CLTypeIDU32:
		return "U32"
	case CLTypeIDU64:
		return "U64"
	case CLTypeIDU128:
		return "U128"
	case CLTypeIDU256:
		return "U256"
	case CLTypeIDU512:
		return "U512"
	case CLTypeIDUnit:
		return "Unit"
	case CLTypeIDString:
		return "String"
	case CLTypeIDKey:
		return "Key"
	case CLTypeIDURef:
		return "URef"
	case CLTypeIDOption:
		return "Option"
	case CLTypeIDList:
		return "List"
	case CLTypeIDByteArray:
		return "ByteArray"
	case CLTypeIDResult:
		return "Result"
	case CLTypeIDMap:
		return "Map"
	case CLTypeIDTuple1:
		return "Tuple1"
	case CLTypeIDTuple2:
		return "Tuple2"
	case CLTypeIDTuple3:
		return "Tuple3"
	case CLTypeIDPublicKey:
		return "PublicKey"
	}

	return "Any"
}

func ClTypeFromBytes(depth uint8, bytes []byte) (CLType, []byte, error) {
	if len(bytes) == 0 {
		return CLType{}, nil, ErrEmptyCLTypeBytes
	}

	if depth >= CLTypeRecursionDepth {
		return CLType{}, nil, ErrExceededRecursionDepth
	}

	depth++
	remainder := make([]byte, len(bytes)-1)
	copy(remainder, bytes[1:])

	switch CLTypeID(bytes[0]) {
	case CLTypeIDBool:
		return NewCLType(CLTypeIDBool), remainder, nil
	case CLTypeIDI32:
		return NewCLType(CLTypeIDI32), remainder, nil
	case CLTypeIDI64:
		return NewCLType(CLTypeIDI64), remainder, nil
	case CLTypeIDU8:
		return NewCLType(CLTypeIDU8), remainder, nil
	case CLTypeIDU32:
		return NewCLType(CLTypeIDU32), remainder, nil
	case CLTypeIDU64:
		return NewCLType(CLTypeIDU64), remainder, nil
	case CLTypeIDU128:
		return NewCLType(CLTypeIDU128), remainder, nil
	case CLTypeIDU256:
		return NewCLType(CLTypeIDU256), remainder, nil
	case CLTypeIDU512:
		return NewCLType(CLTypeIDU512), remainder, nil
	case CLTypeIDUnit:
		return NewCLType(CLTypeIDUnit), remainder, nil
	case CLTypeIDString:
		return NewCLType(CLTypeIDString), remainder, nil
	case CLTypeIDKey:
		return NewCLType(CLTypeIDKey), remainder, nil
	case CLTypeIDURef:
		return NewCLType(CLTypeIDURef), remainder, nil
	case CLTypeIDOption:
		innerType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeIDOption)
		optionCLType := CLTypeOption{
			CLTypeInner: innerType,
		}
		clValue.CLTypeOption = &optionCLType
		return clValue, remainder, nil
	case CLTypeIDList:
		innerType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeIDList)
		listCLType := CLTypeList{
			CLTypeInner: innerType,
		}
		clValue.CLTypeList = &listCLType
		return clValue, remainder, nil
	case CLTypeIDByteArray:
	case CLTypeIDResult:
	case CLTypeIDMap:
		keyType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}

		valueType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeIDMap)
		mapCLType := СLTypeMap{
			CLTypeKey:   keyType,
			CLTypeValue: valueType,
		}

		clValue.CLTypeMap = &mapCLType
		return clValue, remainder, nil
	case CLTypeIDTuple1:
		element1Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}

		clValue := NewCLType(CLTypeIDTuple1)
		tuple1CLType := CLTypeTuple1{
			CLTypeElement: element1Type,
		}

		clValue.CLTypeTuple1 = &tuple1CLType
		return clValue, remainder, nil
	case CLTypeIDTuple2:
		element1Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}

		element2Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeIDTuple2)
		tuple2CLType := CLTypeTuple2{
			CLTypeElement1: element1Type,
			CLTypeElement2: element2Type,
		}

		clValue.CLTypeTuple2 = &tuple2CLType
		return clValue, remainder, nil
	case CLTypeIDTuple3:
		element1Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}

		element2Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}

		element3Type, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeIDTuple3)
		tuple3CLType := CLTypeTuple3{
			CLTypeElement1: element1Type,
			CLTypeElement2: element2Type,
			CLTypeElement3: element3Type,
		}

		clValue.CLTypeTuple3 = &tuple3CLType
		return clValue, remainder, nil

	case CLTypeIDAny:
		return NewCLType(CLTypeIDAny), remainder, nil
	case CLTypeIDPublicKey:
	}

	return CLType{}, nil, errors.New("invalid CLType provided")
}

// ParseBytesWithReminder looks first bytes to detect length of bytes, extract it and return reminder
func ParseBytesWithReminder(data []byte) ([]byte, []byte, error) {
	if len(data) < 4 {
		return nil, nil, errors.New("invalid length value")
	}

	length := binary.LittleEndian.Uint32(data)
	if length == 0 || int(length) > len(data) {
		return nil, nil, errors.New("invalid length value")
	}

	// without uint32
	data = data[4:]
	if int(length) > len(data) {
		return nil, nil, errors.New("invalid length value")
	}

	bytes := make([]byte, length)
	copy(bytes, data[:length])

	reminder := make([]byte, len(data)-int(length))
	copy(reminder, data[length:])
	return bytes, reminder, nil
}
