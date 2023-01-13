package types

import (
	"encoding/binary"
	"errors"
	"strconv"

	"casper-dao-middleware/pkg/casper/types"
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
