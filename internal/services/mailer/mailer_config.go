package mailer

import (
	"bytes"

	"github.com/dtg-lucifer/everato/internal/utils"
)

// Enums
const (
	MailerSuccess = iota // Indicates that the email was sent successfully
	MailerFailure        // Indicates that there was an error sending the email
)

// Gomailer configuration options
type MailerOptions struct {
	Host        string // SMTP server host
	Port        uint16 // SMTP server port
	SenderEmail string // Email address of the sender
	AppPassword string // Password for the sender's email account
}

// This config is a must to send inside the SendMail method
type MailerParameters struct {
	To      string         // Recipient email address
	Subject string         // Email subject
	Body    *bytes.Buffer  // Email body content
	Options *MailerOptions // Configuration for gomail
}

// Mailer interface defines the methods for sending emails.
type Mailer interface {
	SendEmail(wr *utils.HttpWriter) (uint8, error)
}

type MailService struct {
	Config *MailerParameters
}

// Explicitly assert the implementation of Mailer interface
var _ Mailer = (*MailService)(nil)
