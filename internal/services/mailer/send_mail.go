// Package mailer provides email notification services for the Everato platform.
// It handles email composition, template rendering, and delivery via SMTP.
package mailer

import (
	"github.com/dtg-lucifer/everato/internal/utils"
	"gopkg.in/gomail.v2"
)

// NewMailService creates and initializes a new email service instance.
// It validates the configuration parameters before creating the service to
// prevent runtime errors during email sending operations.
//
// Parameters:
//   - config: Email configuration including recipient, subject, body, and SMTP settings
//
// Returns:
//   - A configured MailService instance ready to send emails, or nil if config is invalid
func NewMailService(config *MailerParameters) *MailService {
	if config == nil {
		return nil // Return nil if the config is nil to avoid panic
	}
	return &MailService{
		Config: config,
	}
}

// SendEmail sends an email using the configured parameters and returns the status of the operation.
//
// This method performs the following operations:
// 1. Creates an SMTP dialer with the configured server details
// 2. Composes an email message with proper headers and content
// 3. Establishes a connection to the SMTP server
// 4. Sends the email message
// 5. Returns the appropriate status code and error information
//
// The email is sent as HTML content with UTF-8 encoding. All email headers
// (From, To, Subject) are properly set based on the service configuration.
//
// Parameters:
//   - wr: Custom HTTP writer for request context (not used in current implementation
//     but included for consistency with other service interfaces)
//
// Returns:
//   - uint8: Status code (MailerSuccess or MailerFailure)
//   - error: Detailed error information if sending fails, nil on success
func (m *MailService) SendEmail(wr *utils.HttpWriter) (uint8, error) {
	// Create the SMTP dialer with server configuration
	// This establishes the connection parameters for the mail server
	dialer := gomail.NewDialer(
		m.Config.Options.Host,
		int(m.Config.Options.Port),
		m.Config.Options.SenderEmail,
		m.Config.Options.AppPassword,
	)

	// Create a new email message with proper structure
	message := gomail.NewMessage()

	// Set required email headers
	message.SetHeader("Subject", m.Config.Subject)
	message.SetHeader("From", m.Config.Options.SenderEmail)
	message.SetHeader("To", m.Config.To)

	// Set the message body as HTML content
	// The body is expected to be a rendered HTML template
	message.SetBody("text/html", m.Config.Body.String())

	// Attempt to connect to the SMTP server and send the message
	// This handles connection, authentication, and message delivery
	err := dialer.DialAndSend(message)
	if err != nil {
		// Return failure status with the error details
		return MailerFailure, err
	}

	// Return success status with no error
	return MailerSuccess, nil
}
