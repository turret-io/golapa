// +build !appengine

package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/turret-io/golapa/golapa"
	"os"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"text/template"
	"net/url"

)

var template_path string

type CustomRequestHandler struct {
	golapa.StandardRequestHandler
}

func (crh *CustomRequestHandler) Post(w http.ResponseWriter,  r *http.Request, tpl *template.Template, data url.Values, template_path string) {
	sender := os.Getenv("EMAIL_SENDER")
	subject := os.Getenv("EMAIL_SUBJECT")
	to := os.Getenv("EMAIL_TO")
	via := os.Getenv("SEND_VIA")

	email := data.Get("email")
	name := data.Get("name")


	if name != "" && email != "" {
		log.Print(name)
		log.Print(email)

		if via == "email" {

			mailer := &golapa.StandardEmailer{template_path}
			msg, err := mailer.CreateMessage(sender, subject, to, r.FormValue("name"), r.FormValue("email"))
			if err != nil {
				log.Print("Could not create email: %v", err)
			}

			// JSON encode
			json, err := json.Marshal(msg)
			if err != nil {
				log.Print(err)
			}

			if err == nil {
				PushToRedis("signup_worker", json)
			}

		}

		if via == "turret_io" {
			json, err := json.Marshal(map[string]string{"name":r.FormValue("name"), "email":r.FormValue("email")})
			if err != nil {
				log.Print(err)
			}

			if err == nil {
				PushToRedis("signup_worker", json)
			}
		}
	}
	tpl.Lookup("post.tpl").Execute(w, data)
}

func init() {
	flag.StringVar(&template_path, "templates", "templates", "Path to templates")
	flag.Parse()
	http.Handle("/", newCustomHandler())
}

func newCustomHandler() (*CustomHandler) {
	ch := &CustomHandler{}
	ch.BaseHandler.TemplatePath = template_path
	return ch
}

type CustomHandler struct {
	golapa.BaseHandler
}

func (aeh *CustomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	aeh.BaseHandler.Serve(w, req, new(CustomRequestHandler))
}

func PushToRedis(list string, json []byte) (error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}

	defer conn.Close()

	conn.Do("LPUSH", list, json)

	return nil
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
