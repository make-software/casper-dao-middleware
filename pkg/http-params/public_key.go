package http_params

import (
	"fmt"
	"net/http"

	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/errors"
)

var (
	ErrFailedToParsePublicKey = errors.NewInvalidInputError("failed to parse `public_key` parameter")
	ErrPublicKeyIsNotProvided = errors.NewInvalidInputError("`public_key` was not provided")
)

func parsePublicKeyByKey(key string, r *http.Request) (*types.PublicKey, bool, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, false, nil
	}

	pubKey, err := types.NewPublicKeyFromHexString(param)
	if err != nil {
		return nil, true, ErrFailedToParsePublicKey
	}
	return &pubKey, true, nil
}

func ParsePublicKey(key string, request *http.Request) (types.PublicKey, error) {
	pubKey, isProvided, err := parsePublicKeyByKey(key, request)
	if err != nil {
		return types.PublicKey{}, err
	}

	if !isProvided {
		return types.PublicKey{}, ErrPublicKeyIsNotProvided
	}

	return *pubKey, nil
}

func ParseOptionalPublicKey(key string, request *http.Request) (*types.PublicKey, error) {
	pubKey, _, err := parsePublicKeyByKey(key, request)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func ParseOptionalPublicKeyList(key string, r *http.Request) ([]types.PublicKey, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, nil
	}

	if !commaListRegexp.MatchString(param) {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	list := commaListRegexp.FindAllString(param, -1)
	res := make([]types.PublicKey, 0, len(list))
	for _, key := range list {
		pubKey, err := types.NewPublicKeyFromHexString(key)
		if err != nil {
			return nil, ErrFailedToParsePublicKey
		}
		res = append(res, pubKey)
	}

	return res, nil
}
