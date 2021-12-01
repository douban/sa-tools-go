package notify

import "net/smtp"

func sendEmail(subject, content, from string, target []string) error {

	// Sender data.
	password := "<Email Password>"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, target, []byte(content))

}
