package types

import (
	"encoding/binary"
	"errors"
	"math/big"
	"sort"
)

type U256 *big.Int

func NewU256FromBytes(val []byte) U256 {
	var result big.Int
	result.SetBytes(val)

	return &result
}

func NewU256FromUint64(val uint64) U256 {
	var result big.Int
	result.SetUint64(val)

	return &result
}

func ParseU256FromBytes(bytes []byte) (U256, []byte, error) {
	if len(bytes) == 0 {
		return nil, nil, errors.New("empty bytes provided")
	}

	// read first bytes as bytes number
	numBytes := bytes[0]
	if int(numBytes) > len(bytes) {
		return nil, nil, errors.New("invalid bytes format: number_bytes is more than bytes slice")
	}

	remainder := bytes[1:]

	value := make([]byte, numBytes)
	copy(value, remainder[:numBytes])

	// here is the tricky part: bytes are coming in little-endian but big.Int.SetBytes expect big-endian
	sort.SliceStable(value, func(i, j int) bool { return i > j })

	var val big.Int
	val.SetBytes(value[:])

	rem := make([]byte, len(remainder)-int(numBytes))
	copy(rem, remainder[numBytes:])

	return &val, rem, nil
}

func ParseStringFromBytes(bytes []byte) (string, []byte, error) {
	if len(bytes) == 0 {
		return "", nil, errors.New("empty bytes provided")
	}

	numBytes := binary.LittleEndian.Uint32(bytes)
	if int(numBytes) > len(bytes) {
		return "nil", nil, errors.New("invalid bytes format: number_bytes is more than bytes slice")
	}
	// shift to 4 bytes (unit32)
	remainder := bytes[4:]

	value := make([]byte, numBytes)
	copy(value, remainder[:numBytes])

	remainder = remainder[numBytes:]
	return string(value), remainder, nil
}
