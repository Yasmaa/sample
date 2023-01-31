package auth

import (
	"api/config"
	"bytes"
	"fmt"
	// em "github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
)

func SendMail(subject string, email string, url string, temp string) bool {

	// email
	to := []string{
		`test@local.test`,
	}

	// auth := smtp.PlainAuth("", config.C.Smtp.USER, config.C.Smtp.PASSWORD, config.C.Smtp.HOST)
	// auth := smtp.CRAMMD5Auth(config.C.Smtp.USER, config.C.Smtp.PASSWORD)
	t, _ := template.ParseFiles("templates/" + temp)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+"Subject: %s \n%s\n\n", config.C.Smtp.FROM, email, subject, mimeHeaders)))

	t.Execute(&body, struct {
		Url string
	}{
		Url: url})

	// e := em.NewEmail()
	// e.From = `test@local.test`
	// e.To = []string{`test@local.test`}
	// e.Subject = `test mail`
	// e.Text = body.Bytes()

	// Sending email.
	err := smtp.SendMail(config.C.Smtp.HOST+":"+config.C.Smtp.PORT, nil , config.C.Smtp.FROM, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
