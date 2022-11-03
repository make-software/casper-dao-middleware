package http_params

import (
	"fmt"
	"net/http"

	"casper-dao-middleware/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func ParseString(key string, r *http.Request) (string, error) {
	rawParam := chi.URLParam(r, key)
	if rawParam == "" {
		return "", errors.NewInvalidInputError(fmt.Sprintf("Empty `%s` value", key))
	}

	return rawParam, nil
}
