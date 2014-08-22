package golapa

import (
	"fmt"
	"bytes"
	"text/template"
	"os"
	"net/smtp"
	"crypto/tls"
	_ "crypto/sha512"
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

func (msg *Message) ToMessage() ([]byte) {
	var b bytes.Buffer
	t := time.Now()
	curDateTime := t.UTC().Format(time.RFC1123Z)

	b.Write([]byte(fmt.Sprintf("Sender: %s\r\n", msg.Sender)))
	b.Write([]byte(fmt.Sprintf("Reply-To: %s\r\n", msg.Sender)))
	b.Write([]byte(fmt.Sprintf("Subject: =?utf-8?B?%s?=\r\n", base64.URLEncoding.EncodeToString([]byte(msg.Subject)))))
	b.Write([]byte(fmt.Sprintf("Date: %s\r\n", curDateTime)))
	b.Write([]byte(fmt.Sprintf("From: %s\r\n", msg.Sender)))
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
	ssl := os.Getenv("EMAIL_SSL")

	hostPart, _, err := net.SplitHostPort(host)
	if err != nil {
		log.Print(err.Error())
	}
	auth := smtp.PlainAuth("", from, password, hostPart)

	if ssl == "" {
		err = smtp.SendMail(host, auth, from, msg.To, msg.ToMessage())
		if err != nil {
			log.Print(err.Error())
		}
	}

	if ssl == "true" {
		err = se.SendMailSSL(host, auth, from, msg.To, msg.ToMessage())
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func (se *StandardEmailer) CreateMessage(sender string, subject string, to string, name string, email string) (*Message, error) {
	tpl, err := template.ParseFiles(fmt.Sprintf("%s/email.tpl", se.TemplatePath))
	if err != nil {
		return nil, err
	}
	kv := map[string]string{}
	kv["name"] = name
	kv["email"] = email

	var body bytes.Buffer
	tpl.Execute(&body, kv)
	log.Print(string(body.Bytes()))
	return &Message{
		Sender: sender,
		To: to,
		Subject: subject,
		Body: string(body.Bytes()),
	}, nil
}

func (se *StandardEmailer) SendMailSSL(host string, auth smtp.Auth, from string, recipients []string, msg []byte) (error) {
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = c.Auth(auth)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = c.Mail(from)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	for _, addr := range recipients {
		err = c.Rcpt(addr)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	d, err := c.Data()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	_, err = d.Write(msg)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = d.Close()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	c.Quit()
	return nil
}
