package mailer

import (
	"github.com/dtg-lucifer/everato/internal/utils"
	"gopkg.in/gomail.v2"
)

// Instantiate the mail service
func NewMailService(config *MailerParameters) *MailService {
	if config == nil {
		return nil // Return nil if the config is nil to avoid panic
	}
	return &MailService{
		Config: config,
	}
}

// SendEmail sends an email using the provided configuration and returns the status of the operation.
//
// Parameters:
//   - params - This is a pointer to a MailerParameters struct
//   - wr - this is a pointer to a custom http response writer
//
// Returns:
//   - string - This will be either 0 or 1 (enumerated by `iota`)
//   - err - In case of failure it will hold the error information
func (m *MailService) SendEmail(wr *utils.HttpWriter) (uint8, error) {
	// Create the dialer
	dialer := gomail.NewDialer(
		m.Config.Options.Host,
		int(m.Config.Options.Port),
		m.Config.Options.SenderEmail,
		m.Config.Options.AppPassword,
	)

	// Instantiate the message struct
	message := gomail.NewMessage()

	message.SetHeader("Subject", m.Config.Subject)
	message.SetHeader("From", m.Config.Options.SenderEmail)
	message.SetHeader("To", m.Config.To)

	message.SetBody("text/html", m.Config.Body.String())

	err := dialer.DialAndSend(message) // Send the email using the dialer
	if err != nil {
		return MailerFailure, err
	}

	return MailerSuccess, nil
}
