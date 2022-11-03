package http_params

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		setTime := time.Now().UTC()
		routeCtx.URLParams.Add("date", setTime.Format(time.RFC3339))
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		actualTime, err := ParseTime("date", req)
		assert.NoError(t, err)
		assert.Equal(t, actualTime.Format(time.RFC3339), setTime.Format(time.RFC3339))
	})

	t.Run("Fail: no `date` value", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseTime("date", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Empty `date` value")
	})

	t.Run("Fail: invalid `date` format", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("date", "invalid-format")
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseTime("date", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid `date` format")
	})
}

func TestParseOptionalTimeList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?date=2023-05-21T00:00:00Z,2025-05-21T00:00:00Z", nil)
		assert.NoError(t, err)

		timeList, err := ParseOptionalTimeList("date", req)
		assert.NoError(t, err)
		assert.Equal(t, len(timeList), 2)
	})

	t.Run("Empty param value", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?date=", nil)
		assert.NoError(t, err)

		_, err = ParseOptionalTimeList("date", req)
		assert.Equal(t, err.Error(), "Empty `date` value")
	})

	t.Run("Invalid time format", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/tests?date=invalid-time", nil)
		assert.NoError(t, err)

		_, err = ParseOptionalTimeList("date", req)
		assert.Contains(t, err.Error(), "Invalid `date` format should")
	})
}

func TestParseDate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		setTime := time.Now().UTC()
		routeCtx.URLParams.Add("date", setTime.Format(time.RFC3339))
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseDate("date", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Invalid `date` format should be 2014-04-26")
	})

	t.Run("Fail: no `date` value", func(t *testing.T) {
		routeCtx := chi.NewRouteContext()
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
		assert.NoError(t, err)

		_, err = ParseDate("date", req)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "Empty `date` value")
	})

	t.Run("Fail: invalid `date` format", func(t *testing.T) {
		formats := []string{"invalid", "2006-01-02T00:00:00Z", "2014-04"}
		for _, format := range formats {
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("date", format)
			ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/tests", nil)
			assert.NoError(t, err)

			_, err = ParseDate("date", req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Invalid `date` format")
		}
	})
}
