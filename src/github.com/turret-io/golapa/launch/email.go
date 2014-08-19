package launch

import (
	"appengine"
	"appengine/mail"
	"bytes"
	"text/template"
)

const sender = "<APPENGINE-EMAIL-SENDER>"
const recipient = "<RECIPIENT>"
const subject = "New signup"

func SendContactNotification(ctx appengine.Context, name string, email string) {
	tpl, err := template.ParseFiles("templates/email.tpl")
	if err != nil {
		ctx.Errorf("Error loading email template: %v", err)
	}
	kv := map[string]string{}
	kv["name"] = name
	kv["email"] = email

	var body bytes.Buffer
	tpl.Execute(&body, kv)
	msg := &mail.Message{
		Sender:  sender,
		To:      []string{recipient},
		Subject: subject,
		Body:    string(body.Bytes()),
	}
	if err := mail.Send(ctx, msg); err != nil {
		ctx.Errorf("Could not send email: %v", err)
	}

}
