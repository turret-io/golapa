package launch

import (
	"appengine"
	"appengine/mail"
	"fmt"
)

const sender = "<APPENGINE-EMAIL-SENDER>"
const recipient = "<RECIPIENT>"

func SendContactNotification(ctx appengine.Context, name string, email string) {
	msg := &mail.Message{
		Sender:  sender,
		To:      []string{recipient},
		Subject: "New Signup",
		Body:    fmt.Sprintf("New signup:\n\n%s\n%s", name, email),
	}
	if err := mail.Send(ctx, msg); err != nil {
		ctx.Errorf("Could not send email: %v", err)
	}

}
