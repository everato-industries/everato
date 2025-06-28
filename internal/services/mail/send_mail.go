package mailer

import "github.com/dtg-lucifer/everato/internal/utils"

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
	// Simulate sending email
	return MailerSuccess, nil
}
