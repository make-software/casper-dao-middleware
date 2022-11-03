package http_params

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"casper-dao-middleware/pkg/errors"
	"github.com/go-chi/chi/v5"
)

var commaListRegexp = regexp.MustCompile("([^,]+)")

// getParamByKey generic read param by key as URL.Query chi.URLParam
func getParamByKey(key string, r *http.Request) (string, bool) {
	if hasKey := r.URL.Query().Has(key); hasKey {
		rawKey := r.URL.Query().Get(key)
		if rawKey != "" {
			return rawKey, true
		}
		return "", true
	}

	if chiParam := chi.URLParam(r, key); chiParam != "" {
		return chiParam, true
	}

	return "", false
}

func ParseOptionalCommaSeparatedList(key string, r *http.Request) ([]string, bool, error) {
	param, isProvided := getParamByKey(key, r)
	if !isProvided {
		return nil, false, nil
	}

	// we are using regexp instead of strings.Split because it is more flexibility
	// with strings.Split we should do additional checks as param=`value`,`` will be parsed as []string{"value",""}, etc
	if !commaListRegexp.MatchString(param) {
		return nil, false, errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	// find all occurrences
	return commaListRegexp.FindAllString(param, -1), true, nil
}

func ParseOptionalBool(key string, r *http.Request) (*bool, error) {
	stringParam, ok := getParamByKey(key, r)
	if !ok {
		return nil, nil
	}

	boolParam, err := strconv.ParseBool(stringParam)
	if err != nil {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("Invalid `%s` bool format", key))
	}

	return &boolParam, nil
}
