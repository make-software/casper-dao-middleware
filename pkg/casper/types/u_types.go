package types

import (
	"errors"
	"math/big"
	"sort"
)

type U256 *big.Int

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
