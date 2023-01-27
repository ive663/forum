package delivery

import (
	"net/http"
	"text/template"

	"github.com/ive663/forum/internal/service"
)

type Handler struct {
	templates *template.Template
	services  *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", h.authenticateUser(h.index))
	mux.HandleFunc("/signin", h.signin)
	mux.HandleFunc("/signup", h.signup)
	mux.HandleFunc("/createpost", h.authenticateUser(h.createpost))
	mux.HandleFunc("/logout", h.logout)
	mux.HandleFunc("/post", h.authenticateUser(h.post))
	mux.HandleFunc("/likepost", h.authenticateUser(h.likePost))
	mux.HandleFunc("/likepostindex", h.authenticateUser(h.likePostIndex))
	mux.HandleFunc("/likecomment", h.authenticateUser(h.likeComment))
	mux.HandleFunc("/dislikecomment", h.authenticateUser(h.dislikeComment))
	mux.HandleFunc("/dislikepost", h.authenticateUser(h.dislikePost))
	mux.HandleFunc("/dislikepostindex", h.authenticateUser(h.dislikePostIndex))
	return mux
}
