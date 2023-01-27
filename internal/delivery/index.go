package delivery

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/ive663/forum/internal/module"
	"github.com/ive663/forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		user_id, ok := r.Context().Value(keyUserID).(int)
		if !ok {
			return
		}
		user_authorization := true
		if user_id == 0 {
			user_authorization = false
		}
		t, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			log.Print("err:delivery:index: ParseFiles", err)
			h.Errors(w, http.StatusInternalServerError, "Error parsing file")
			return
		}
		var posts module.PostList
		if len(r.URL.Query()) == 0 {
			posts, err = h.services.GetNewPosts()
			if err != nil {
				log.Print("err:delivery:index: GetNewPosts")
				h.Errors(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			posts, err = h.services.GetAllPostBy(user_id, r.URL.Query())
			if err != nil {
				log.Print("err:delivery:index: GetAllPostBy")
				if errors.Is(err, service.ErrInvalidQueryRequest) {
					h.Errors(w, http.StatusNotFound, "Invalid query request")
					return
				}
				h.Errors(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		u := module.User{
			Posts:         posts.PrepToView(),
			Authorization: user_authorization,
		}

		if err = t.Execute(w, u); err != nil {
			log.Print(err)
			log.Print("err:delivery:index: Execute")
			h.Errors(w, http.StatusInternalServerError, "Error executing file")
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
