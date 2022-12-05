package utils

import (
	"casper-dao-middleware/pkg/casper/types"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"strconv"
	"strings"
)

type RecordValue struct {
	U64Value  *uint64
	UValue    *types.U256
	BoolValue *bool
}

type FutureValue struct {
	Value          RecordValue
	ActivationTime uint64
}

type Record struct {
	Value       RecordValue
	FutureValue *FutureValue
}

func NewRecordFromBytes(rawBytes []byte) (Record, error) {
	recordValue, reminder, err := NewRecordValueFromBytesWithReminder(rawBytes)
	if err != nil {
		return Record{}, nil
	}

	record := Record{
		Value: recordValue,
	}

	if len(reminder) == 0 {
		return Record{}, errors.New("invalid record format")
	}

	if reminder[0] == 1 {
		recordValue, reminder, err = NewRecordValueFromBytesWithReminder(reminder[1:])
		if err != nil {
			return Record{}, nil
		}

		activationTime := binary.LittleEndian.Uint64(reminder)

		record.FutureValue = &FutureValue{
			Value:          recordValue,
			ActivationTime: activationTime,
		}
	}
	return record, nil
}

func ToDictionaryItemKey(key string) (string, error) {
	res := make([]byte, 0)
	blake, err := blake2b.New256(res)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(key)
	blake.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(keyBytes))))
	blake.Write(keyBytes)

	return hex.EncodeToString(blake.Sum(nil)), nil
}

func ToDictionaryKey(eventsUref string, index uint32) (string, error) {
	urefsParts := strings.Split(eventsUref, "-")
	// uref format uref-d1a68e4ae2c8ffe65cafcfc172caf1179bc5fa820d25eb4574a48f89225820a0-007
	if len(urefsParts) != 3 {
		return "", errors.New("invalid uref format provided")
	}
	urefHashBytes, err := hex.DecodeString(urefsParts[1])
	if err != nil {
		return "", err
	}

	res := make([]byte, 0)
	key, err := blake2b.New256(res)
	if err != nil {
		return "", err
	}

	key.Write(urefHashBytes)
	key.Write(calculateDictionaryIndexHash(index))
	dictionaryKey := fmt.Sprintf("dictionary-%s", hex.EncodeToString(key.Sum(nil)))
	return dictionaryKey, nil
}

func ParseDAOCLValueFromBytes(data string) (types.CLValue, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return types.CLValue{}, err
	}

	bytes, reminder, err := types.ParseBytesWithReminder(decoded)
	if err != nil {
		return types.CLValue{}, err
	}

	clType, reminder, err := types.ClTypeFromBytes(0, reminder)
	if err != nil {
		return types.CLValue{}, err
	}

	return types.CLValue{
		Type:  clType,
		Bytes: bytes,
	}, nil
}

func NewRecordValueFromBytesWithReminder(rawBytes []byte) (RecordValue, []byte, error) {
	numBytes := binary.LittleEndian.Uint32(rawBytes)
	// shift 4 bytes (uint32)
	reminder := rawBytes[4:]

	// numBytes == 8 could be or pure u64 or U256/U512 coded in 8 bytes
	// but U256/U512 bytes representation is not equal to u64
	// so if the numBytes == 8 and U256/U512 it should be the following bytes representation:
	// 8 0 0 0  ==> numBytes + 7  ==> numBytes of internal data + internal data(7 bytes)
	// U256/U512 =  8 0 0 0 7 1 1 1 1 1 1 1
	// u64 =  8 0 0 0 1 1 1 1 1 1 1 1
	if numBytes == 8 && reminder[0] != 7 {
		val := binary.LittleEndian.Uint64(reminder)
		return RecordValue{
			U64Value: &val,
		}, reminder[8:], nil
	}

	if numBytes == 1 {
		boolVal := reminder[0] == 1
		return RecordValue{
			BoolValue: &boolVal,
		}, reminder[1:], nil
	}

	val, _, err := types.ParseU256FromBytes(reminder)
	if err != nil {
		return RecordValue{}, nil, err
	}

	return RecordValue{
		UValue: &val,
	}, reminder[numBytes:], nil
}

func (r RecordValue) String() string {
	switch {
	case r.UValue != nil:
		val := *r.UValue
		return (*val).String()
	case r.U64Value != nil:
		return strconv.FormatInt(int64(*r.U64Value), 10)
	case r.BoolValue != nil:
		return strconv.FormatBool(*r.BoolValue)
	}

	return ""
}

func calculateDictionaryIndexHash(index uint32) []byte {
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, index)

	indexHash := blake2b.Sum256(indexBytes)
	return []byte(hex.EncodeToString(indexHash[:]))
}
