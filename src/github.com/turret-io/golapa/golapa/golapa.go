package golapa

import (
	 "fmt"
	"net/http"
	"text/template"
	"net/url"
	"log"
	"os"
	"time"
)

type RequestHandler interface {
	Get(w http.ResponseWriter, r *http.Request, name string, tpl *template.Template)
	Post(w http.ResponseWriter, r *http.Request, tpl *template.Template, data url.Values, template_path string)
}

type StandardRequestHandler struct {}

func (srh *StandardRequestHandler) Get(w http.ResponseWriter, r *http.Request, name string, tpl *template.Template) {
	switch name {
	case "c/terms":
		tpl.Lookup("terms.tpl").Execute(w, nil)
	case "c/privacy":
		tpl.Lookup("privacy.tpl").Execute(w, nil)
	default:
		tpl.Lookup("main.tpl").Execute(w, nil)

	}
}

func (srh *StandardRequestHandler) Post(w http.ResponseWriter,  r *http.Request, tpl *template.Template, data url.Values, template_path string) {
	sender := os.Getenv("EMAIL_SENDER")
	subject := os.Getenv("EMAIL_SUBJECT")

	email := data.Get("email")
	name := data.Get("name")



	if name != "" && email != "" {
		log.Print(name)
		log.Print(email)


		// Create goroutine to send email via
		go func(){
			mailer := &StandardEmailer{template_path}
			msg, err := mailer.CreateMessage(sender, subject, r.FormValue("name"), r.FormValue("email"))
			if err != nil {
				log.Print("Could not create email: %v", err)
			}
			mailer.Send(msg)
			time.Sleep(5000 * time.Millisecond)
			
		}()

		// OK
	}
	tpl.Lookup("post.tpl").Execute(w, data)



type BaseHandler struct {
	TemplatePath	string
}

func (bh *BaseHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	bh.serve(w, req, new(StandardRequestHandler))
}

func (bh *BaseHandler) serve(w http.ResponseWriter, r *http.Request, rh RequestHandler) {
	tpl, err := template.ParseFiles(fmt.Sprintf("%s/main.tpl", bh.TemplatePath), fmt.Sprintf("%s/post.tpl", bh.TemplatePath), fmt.Sprintf("%s/terms.tpl", bh.TemplatePath), fmt.Sprintf("%s/privacy.tpl", bh.TemplatePath), fmt.Sprintf("%s/header.tpl", bh.TemplatePath), fmt.Sprintf("%s/footer.tpl", bh.TemplatePath))
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	switch r.Method {
	case "GET":
		rh.Get(w, r, r.URL.Path[len("/"):], tpl)

	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		rh.Post(w, r, tpl, r.PostForm, bh.TemplatePath)
	default:
		http.Error(w, "Method not implemented", http.StatusInternalServerError)
	}
}
