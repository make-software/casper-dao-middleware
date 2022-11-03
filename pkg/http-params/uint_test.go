package http_params

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseUint64(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		setUint := 10
		routeCtx.URLParams.Add("param", strconv.Itoa(setUint))
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		actualUint, err := ParseUint64("param", req)
		assert.NoError(t, err)
		assert.Equal(t, actualUint, uint64(setUint))
	})

	t.Run("Fail: no `uint` value", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseUint64("uint", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Empty `uint` value")
	})

	t.Run("Fail: invalid `param` format", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("param", "invalid-format")
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseUint64("param", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid `param` format")
	})
}

func TestParseOptionalUint32List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?param=100,123", nil)
		assert.NoError(t, err)

		uintList, err := ParseOptionalUint32List("param", req)
		assert.NoError(t, err)
		assert.Equal(t, len(uintList), 2)
	})

	t.Run("Empty param value", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?param=", nil)
		assert.NoError(t, err)

		_, err = ParseOptionalUint32List("param", req)
		assert.Equal(t, err.Error(), "Empty `param` value")
	})

	t.Run("Invalid uint format", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?param=invalid", nil)
		assert.NoError(t, err)

		_, err = ParseOptionalUint32List("param", req)
		assert.Contains(t, err.Error(), "invalid syntax")
	})
}
