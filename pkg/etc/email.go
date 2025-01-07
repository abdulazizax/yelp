package etc

import (
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

type Otp struct {
	Code string
}

// generateEmailBody generates the HTML email body for a student
func GenerateOtpEmailBody(otp string) (string, error) {
	templateString := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f7fa;
        }
        .email-container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .email-header {
            text-align: center;
            margin-bottom: 20px;
        }
        .email-header h2 {
            color: #333333;
        }
        .email-body {
            font-size: 16px;
            color: #555555;
            line-height: 1.5;
        }
        .otp-code {
            font-size: 20px;
            font-weight: bold;
            color: #007BFF;
            padding: 10px;
            background-color: #f0f8ff;
            border-radius: 4px;
            margin-top: 15px;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 14px;
            color: #888888;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <h2>Welcome to Yelp!</h2>
        </div>
        <div class="email-body">
            <p>Hi there,</p>
            <p>Thank you for signing up to Yelp! To complete your registration and verify your account, please use the OTP (One-Time Password) below:</p>
            <div class="otp-code">{{.Code}}</div>
            <p>This code is valid for 10 minutes. If you did not request this, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Yelp. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	otpData := Otp{otp}

	var builder strings.Builder
	err = tmpl.Execute(&builder, otpData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return builder.String(), nil
}

// sendEmail sends an email using SMTP
func SendEmail(smtpHost, smtpPort, from, password, to, body string) error {
	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("Subject: Otp Code for Yelp Account Verification\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"\r\n%s", from, to, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
