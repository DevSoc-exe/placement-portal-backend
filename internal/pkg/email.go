package pkg

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(to, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASSWORD")
	if senderEmail == "" || password == "" {
		return fmt.Errorf("email or password not found in environment variables")
	}

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	subject := "TPC Registration Email Verification"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to the SMTP server: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to upgrade connection to TLS: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate to the SMTP server: %w", err)
	}

	if err = client.Mail(senderEmail); err != nil {
		return fmt.Errorf("failed to set the sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set the recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email body: %w", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	if err = client.Quit(); err != nil {
		return fmt.Errorf("failed to quit the SMTP client session: %w", err)
	}

	fmt.Println("Email sent successfully.")
	return nil
}


func CreateMailMessageWithVerificationToken(token string, userID string) string {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:3000"
	}

	return "Please verify your email by clicking on this Verification link: " + domain + "/user/verify/" + token + "?uid=" + userID
}
