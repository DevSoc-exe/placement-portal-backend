package pkg

import (
	"fmt"
	"net/smtp"
	"os"
)

type Email struct {
	body    string
	to      string
	subject string
	mime    string
}

type MultipleEmail struct {
	body    string
	to      []string
	subject string
	mime    string
}

type CompanyCrux struct {
	Name     string
	Deadline string
	ID       string
}

func (e *Email) SendEmail() error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASSWORD")
	if senderEmail == "" || password == "" {
		return fmt.Errorf("email or password not found in environment variables")
	}

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	msg := []byte(fmt.Sprintf("%sTo: %s\r\n%s\r\n%s", e.subject, e.to, e.mime, e.body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{e.to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully.")
	return nil
}

func (e *MultipleEmail) SendEmail() error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	senderEmail := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASSWORD")
	if senderEmail == "" || password == "" {
		return fmt.Errorf("email or password not found in environment variables")
	}

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	msg := []byte(fmt.Sprintf("%sTo: %s\r\n%s\r\n%s", e.subject, e.to, e.mime, e.body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, e.to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully.")
	return nil
}

func CreateMailMessageWithVerificationToken(token string, userID string, userEmail string) *Email {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:3000"
	}

	return &Email{
		body:    "Please verify your email by clicking on this Verification link: " + domain + "/user/verify/" + token + "?uid=" + userID,
		to:      userEmail,
		subject: "Subject: TPC Registration Email Verification\r\n",
	}
}

func CreateOTPEmail(otp int, name string, email string) *Email {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:3000"
	}

	mail := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your OTP Code</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding: 10px 0;
            background-color: #4CAF50;
            color: #ffffff;
            border-radius: 8px 8px 0 0;
        }
        .content {
            text-align: center;
            padding: 20px;
        }
        .otp {
            font-size: 24px;
            font-weight: bold;
            color: #333333;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #888888;
            padding-top: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Your OTP Code</h2>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>Please use the following One-Time Password (OTP) to complete your login:</p>
            <div class="otp">%d</div>
            <p>This OTP is valid for 10 minutes. Please do not share it with anyone.</p>
        </div>
        <div class="footer">
            <p>If you did not request this, please contact our support team immediately.</p>
        </div>
    </div>
</body>
</html>
`, name, otp)
	return &Email{
		subject: "Subject: OTP for TPC Portal Login\r\n",
		body:    mail,
		to:      email,
		mime:    "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n",
	}
}

func CreateDriveUpdateNotificationEmail(email []string, companyInfo CompanyCrux) *MultipleEmail {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:3000"
	}

	mail := fmt.Sprintf(`<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Update Notification</title>
    <style>
				a {
				color:white;
				}
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding: 10px 0;
            background-color: #007BFF;
            color: #ffffff;
            border-radius: 8px 8px 0 0;
        }
        .content {
            text-align: center;
            padding: 20px;
        }
        .update-info {
            font-size: 16px;
            color: #333333;
            margin: 20px 0;
        }
        .cta-button {
            display: inline-block;
            padding: 10px 20px;
            font-size: 16px;
            color: #ffffff;
            background-color: #007BFF;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .cta-button:hover {
            background-color: #0056b3;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #888888;
            padding-top: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>%s drive is live!</h2>
        </div>
        <div class="content">
            <p>%s is conducting a placement drive</p>
            <p>Interested students to appply by %s</p>
            <a href="http://localhost:3000/user/drive/%s" class="cta-button">Check Out!</a>
        </div>
        <div class="footer">
            <p>Please read about the company in detail and Job Description carefully. Only those students who are genuinely interested in this opportunity and are willing to join should fill. Will not entertain any last minute drop out cases.</p>
        </div>
    </div>
</body>
</html>
`, companyInfo.Name, companyInfo.Name, companyInfo.Deadline, companyInfo.ID)

	return &MultipleEmail{
		subject: "Subject: New Placement Drive update!\r\n",
		body:    mail,
		to:      email,
		mime:    "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n",
	}
}
