// +build appengine

package golapa

import (
	"appengine"
	"appengine/taskqueue"
	"appengine/mail"
	_ "fmt"
	"net/url"
	"net/http"
	"text/template"
)

const sender = "<APPENGINE-EMAIL-SENDER>"
const to = "<RECIPIENT>"
const subject = "New signup"

func init() {
	http.Handle("/", newAppEngineHandler())
	http.Handle("/worker", newAppEngineEmailHandler())
}

func newAppEngineHandler() (*AppEngineHandler) {
	aeh := &AppEngineHandler{}
	aeh.BaseHandler.TemplatePath = "../templates/"
	return aeh
}

func newAppEngineEmailHandler() (*AppEngineEmailer) {
	aee := &AppEngineEmailer{}
	aee.TemplatePath = "../templates/"
	return aee
}


type AppEngineEmailer struct {
	StandardEmailer
	ctx appengine.Context
}

func (aee *AppEngineEmailer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle worker requests
	ctx := appengine.NewContext(req)
	aee.ctx = ctx
	msg, err := aee.CreateMessage(sender, subject, to, req.FormValue("name"), req.FormValue("email"))
	if err != nil {
		ctx.Errorf("Could not create email: %v", err)
	}
	aee.Send(msg)
}

func (aee *AppEngineEmailer) Send(msg *Message) {
	message := &mail.Message{
		Sender: msg.Sender,
		To: []string{msg.To},
		Subject: msg.Subject,
		Body: msg.Body,
	}
	if err := mail.Send(aee.ctx, message); err != nil {
		aee.ctx.Errorf("Could not send email: %v", err)
	}
}

type AppEngineRequestHandler struct {
	StandardRequestHandler
}

func (aerh *AppEngineRequestHandler) Get(w http.ResponseWriter, r *http.Request, name string, tpl *template.Template) {
	ctx := appengine.NewContext(r)
	ctx.Infof(r.URL.Path[len("/"):])
	aerh.StandardRequestHandler.Get(w, r, name, tpl)
}

func (aerh *AppEngineRequestHandler) Post(w http.ResponseWriter, r *http.Request, tpl *template.Template, data url.Values, template_path string) {
	ctx := appengine.NewContext(r)
	ctx.Infof(r.URL.Path[len("/"):])

	email := data.Get("email")
	name := data.Get("name")
	if name != "" && email != "" {
		t := taskqueue.NewPOSTTask("/worker", map[string][]string{"name": {name}, "email": {email}})
		if _, err := taskqueue.Add(ctx, t, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// OK
	}
	tpl.Lookup("post.tpl").Execute(w, data)
}

type AppEngineHandler struct {
	BaseHandler
}

func (aeh *AppEngineHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	aeh.BaseHandler.Serve(w, req, new(AppEngineRequestHandler))
}
