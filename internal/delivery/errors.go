package delivery

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Error struct {
	Statusint        int
	StatusText       string
	StatusTextandInt string
	MessageError     string
}

func (h *Handler) Errors(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	t, err := template.ParseFiles("templates/errors.html")
	if err != nil {
    log.Print(err)
		http.Error(w, strconv.Itoa(http.StatusInternalServerError)+" "+"Error parsing file", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	error1 := Error{status, http.StatusText(status), strconv.Itoa(status) + " " + http.StatusText(status), message}
  log.Print(error1)
	if err := t.Execute(w, error1); err != nil {
		log.Print(err)
		http.Error(w, strconv.Itoa(http.StatusInternalServerError)+" "+"Error executing file", http.StatusInternalServerError)
		return
	}
}
