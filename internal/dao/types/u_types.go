package types

import (
	"encoding/binary"
	"errors"
	"math/big"
	"sort"
)

// TODO: remove it and use SDK types
type UType interface {
	U128 | U256 | U512
}

type (
	U128 big.Int
	U256 big.Int
	U512 big.Int
)

func ParseUTypeFromBytes[T UType](bytes []byte) (T, []byte, error) {
	if len(bytes) == 0 {
		return T{}, nil, errors.New("empty bytes provided")
	}

	// read first bytes as bytes number
	numBytes := bytes[0]
	if int(numBytes) > len(bytes) {
		return T{}, nil, errors.New("invalid bytes format: number_bytes is more than bytes slice")
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

	return T(val), rem, nil
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

func (u U128) Into() *big.Int {
	val := big.Int(u)
	return &val
}

func (u U256) Into() *big.Int {
	val := big.Int(u)
	return &val
}

func (u U512) Into() *big.Int {
	val := big.Int(u)
	return &val
}
