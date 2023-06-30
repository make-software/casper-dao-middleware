package types

import (
	"encoding/binary"
	"strconv"
)

type RecordValue struct {
	U64Value  *uint64
	UValue    *U256
	Address   *Address
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

func NewRecordValueFromBytesWithReminder(rawBytes []byte) (RecordValue, []byte, error) {
	numBytes := binary.LittleEndian.Uint32(rawBytes)
	// shift 4 bytes (uint32)
	reminder := rawBytes[4:]

	// length 33 represent Key CLValue
	if numBytes == 33 {
		key, reminder, err := ParseKeyFromBytes(reminder)
		if err != nil {
			return RecordValue{}, nil, err
		}

		var address Address
		if key.AccountHash != nil {
			address.AccountHash = key.AccountHash
		} else {
			address.ContractPackageHash = key.Hash
		}

		return RecordValue{
			Address: &address,
		}, reminder, nil
	}

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

	val, _, err := ParseUTypeFromBytes[U256](reminder)
	if err != nil {
		return RecordValue{}, nil, err
	}

	return RecordValue{
		UValue: &val,
	}, reminder[numBytes:], nil
}

func NewRecordValueFromBytes(rawBytes []byte) (RecordValue, error) {
	numBytes := len(rawBytes)

	// length 33 represent Key CLValue
	if numBytes == 33 {
		key, _, err := ParseKeyFromBytes(rawBytes)
		if err != nil {
			return RecordValue{}, err
		}

		var address Address
		if key.AccountHash != nil {
			address.AccountHash = key.AccountHash
		} else {
			address.ContractPackageHash = key.Hash
		}

		return RecordValue{
			Address: &address,
		}, nil
	}

	// numBytes == 8 could be or pure u64 or U256/U512 coded in 8 bytes
	// but U256/U512 bytes representation is not equal to u64
	// so if the numBytes == 8 and U256/U512 it should be the following bytes representation:
	// 8 0 0 0  ==> numBytes + 7  ==> numBytes of internal data + internal data(7 bytes)
	// U256/U512 =  8 0 0 0 7 1 1 1 1 1 1 1
	// u64 =  8 0 0 0 1 1 1 1 1 1 1 1
	if numBytes == 8 && rawBytes[0] != 7 {
		val := binary.LittleEndian.Uint64(rawBytes)
		return RecordValue{
			U64Value: &val,
		}, nil
	}

	if numBytes == 1 {
		boolVal := rawBytes[0] == 1
		return RecordValue{
			BoolValue: &boolVal,
		}, nil
	}

	val, _, err := ParseUTypeFromBytes[U256](rawBytes)
	if err != nil {
		return RecordValue{}, err
	}

	return RecordValue{
		UValue: &val,
	}, nil
}

func (r RecordValue) String() string {
	switch {
	case r.UValue != nil:
		val := *r.UValue
		return val.Into().String()
	case r.U64Value != nil:
		return strconv.FormatInt(int64(*r.U64Value), 10)
	case r.BoolValue != nil:
		return strconv.FormatBool(*r.BoolValue)
	case r.Address != nil:
		if r.Address.ContractPackageHash != nil {
			return r.Address.ContractPackageHash.ToHex()
		}
		if r.Address.AccountHash != nil {
			return r.Address.AccountHash.ToHex()
		}
	}

	return ""
}
