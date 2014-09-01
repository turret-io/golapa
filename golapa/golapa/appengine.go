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
const template_path = "../../templates/"
const api_key = "<TURRET-IO-API-KEY>"
const api_secret = "<TURRET-IO-API-SECRET>"

func init() {
	http.Handle("/", newAppEngineHandler())
	http.Handle("/worker", newAppEngineEmailHandler())
}

func newAppEngineHandler() (*AppEngineHandler) {
	aeh := &AppEngineHandler{}
	aeh.BaseHandler.TemplatePath = template_path 
	return aeh
}

func newAppEngineEmailHandler() (*AppEngineEmailer) {
	aee := &AppEngineEmailer{}
	aee.TemplatePath = template_path 
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

func newAppEngineTurretHandler() (*AppEngineTurret) {
    aet := &AppEngineTurret{}
    return aet
}

type AppEngineHandler struct {
	BaseHandler
}

func (aeh *AppEngineHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	aeh.BaseHandler.Serve(w, req, new(AppEngineRequestHandler))
}

type AppEngineTurret struct {
    ctx appengine.Context
}

/*
func (aet *AppEngineTurret) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // handle worker requests
    ctx := appengine.NewContext(req)
    aet.ctx = ctx
    attr_map := make(map[string]string)
    prop_map := make(map[string]string)
    attr_map["contact_name"] = req.FormValue("name") 
    attr_map["signedup"] = "2"
    turret := turretIO.NewAppEngineTurretIO(api_key, api_secret, ctx)
    inst := turretIO.NewUser(turret)
    resp, err := inst.Set(req.FormValue("email"), attr_map, prop_map)
    if err != nil {
        aet.ctx.Errorf(err.Error())
    }
}
*/
