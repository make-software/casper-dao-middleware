package http_response

import (
	"encoding/json"
	"net/http"

	"casper-dao-middleware/pkg/errors"
	"casper-dao-middleware/pkg/pagination"

	"go.uber.org/zap"
)

type ErrorResult struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

type ErrorResponse struct {
	Error ErrorResult `json:"error,omitempty"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type PaginatedResponse struct {
	SuccessResponse
	ItemCount uint64 `json:"item_count"`
	PageCount uint64 `json:"page_count"`
}

func NewPaginatedResponseFromResult(result *pagination.Result) PaginatedResponse {
	return PaginatedResponse{
		SuccessResponse: SuccessResponse{
			Data: result.Data,
		},
		ItemCount: result.ItemCount,
		PageCount: result.PageCount,
	}
}

func Success(w http.ResponseWriter, data interface{}) {
	if result, ok := data.(*pagination.Result); ok {
		WriteJSON(w, http.StatusOK, NewPaginatedResponseFromResult(result))
		return
	}

	WriteJSON(w, http.StatusOK, SuccessResponse{
		Data: data,
	})
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	var (
		response ErrorResponse
		httpCode int
	)
	switch castedErr := err.(type) {
	case errors.Error:
		httpCode = castedErr.GetHTTPCode()
		response.Error = ErrorResult{
			Code:        castedErr.GetCode(),
			Message:     castedErr.GetMessage(),
			Description: castedErr.GetDescription(),
		}
	default:
		httpCode = http.StatusInternalServerError
		response.Error = ErrorResult{
			Code: "internal_error",
		}
		zap.S().With(zap.Error(err)).Errorw("unhandled server error", "urlPath", r.URL.Path)
	}

	WriteJSON(w, httpCode, response)
}

func FromFunction[T any, F func() (T, error)](function F, w http.ResponseWriter, r *http.Request) {
	result, err := function()
	if err != nil {
		Error(w, r, err)
		return
	}

	Success(w, result)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error while render response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(marshaled)
}
