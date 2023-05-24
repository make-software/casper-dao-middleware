package types

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
)

type Tuple2 struct {
	Element1 string
	Element2 uint32
}

func ParseTuple2U512MapFromCLValue(clValue clvalue.CLValue) (map[Tuple2]clvalue.UInt512, error) {
	result := make(map[Tuple2]clvalue.UInt512, len(clValue.Map.Map()))
	for _, mapVal := range clValue.Map.Data() {

		if mapVal.Inner1.Tuple2 == nil {
			return nil, errors.New("expect Tuple2 key in map")
		}

		keyValue := mapVal.Inner1.Tuple2

		if keyValue.Inner1.Key == nil {
			return nil, errors.New("expect Key element1 in Tuple2 key in map")
		}

		if keyValue.Inner2.UI32 == nil {
			return nil, errors.New("expect U32 element2 in Tuple2 key in map")
		}

		if mapVal.Inner2.UI512 == nil {
			return nil, errors.New("expect UI512 element2 in map value")
		}

		var el1 string
		if keyValue.Inner1.Key.Account != nil {
			el1 = keyValue.Inner1.Key.Account.ToHex()
		} else {
			el1 = keyValue.Inner1.Key.Hash.ToHex()
		}

		tupleKey := Tuple2{
			Element1: el1,
			Element2: keyValue.Inner2.UI32.Value(),
		}

		result[tupleKey] = *mapVal.Inner2.UI512
	}

	return result, nil
}
