// +build !appengine

package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/turret-io/golapa/golapa"
)

var template_path string

func init() {
	flag.StringVar(&template_path, "templates", "templates", "Path to templates")
	flag.Parse()
	http.Handle("/", newBaseHandler())
}

func newBaseHandler() (*golapa.BaseHandler) {
	bh := &golapa.BaseHandler{template_path}
	return bh
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
