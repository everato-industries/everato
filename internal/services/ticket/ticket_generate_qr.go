package ticket

import (
	"bytes"
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRImage encodes the given content string into a QR code PNG image
// and returns the base64-encoded PNG bytes (suitable for embedding in HTML/JSON).
//
// Parameters:
//   - content: The string to encode (typically the signed QR payload JSON)
//   - size:    Pixel size of the QR code image (256 is a good default)
//
// Returns:
//   - base64-encoded PNG string, or error
func GenerateQRImage(content string, size int) (string, error) {
	if size <= 0 {
		size = 256
	}

	// Encode QR code into a byte buffer
	var buf bytes.Buffer
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return "", fmt.Errorf("qrcode.Encode: %w", err)
	}

	buf.Write(png)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DataURLForQR returns a complete data: URL for embedding the QR image in HTML/email.
//
// Example output: "data:image/png;base64,iVBORw0KGgo..."
func DataURLForQR(content string, size int) (string, error) {
	b64, err := GenerateQRImage(content, size)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + b64, nil
}
