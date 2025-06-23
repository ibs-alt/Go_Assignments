package todo

import (
	"html/template"
	"net/http"
)

func ListHandler(reads chan<- readReq) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/list.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		req := readReq{resp: make(chan []Todo)}
		reads <- req
		tmpl.Execute(w, <-req.resp)
	}
}
