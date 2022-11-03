package http_params

import (
	"fmt"
	"net/http"

	"casper-dao-middleware/pkg/errors"
	"casper-dao-middleware/pkg/types"
)

func ParseOptionalData(key string, request *http.Request) (*types.OptionalData, error) {
	includesData, isProvided := getParamByKey(key, request)
	if !isProvided {
		return &types.OptionalData{}, nil
	}

	optionalData, err := types.ParseOptionalData(includesData)
	if err != nil {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid '%s' param", key)).
			SetDescription(err.Error())
	}

	return optionalData, nil
}
