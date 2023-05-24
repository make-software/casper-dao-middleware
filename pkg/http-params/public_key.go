package http_params

import (
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/pkg/errors"
)

var (
	ErrFailedToParsePublicKey = errors.NewInvalidInputError("failed to parse `public_key` parameter")
	ErrPublicKeyIsNotProvided = errors.NewInvalidInputError("`public_key` was not provided")
)

func parsePublicKeyByKey(key string, r *http.Request) (*casper.PublicKey, bool, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, false, nil
	}

	pubKey, err := casper.NewPublicKey(param)
	if err != nil {
		return nil, true, ErrFailedToParsePublicKey
	}
	return &pubKey, true, nil
}

func ParsePublicKey(key string, request *http.Request) (casper.PublicKey, error) {
	pubKey, isProvided, err := parsePublicKeyByKey(key, request)
	if err != nil {
		return casper.PublicKey{}, err
	}

	if !isProvided {
		return casper.PublicKey{}, ErrPublicKeyIsNotProvided
	}

	return *pubKey, nil
}

func ParseOptionalPublicKey(key string, request *http.Request) (*casper.PublicKey, error) {
	pubKey, _, err := parsePublicKeyByKey(key, request)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func ParseOptionalPublicKeyList(key string, r *http.Request) ([]casper.PublicKey, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, nil
	}

	if !commaListRegexp.MatchString(param) {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	list := commaListRegexp.FindAllString(param, -1)
	res := make([]casper.PublicKey, 0, len(list))
	for _, key := range list {
		pubKey, err := casper.NewPublicKey(key)
		if err != nil {
			return nil, ErrFailedToParsePublicKey
		}
		res = append(res, pubKey)
	}

	return res, nil
}
