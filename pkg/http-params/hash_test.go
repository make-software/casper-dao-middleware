package http_params

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseHash(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()

		setHash := "9fccaf372ec5f61ac851fcec593d159f928a26df8f2af5aa3522ed9e0b7cbb36"
		routeCtx.URLParams.Add("hash", setHash)
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/hashes", nil)
		assert.NoError(t, err)

		actualHash, err := ParseHash("hash", req)
		assert.NoError(t, err)
		assert.Equal(t, actualHash.ToHex(), setHash)
	})

	t.Run("Fail: no `hash` value", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/hashes", nil)
		assert.NoError(t, err)

		_, err = ParseHash("hash", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Empty `hash` value")
	})

	t.Run("Fail: invalid `hash` format", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("hash", "invalid-format")
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/hashes", nil)
		assert.NoError(t, err)

		_, err = ParseUint64("hash", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid `hash` format")
	})
}
