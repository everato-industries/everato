package utils

import (
	"encoding/json"
	"net/http"
)

type M map[string]any

// HttpWriter is a utility struct to handle HTTP responses
// it helps write response to the http stream more easily
type HttpWriter struct {
	W          http.ResponseWriter
	R          *http.Request
	StatusCode int
}

// Returns a new HttpWriter instance
func NewHttpWriter(w http.ResponseWriter, r *http.Request) *HttpWriter {
	return &HttpWriter{
		W:          w,
		R:          r,
		StatusCode: http.StatusOK, // Default status code
	}
}

// Status sets the HTTP status code for the response
func (hw *HttpWriter) Status(code int) *HttpWriter {
	hw.StatusCode = code
	hw.W.WriteHeader(code)
	return hw
}

// Json writes a JSON response to the HTTP response writer
func (hw *HttpWriter) Json(data M) {
	// Convert data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		hw.W.WriteHeader(http.StatusInternalServerError)
		hw.W.Header().Set("Content-Type", "text/plain")
		http.Error(hw.W, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set content type header
	hw.W.Header().Set("Content-Type", "application/json")

	// Write JSON data
	s, err := hw.W.Write(jsonData)
	if err != nil && s != len(jsonData) {
		hw.W.WriteHeader(http.StatusInternalServerError)
		http.Error(hw.W, "Failed to write JSON response", http.StatusInternalServerError)
		return
	}
}

// Text writes a TEXT response to the HTTP response writer
func (hw *HttpWriter) Text(text string) {
	// Set content type header
	hw.W.Header().Set("Content-Type", "text/plain")

	// Write text data
	_, err := hw.W.Write([]byte(text))
	if err != nil {
		hw.W.WriteHeader(http.StatusInternalServerError)
		http.Error(hw.W, "Failed to write text response", http.StatusInternalServerError)
		return
	}
}

// Error writes a Error message to the Http response stream
// in a more strucutured way
func (hw *HttpWriter) Error(err error) {
	// Set content type header
	hw.W.WriteHeader(http.StatusInternalServerError)
	hw.W.Header().Set("Content-Type", "text/plain")

	// Write error response
	http.Error(hw.W, err.Error(), hw.StatusCode)
}
