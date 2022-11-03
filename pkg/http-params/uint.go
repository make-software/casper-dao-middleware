package http_params

import (
	"fmt"
	"net/http"
	"strconv"

	"casper-dao-middleware/pkg/errors"
)

func ParseUint64(key string, r *http.Request) (uint64, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return 0, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	param, err := strconv.ParseUint(stringParam, 10, 0)
	if err != nil {
		return 0, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}

	return param, nil
}

func ParseUint16(key string, r *http.Request) (uint16, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return 0, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	param, err := strconv.ParseUint(stringParam, 10, 0)
	if err != nil {
		return 0, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}
	return uint16(param), nil
}

func ParseOptionalUint16(key string, r *http.Request) (*uint16, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return nil, nil
	}

	param, err := strconv.ParseUint(stringParam, 10, 0)
	if err != nil {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}

	parsed := uint16(param)
	return &parsed, nil
}

func ParseOptionalUint32(key string, r *http.Request) (*uint32, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return nil, nil
	}

	param, err := strconv.ParseUint(stringParam, 10, 0)
	if err != nil {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}

	parsed := uint32(param)
	return &parsed, nil
}

func ParseOptionalUint16List(key string, r *http.Request) ([]uint16, error) {
	params, _, err := ParseOptionalCommaSeparatedList(key, r)
	if err != nil {
		return nil, err
	}

	result := make([]uint16, 0, len(params))
	for _, param := range params {
		parsed, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
		}
		result = append(result, uint16(parsed))
	}
	return result, nil
}

func ParseOptionalUint32List(key string, r *http.Request) ([]uint32, error) {
	params, _, err := ParseOptionalCommaSeparatedList(key, r)
	if err != nil {
		return nil, err
	}

	result := make([]uint32, 0, len(params))
	for _, param := range params {
		parsed, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, uint32(parsed))
	}
	return result, nil
}

func ParseUint32(key string, r *http.Request) (uint32, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return 0, nil
	}

	param, err := strconv.ParseUint(stringParam, 10, 0)
	if err != nil {
		return 0, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format", key))
	}

	parsed := uint32(param)
	return parsed, nil
}
