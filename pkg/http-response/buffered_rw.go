package http_response

import (
	"bytes"
	"net/http"
)

// BufferedResponseWriter implementation of http.ResponseWriter interface
// used to save produced by handler response into the buffer
type BufferedResponseWriter struct {
	StatusCode int
	Buf        *bytes.Buffer
}

func (w *BufferedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *BufferedResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
}

func (w *BufferedResponseWriter) Write(p []byte) (int, error) {
	return w.Buf.Write(p)
}
