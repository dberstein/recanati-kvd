package rw

import "net/http"

type ResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
}

func New(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 0}
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}

func (w *ResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	// receive status code from this method
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)

	return
}
