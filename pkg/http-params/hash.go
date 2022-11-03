package http_params

import (
	"fmt"
	"net/http"

	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func ParseHash(key string, r *http.Request) (types.Hash, error) {
	rawHash := chi.URLParam(r, key)
	if rawHash == "" {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	hash, err := types.NewHashFromHexString(rawHash)
	if err != nil {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}

	return hash, nil
}

func ParseOptionalHash(key string, r *http.Request) (*types.Hash, error) {
	hash, _, err := parseHashByKey(key, r)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ParseOptionalHashList(key string, r *http.Request) ([]types.Hash, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, nil
	}

	if !commaListRegexp.MatchString(param) {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	list := commaListRegexp.FindAllString(param, -1)
	res := make([]types.Hash, 0, len(list))
	for _, key := range list {
		pubKey, err := types.NewHashFromHexString(key)
		if err != nil {
			return nil, ErrFailedToParsePublicKey
		}
		res = append(res, pubKey)
	}

	return res, nil
}

func parseHashByKey(key string, r *http.Request) (*types.Hash, bool, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, false, nil
	}

	hash, err := types.NewHashFromHexString(param)
	if err != nil {
		return nil, true, errors.NewInvalidInputError(fmt.Sprintf("failed to parse %s parameter", key))
	}
	return &hash, true, nil
}
