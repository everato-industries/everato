package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dtg-lucifer/everato/pkg"
)

// Constants for header values
const (
	HeaderContentTypeName = "Content-Type"
	HeaderContentTypeJson = "application/json"
	HeaderContentTypeText = "text/plain"
)

// Mapping of the default map of a string to any value to custom type
type M map[string]any

// HttpWriter is a utility struct to handle HTTP responses
// it helps write response to the http stream more easily
type HttpWriter struct {
	W          http.ResponseWriter
	R          *http.Request
	StatusCode int
	BufferSize uint
}

// Returns a new HttpWriter instance
func NewHttpWriter(w http.ResponseWriter, r *http.Request) *HttpWriter {
	return &HttpWriter{
		W:          w,               // http.ResponseWriter
		R:          r,               // *http.Request
		StatusCode: http.StatusOK,   // Default status code
		BufferSize: 5 * 1024 * 1024, // Default size of request body 5 MB
	}
}

// Status sets the HTTP status code for the response
func (hw *HttpWriter) Status(code int) *HttpWriter {
	hw.StatusCode = code // Set the status code to the struct
	// -					to use it on chained operations

	// We'll write the header just once in the Json/Text/Error methods
	return hw
}

// Json writes a JSON response to the HTTP response writer
func (hw *HttpWriter) Json(data M) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Append the request id with the original data
	if r := hw.W.Header().Get("X-Request-ID"); r == "" {
		data["request_id"] = "unknown"
	} else {
		data["request_id"] = r
	}

	// Convert data to JSON bytes
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Set headers before writing status
		hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeText)
		hw.W.WriteHeader(http.StatusInternalServerError)
		hw.W.Write([]byte("Failed to marshal JSON"))
		return
	}

	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeJson)

	// Write status code that was set with Status()
	hw.W.WriteHeader(hw.StatusCode)

	// Write JSON data
	_, err = hw.W.Write(jsonData)
	if err != nil {
		logger.StdoutLogger.Error("Error writing JSON to the response", "err", err.Error())
		logger.FileLogger.Error("Error writing JSON to the response", "err", err.Error())
		return
	}
}

// ParseBody method takes pointer to either a map or a struct
func (hw *HttpWriter) ParseBody(body any) error {
	// Check if the body is not provided
	if hw.R.Body == nil {
		return errors.New("the request doesn't have a body")
	}

	// Check if the request doesn't have a proper JSON body
	contentType := hw.R.Header.Get(HeaderContentTypeName)
	if contentType == "" || !strings.Contains(contentType, HeaderContentTypeJson) {
		return errors.New("the request should have a proper JSON body")
	}

	raw := hw.R.Body  // Getting the raw body
	defer raw.Close() // Close the body

	decoder := json.NewDecoder(raw) // Decoding the raw body
	err := decoder.Decode(body)     // Into the actual map / struct
	if err != nil {
		return errors.New("Failed to parse JSON body: " + err.Error())
	}

	return nil
}

// Text writes a TEXT response to the HTTP response writer
func (hw *HttpWriter) Text(text string) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeText)

	// Append the request id with the response
	reqId := hw.W.Header().Get("X-Request-ID")
	res := strings.Join([]string{text, ("RequestId=" + reqId)}, ";")

	// Write status code that was set with Status()
	hw.W.WriteHeader(hw.StatusCode)

	// Write text data
	_, err := hw.W.Write([]byte(res))
	if err != nil {
		logger.StdoutLogger.Error("Error writing TEXT to the response", "err", err.Error())
		return
	}
}

// Error writes a Error message to the Http response stream
// in a more strucutured way
func (hw *HttpWriter) Error(err error, status_code ...int) {
	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeText)

	// If no status code is explicitly set, use InternalServerError
	if hw.StatusCode == http.StatusOK {
		if status_code[0] < 300 {
			hw.StatusCode = http.StatusInternalServerError
		} else {
			hw.StatusCode = status_code[0]
		}
	}

	// Write status code
	hw.W.WriteHeader(hw.StatusCode)

	// Append the request id with the error
	reqId := hw.W.Header().Get("X-Request-ID")
	res := strings.Join([]string{err.Error(), ("RequestId=" + reqId)}, ";")

	// Write error message
	hw.W.Write([]byte(res))
}
