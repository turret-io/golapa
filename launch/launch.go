package launch

import (
	"appengine"
	"appengine/taskqueue"
	_ "appengine/urlfetch"
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"html/template"
	"net/http"
	"net/url"
)

const api_key = string("")
const api_secret = string("")

func init() {
	http.HandleFunc("/", MainHandler)
	http.HandleFunc("/worker", EmailSubmitter)
}

func handlePost(w http.ResponseWriter, c appengine.Context, tpl *template.Template, data url.Values) {

	email := data.Get("email")
	name := data.Get("name")
	if name != "" && email != "" {
		t := taskqueue.NewPOSTTask("/worker", map[string][]string{"name": {name}, "email": {email}})
		if _, err := taskqueue.Add(c, t, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// OK
	}
	tpl.Lookup("post.tpl").Execute(w, data)

}

func handleGet(w http.ResponseWriter, name string, tpl *template.Template) {
	switch name {
	case "c/terms":
		tpl.Lookup("terms.tpl").Execute(w, nil)
	case "c/privacy":
		tpl.Lookup("privacy.tpl").Execute(w, nil)
	default:
		tpl.Lookup("main.tpl").Execute(w, nil)

	}
}

func signString(s string, k string) string {
	h := hmac.New(sha512.New, []byte(k))
	return string(h.Sum([]byte(s)))
}

func EmailSubmitter(w http.ResponseWriter, r *http.Request) {
	// Email new contact
	c := appengine.NewContext(r)
	SendContactNotification(c, r.FormValue("name"), r.FormValue("email"))

}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	tpl, err := template.ParseFiles("templates/main.tpl", "templates/post.tpl", "templates/terms.tpl", "templates/privacy.tpl", "templates/header.tpl", "templates/footer.tpl")
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		ctx := appengine.NewContext(r)
		ctx.Infof(r.URL.Path[len("/"):])
		handleGet(w, r.URL.Path[len("/"):], tpl)

	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		handlePost(w, c, tpl, r.PostForm)
	default:
		http.Error(w, "Method not implemented", http.StatusInternalServerError)
	}
}
