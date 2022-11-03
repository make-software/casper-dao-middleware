package types

import (
	"encoding/binary"
	"errors"
)

const CLTypeRecursionDepth = 50

const (
	CLTypeBool CLTypeID = iota
	CLTypeI32
	CLTypeI64
	CLTypeU8
	CLTypeU32
	CLTypeU64
	CLTypeU128
	CLTypeU256
	CLTypeU512
	CLTypeUnit
	CLTypeString
	CLTypeKey
	CLTypeURef
	CLTypeOption
	CLTypeList
	CLTypeByteArray
	CLTypeResult
	CLTypeMap
	CLTypeTuple1
	CLTypeTuple2
	CLTypeTuple3
	CLTypeAny
	CLTypePublicKey
)

var ErrExceededRecursionDepth = errors.New("recursion depth exceeded during CLType parsing")

type (
	CLTypeID byte

	CLType struct {
		CLType   *CLType
		CLTypeID CLTypeID
	}

	CLValue struct {
		Type  CLType
		Bytes []byte
	}
)

func NewCLType(clTypeID CLTypeID) CLType {
	return CLType{
		CLTypeID: clTypeID,
	}
}

func (cv CLType) ToString() string {
	if cv.CLType == nil {
		return cv.CLTypeID.ToString()
	}

	result := cv.CLTypeID.ToString()
	return concatClTypes(result, *cv.CLType)
}

func (t CLTypeID) ToString() string {
	switch t {
	case CLTypeBool:
		return "Bool"
	case CLTypeI32:
		return "I32"
	case CLTypeI64:
		return "I64"
	case CLTypeU8:
		return "U8"
	case CLTypeU32:
		return "U32"
	case CLTypeU64:
		return "U64"
	case CLTypeU128:
		return "U128"
	case CLTypeU256:
		return "U256"
	case CLTypeU512:
		return "U512"
	case CLTypeUnit:
		return "Unit"
	case CLTypeString:
		return "String"
	case CLTypeKey:
		return "Key"
	case CLTypeURef:
		return "URef"
	case CLTypeOption:
		return "Option"
	case CLTypeList:
		return "List"
	case CLTypeByteArray:
		return "ByteArray"
	case CLTypeResult:
		return "Result"
	case CLTypeMap:
		return "Map"
	case CLTypeTuple1:
		return "Tuple1"
	case CLTypeTuple2:
		return "Tuple2"
	case CLTypeTuple3:
		return "Tuple3"
	case CLTypePublicKey:
		return "PublicKey"
	}

	return "Any"
}

func ClTypeFromBytes(depth uint8, bytes []byte) (CLType, []byte, error) {
	if depth >= CLTypeRecursionDepth {
		return CLType{}, nil, ErrExceededRecursionDepth
	}

	depth = depth + 1
	remainder := make([]byte, len(bytes)-1)
	copy(remainder, bytes[1:])

	switch CLTypeID(bytes[0]) {
	case CLTypeBool:
		return NewCLType(CLTypeBool), remainder, nil
	case CLTypeI32:
		return NewCLType(CLTypeI32), remainder, nil
	case CLTypeI64:
		return NewCLType(CLTypeI64), remainder, nil
	case CLTypeU8:
		return NewCLType(CLTypeU8), remainder, nil
	case CLTypeU32:
		return NewCLType(CLTypeU32), remainder, nil
	case CLTypeU64:
		return NewCLType(CLTypeU64), remainder, nil
	case CLTypeU128:
		return NewCLType(CLTypeU128), remainder, nil
	case CLTypeU256:
		return NewCLType(CLTypeU256), remainder, nil
	case CLTypeU512:
		return NewCLType(CLTypeU512), remainder, nil
	case CLTypeUnit:
		return NewCLType(CLTypeUnit), remainder, nil
	case CLTypeString:
		return NewCLType(CLTypeString), remainder, nil
	case CLTypeKey:
		return NewCLType(CLTypeKey), remainder, nil
	case CLTypeURef:
		return NewCLType(CLTypeURef), remainder, nil
	case CLTypeOption:
		innerType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeOption)
		clValue.CLType = &innerType
		return clValue, remainder, nil
	case CLTypeList:
		innerType, remainder, err := ClTypeFromBytes(depth, remainder)
		if err != nil {
			return CLType{}, nil, err
		}
		clValue := NewCLType(CLTypeList)
		clValue.CLType = &innerType
		return clValue, remainder, nil
	case CLTypeByteArray:
	case CLTypeResult:
	case CLTypeMap:
	case CLTypeTuple1:
	case CLTypeTuple2:
	case CLTypeTuple3:
	case CLTypeAny:
	case CLTypePublicKey:

	}

	return CLType{}, nil, nil
}

// ParseBytesWithReminder looks first bytes to detect length of bytes, extract it and return reminder
func ParseBytesWithReminder(data []byte) ([]byte, []byte, error) {
	length := binary.LittleEndian.Uint32(data)
	if length == 0 || int(length) > len(data) {
		return nil, nil, errors.New("invalid length value")
	}

	// without uint32
	data = data[4:]

	bytes := make([]byte, length)
	copy(bytes, data[:length])

	reminder := make([]byte, len(data)-int(length))
	copy(reminder, data[length:])
	return bytes, reminder, nil
}

func concatClTypes(result string, clType CLType) string {
	result += "("
	if clType.CLType == nil {
		return "(" + clType.CLTypeID.ToString() + ")"
	}

	result += clType.CLTypeID.ToString()
	result += concatClTypes(result, *clType.CLType)
	result += ")"
	return result
}
