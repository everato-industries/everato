// Package mailer provides email notification services for the Everato platform.
// It handles email composition, template rendering, and delivery via SMTP.
package mailer

import (
	"bytes"

	"github.com/dtg-lucifer/everato/internal/utils"
)

// Status constants for email delivery results
// These provide clear return values from email sending operations
const (
	MailerSuccess = iota // Indicates that the email was sent successfully
	MailerFailure        // Indicates that there was an error sending the email
)

// MailerOptions defines the SMTP server configuration for email delivery.
// It contains all necessary credentials and connection details for establishing
// a secure connection to the SMTP server.
type MailerOptions struct {
	Host        string // SMTP server hostname (e.g., smtp.gmail.com)
	Port        uint16 // SMTP server port (typically 587 for TLS or 465 for SSL)
	SenderEmail string // Email address used as the sender (From header)
	AppPassword string // Password or app-specific password for SMTP authentication
}

// MailerParameters encapsulates all parameters required to send an email.
// This structure must be provided to the SendEmail method with all fields populated.
// It combines both message content and delivery configuration in one structure.
type MailerParameters struct {
	To      string         // Recipient email address
	Subject string         // Email subject line
	Body    *bytes.Buffer  // Email body content as rendered HTML template
	Options *MailerOptions // SMTP server configuration and credentials
}

// Mailer interface defines the contract for email delivery services.
// Any email delivery implementation in the application must satisfy this interface,
// allowing for different email providers or mock implementations for testing.
type Mailer interface {
	// SendEmail sends an email using the configured parameters
	// Returns a status code (MailerSuccess or MailerFailure) and an error if applicable
	SendEmail(wr *utils.HttpWriter) (uint8, error)
}

// MailService implements the Mailer interface using the gomail package.
// It provides email delivery capabilities through standard SMTP protocols
// and supports HTML email templates with proper encoding.
type MailService struct {
	Config *MailerParameters // Email delivery configuration and content
}

// Explicitly assert the implementation of Mailer interface
var _ Mailer = (*MailService)(nil)
