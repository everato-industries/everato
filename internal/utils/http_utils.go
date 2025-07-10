// Package utils provides utility functions for the Everato application.
// It includes helpers for HTTP handling, environment variables, and data conversion.
package utils

import (
	"encoding/json"
	"errors"
	"github.com/dtg-lucifer/everato/pkg"
	"net/http"
	"strings"
)

// Constants for common HTTP header values
// These are used throughout the application for consistent content type handling
const (
	HeaderContentTypeName = "Content-Type"             // Standard Content-Type header name
	HeaderContentTypeJson = "application/json"         // JSON content type value
	HeaderContentTypeText = "text/plain"               // Plain text content type value
	HeaderContentTypeHtml = "text/html; charset=utf-8" // HTML content type with UTF-8 charset
)

// CookieParams provides a structured way of passing parameters to the SetCookie method.
// It encapsulates all standard cookie attributes in a single structure for convenience.
type CookieParams struct {
	Name     string        // Name of the cookie
	Value    string        // Value of the cookie
	MaxAge   int           // Maximum age of the cookie in seconds
	Path     string        // Path for which the cookie is valid
	Domain   string        // Domain for which the cookie is valid
	Secure   bool          // Whether the cookie should be secure (only sent over HTTPS)
	HttpOnly bool          // Whether the cookie should be HTTP-only (not accessible via JavaScript)
	SameSite http.SameSite // SameSite attribute for the cookie (None, Lax, or Strict)
}

// M is a type alias for map[string]any, providing a concise way to represent
// key-value pairs used in JSON responses and template data.
type M map[string]any

// HttpWriter is a utility struct that wraps standard http.ResponseWriter and http.Request
// to provide a more convenient fluent API for writing HTTP responses.
// It includes methods for writing JSON, HTML, and error responses with chainable calls.
type HttpWriter struct {
	W          http.ResponseWriter // Underlying HTTP response writer
	R          *http.Request       // Associated HTTP request
	StatusCode int                 // HTTP status code to use for the response
	BufferSize uint                // Maximum size of request body in bytes (5MB default)
}

// NewHttpWriter creates and returns a new HttpWriter instance.
// It wraps the standard ResponseWriter and Request objects with additional functionality.
//
// Parameters:
//   - w: Standard HTTP response writer
//   - r: HTTP request object
//
// Returns:
//   - A configured HttpWriter with default status code and buffer size
func NewHttpWriter(w http.ResponseWriter, r *http.Request) *HttpWriter {
	return &HttpWriter{
		W:          w,               // http.ResponseWriter
		R:          r,               // *http.Request
		StatusCode: http.StatusOK,   // Default status code
		BufferSize: 5 * 1024 * 1024, // Default size of request body 5 MB (5242880 bytes)
	}
}

// Status sets the HTTP status code for the response.
// This method supports method chaining for fluent API usage.
//
// Parameters:
//   - code: HTTP status code (e.g., http.StatusOK, http.StatusBadRequest)
//
// Returns:
//   - The HttpWriter instance for method chaining
func (hw *HttpWriter) Status(code int) *HttpWriter {
	hw.StatusCode = code // Set the status code to the struct to use it on chained operations

	// We'll write the header just once in the Json/Text/Error methods
	return hw
}

// Json writes a JSON response to the HTTP response writer.
// It automatically includes the request ID in the response for traceability.
//
// This method:
// 1. Adds request ID to the response data
// 2. Marshals the data to JSON
// 3. Sets appropriate content type and status headers
// 4. Writes the JSON data to the response
//
// Parameters:
//   - data: Map of data to be serialized as JSON
func (hw *HttpWriter) Json(data M) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Append the request id with the original data for traceability
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

// Html renders an HTML template and writes the result to the HTTP response.
// It loads the template from the filesystem, sets the appropriate content type,
// and executes the template with the provided data.
//
// Parameters:
//   - view: Path to the HTML template file to render
//   - data: Data to pass to the template for rendering
func (hw *HttpWriter) Html(view string, data any) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Check if the view path exists and load the template
	tmp, err := pkg.GetTemplate(view)
	if err != nil {
		logger.StdoutLogger.Error("Error loading HTML template", "template", view, "err", err.Error())
		hw.Status(http.StatusInternalServerError).Text("Error loading template")
		return
	}

	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeHtml)

	// Write status code that was set with Status()
	hw.W.WriteHeader(hw.StatusCode)

	// Execute the template with the provided data
	err = tmp.Execute(hw.W, data)
	if err != nil {
		logger.StdoutLogger.Error("Error executing HTML template", "err", err.Error())
		logger.FileLogger.Error("Error executing HTML template", "err", err.Error())
		return
	}
}

// ParseBody parses the JSON request body into the provided target struct or map.
// It validates that the request has a body and the proper Content-Type header.
//
// Parameters:
//   - body: Pointer to a struct or map where the parsed JSON will be stored
//
// Returns:
//   - error: If the request has no body, invalid content type, or parsing fails
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

// Text writes a plain text response to the HTTP response writer.
// It automatically appends the request ID to the response for traceability.
//
// Parameters:
//   - text: The text content to write in the response
func (hw *HttpWriter) Text(text string) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeText)

	// Append the request id with the response for traceability
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

// Error writes an error response to the HTTP response writer.
// It sets an appropriate status code and includes the request ID for traceability.
//
// Parameters:
//   - err: The error to include in the response
//   - status_code: Optional status code to use (defaults to InternalServerError if not specified or invalid)
func (hw *HttpWriter) Error(err error, status_code ...int) {
	// Set content type header - must be set BEFORE WriteHeader
	hw.W.Header().Set(HeaderContentTypeName, HeaderContentTypeText)

	// If no status code is explicitly set, use InternalServerError
	if hw.StatusCode == http.StatusOK {
		if len(status_code) > 0 {
			if status_code[0] < 300 {
				hw.StatusCode = http.StatusInternalServerError
			} else {
				hw.StatusCode = status_code[0]
			}
		} else {
			hw.StatusCode = http.StatusInternalServerError
		}
	}

	// Write status code
	hw.W.WriteHeader(hw.StatusCode)

	// Append the request id with the error for traceability
	reqId := hw.W.Header().Get("X-Request-ID")
	res := strings.Join([]string{err.Error(), ("\nRequestId=" + reqId + ";")}, ";")

	// Write error message
	hw.W.Write([]byte(res))
}

// SetCookie sets a cookie in the HTTP response using the provided parameters.
// This method provides a convenient way to set cookies with all common attributes.
//
// Parameters:
//   - params: CookieParams struct containing all cookie attributes
//
// Example usage:
//
//	hw.SetCookie(utils.CookieParams{
//	    Name:     "session",
//	    Value:    sessionToken,
//	    MaxAge:   3600,               // 1 hour
//	    Path:     "/",
//	    Secure:   true,
//	    HttpOnly: true,
//	    SameSite: http.SameSiteLaxMode,
//	})
func (hw *HttpWriter) SetCookie(params CookieParams) {
	cookie := &http.Cookie{
		Name:     params.Name,
		Value:    params.Value,
		MaxAge:   params.MaxAge,
		Path:     params.Path,
		Domain:   params.Domain,
		Secure:   params.Secure,
		HttpOnly: params.HttpOnly,
		SameSite: params.SameSite,
	}

	http.SetCookie(hw.W, cookie)
}
