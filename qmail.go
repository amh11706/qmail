package qmail

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"

	mailyak "github.com/domodwyer/mailyak/v3"
)

var smtpServer = os.Getenv("SMTP_SERVER")
var from = os.Getenv("SMTP_USER")
var password = os.Getenv("SMTP_PASSWORD")

func init() {
	fmt.Println("SMTP_SERVER:", smtpServer)
	fmt.Println("SMTP_USER:", from)
}

var auth = smtp.PlainAuth("", from, password, smtpServer)

func newMessage(to, subject, body string) *mailyak.MailYak {
	m, err := mailyak.NewWithTLS(smtpServer+":465", auth, nil)
	if err != nil {
		fmt.Println("Error creating email:", err)
		return nil
	}
	m.To(to)
	m.From(from)
	m.FromName("Captain Quacker")
	m.Subject(subject)
	m.Plain().Set(body)
	return m
}

func prettyPrintJSON(report string) string {
	var data interface{}
	err := json.Unmarshal([]byte(report), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return report
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return report
	}

	return string(prettyJSON)
}

func SendCspAlert(report string) error {
	m := newMessage("amh11706@gmail.com", "CSP Violation", fmt.Sprintf(
		`A CSP violation has occurred on superquacken.com. The report is as follows:
%s`, prettyPrintJSON(report)))
	return m.Send()
}

func SendCrashReport(err error, stack string) error {
	m := newMessage("amh11706@gmail.com", "Quacken Crash Report", fmt.Sprintf("Error: %v\n\nStack: %v", err, stack))
	return m.Send()
}

func SendReset(to, token string) error {
	m := newMessage(to, "Account Recovery", fmt.Sprintf(
		`You are receiving this email because someone requested a link to reset your Quacken password.
This link is valid for 30 minutes from the time it was sent.

https://www.superquacken.com/#/auth/reset/%s`, token))
	return m.Send()
}

func SendRestore(to, token string) error {
	m := newMessage(to, "Account Email Changed", fmt.Sprintf(
		`You are receiving this message because the email associated with your Quacken account has been changed.
You may use this link at any time to revert the change.

https://www.superquacken.com/#/auth/restore/%s`, token))
	return m.Send()
}

func SendDelete(to, token string) error {
	m := newMessage(to, "Account Deleted", fmt.Sprintf(
		`You are receiving this message because your Quacken account has been deleted.
You may use this link at any time to restore your account. Your email address has been deleted from our database, but we will keep a
backup of the rest of the account in case you change your mind.

https://www.superquacken.com/#/auth/restore/%s`, token))
	return m.Send()
}
