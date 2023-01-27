package delivery

import (
	"errors"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ive663/forum/internal/module"

	"github.com/ive663/forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/signup.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error parsing file")
			return
		}
		if err = t.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error executing")
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			log.Println("delivery:error: invalid email...")
			h.Errors(w, http.StatusBadRequest, "Please... use jon@smith.com format")
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			log.Println("delivery:error: invalid username... use latin letters")
			h.Errors(w, http.StatusBadRequest, "Please... use latin letters")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			log.Println("delivery:error: invalid password")
			h.Errors(w, http.StatusBadRequest, "Please use stronger password")
		}
		user := &module.User{
			Login:    username[0],
			Password: password[0],
			Email:    email[0],
		}
		newUser, err := h.services.Auth.CreateNewUser(user)
		if err != nil {
			if errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrInvalidPassword) || errors.Is(err, service.ErrInvalidUserName) {
				w.WriteHeader(http.StatusBadRequest)
				h.Errors(w, http.StatusBadRequest, err.Error())
				log.Println(":delivery:error: invalid email or password or username")
				return
			}
			h.Errors(w, http.StatusUnprocessableEntity, err.Error())
			log.Println("delivery:error: can't create user")
			return
		}
		tkn, err := h.services.Auth.GenerateSessionToken(newUser.Login, password[0])
		password[0] = ""
		user.Password = ""
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			log.Println("delivery:error: can't generate session token")
			return
		}
		if tkn == "" {
			h.Errors(w, http.StatusInternalServerError, "Error in token")
			log.Println("delivery:error: token is empty")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   tkn,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(12 * time.Hour),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
	}
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/signin.html")
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			t.Execute(w, nil)
			return
		}
		if len(cookie.Value) != 0 {
			log.Println("error: cookie is not empty")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		username, ok := r.Form["username"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please write latin password or login.")
			log.Println("error: invalid username")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please write latin password or login")
			log.Println("error: invalid password")
			return
		}
		token, err := h.services.Auth.GenerateSessionToken(username[0], password[0])
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				log.Println("error: user not found. can't generate token")
				h.Errors(w, http.StatusUnauthorized, err.Error())
				return
			}
			log.Println("error: user not authorized")
			h.Errors(w, http.StatusUnauthorized, err.Error())
			return
		}
		if token == "" {
			log.Println("log:session: token is empty, new session not created")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   token,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(12 * time.Hour),
		})
		password[0] = ""
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		log.Println("error: method not allowed")
	}
}
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	if r.Method != http.MethodGet {
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("error: nil cookies")
			h.Errors(w, http.StatusUnauthorized, "Error in cookie")
			return
		}
		log.Println("error: can't get cookie")
		h.Errors(w, http.StatusBadRequest, "Error in cookie")
		return
	}
	if err := h.services.Auth.DeleteSessionToken(c.Value); err != nil {
		h.Errors(w, http.StatusInternalServerError, err.Error())
		log.Println("error: can't delete session token")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
