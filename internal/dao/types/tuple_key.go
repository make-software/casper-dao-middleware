package types

import (
	"errors"

	casper_types "casper-dao-middleware/pkg/casper/types"
)

type Tuple2 struct {
	Element1 string
	Element2 uint32
}

func ParseTuple2U512MapFromCLValue(clValue casper_types.CLValue) (map[Tuple2]casper_types.U512, error) {
	result := make(map[Tuple2]casper_types.U512, len(clValue.Map.Data))
	for mapKey, mapVal := range clValue.Map.Data {
		if mapVal.U512 == nil {
			return nil, errors.New("expect not nil U512 value in map")
		}

		if mapKey.Tuple2 == nil {
			return nil, errors.New("expect not nil Tuple2 key in map")
		}

		if mapKey.Tuple2[0].Key == nil {
			return nil, errors.New("expect Key element1 in Tuple2 key in map")
		}

		if mapKey.Tuple2[1].U32 == nil {
			return nil, errors.New("expect U32 element2 in Tuple2 key in map")
		}

		var el1 string
		if mapKey.Tuple2[0].Key.AccountHash != nil {
			el1 = mapKey.Tuple2[0].Key.AccountHash.ToHex()
		} else {
			el1 = mapKey.Tuple2[0].Key.Hash.ToHex()
		}

		tupleKey := Tuple2{
			Element1: el1,
			Element2: *mapKey.Tuple2[1].U32,
		}

		result[tupleKey] = *mapVal.U512
	}

	return result, nil
}
