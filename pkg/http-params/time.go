package http_params

import (
	"fmt"
	"net/http"
	"time"

	"casper-dao-middleware/pkg/errors"
	"casper-dao-middleware/pkg/types"

	"github.com/go-chi/chi/v5"
)

func ParseTime(key string, r *http.Request) (time.Time, error) {
	rawTime := chi.URLParam(r, key)
	if rawTime == "" {
		return time.Time{}, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	parsed, err := time.Parse(time.RFC3339, rawTime)
	if err != nil {
		return time.Time{}, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format should be %s", key, time.RFC3339))
	}

	return parsed, nil
}

func ParseDate(key string, r *http.Request) (types.Date, error) {
	rawDate := chi.URLParam(r, key)
	if rawDate == "" {
		return types.Date{}, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	parsed, err := types.ParseDateFromString(rawDate)
	if err != nil {
		return types.Date{}, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format should be %s", key, "2014-04-26"))
	}

	return parsed, nil
}

func ParseOptionalDateList(key string, r *http.Request) ([]types.Date, error) {
	params, _, err := ParseOptionalCommaSeparatedList(key, r)
	if err != nil {
		return nil, err
	}

	res := make([]types.Date, 0, len(params))
	for _, param := range params {
		parsed, err := types.ParseDateFromString(param)
		if err != nil {
			//TODO: check exact error
			return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format should be %s", key, "2014-04-26"))
		}
		res = append(res, parsed)
	}
	return res, nil
}

func ParseOptionalTimeList(key string, r *http.Request) ([]time.Time, error) {
	params, _, err := ParseOptionalCommaSeparatedList(key, r)
	if err != nil {
		return nil, err
	}

	res := make([]time.Time, 0, len(params))
	for _, param := range params {
		parsed, err := time.Parse(time.RFC3339, param)
		if err != nil {
			return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` format should be %s", key, time.RFC3339))
		}
		res = append(res, parsed)
	}
	return res, nil
}
