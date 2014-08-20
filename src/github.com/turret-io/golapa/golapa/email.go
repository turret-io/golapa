package golapa

import (
	"fmt"
	"bytes"
	"text/template"
	"os"
	"net/smtp"
	"net"
	"time"
	"log"
	"encoding/base64"
)

type Message struct {
	Sender string
	To	[]string
	Subject string
	Body	string
}

func (msg *Message) toMessage() ([]byte) {
	var b bytes.Buffer
	t := time.Now()
	curDateTime := t.UTC().Format(time.RFC1123Z)

	b.Write([]byte(fmt.Sprintf("Sender: %s \r\n", msg.Sender)))
	b.Write([]byte(fmt.Sprintf("Reply-To: %s\r\n", msg.Sender)))
	b.Write([]byte(fmt.Sprintf("Subject: =?utf-8?B?%s?=\r\n", base64.URLEncoding.EncodeToString([]byte(msg.Subject)))))
	b.Write([]byte(fmt.Sprintf("Date: %s\r\n", curDateTime)))
	b.Write([]byte(fmt.Sprintf("From: %s %s\r\n", msg.Sender)))
	b.Write([]byte(fmt.Sprintf("To: %s\r\n", msg.To)))

	return b.Bytes()
}

type Emailer interface {
	Send(msg *Message)
}

type StandardEmailer struct {
	TemplatePath	string
}

func (se *StandardEmailer) Send(msg *Message) {
	host := os.Getenv("EMAIL_HOST")
	password := os.Getenv("EMAIL_PASSWORD")
	from := os.Getenv("EMAIL_FROM")
	hostPart, _, err := net.SplitHostPort(host)
	if err != nil {
		log.Print(err.Error())
	}
	auth := smtp.PlainAuth("", from, password, hostPart)

	err = smtp.SendMail(host, auth, from, msg.To, msg.toMessage())
	if err != nil {
		log.Print(err.Error())
	}
}

func (se *StandardEmailer) CreateMessage(sender string, subject string, name string, email string) (*Message, error) {
	tpl, err := template.ParseFiles(fmt.Sprintf("%s/email.tpl", se.TemplatePath))
	if err != nil {
		return nil, err
	}
	kv := map[string]string{}
	kv["name"] = name
	kv["email"] = email

	var body bytes.Buffer
	tpl.Execute(&body, kv)

	return &Message{
		Sender: sender,
		To: []string{email},
		Subject: subject,
		Body: string(body.Bytes()),
	}, nil
}
