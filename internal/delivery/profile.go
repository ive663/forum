package delivery

// import (
// 	"log"
// 	"net/http"
// 	"text/template"
// )

// func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	user_is_authoryz, ok := r.Context().Value(keyUserBool).(bool)
// 	if !ok {
// 		h.Errors(w, http.StatusForbidden, "not authorized")
// 		return
// 	}
// 	if r.Method == "GET" {
// 		t, err := template.ParseFiles("./templates/profile.html")
// 		if err != nil {
// 			log.Print(err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		cookie, err := r.Cookie("session")
// 		if err!= nil {
//             log.Print(err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		user, err := h.services.ParseSessionToken(cookie.Value)
// 		if err!= nil {
//             log.Print(err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		if user_is_authoryz {
// 			module.G

// 		}
// 	}
// }
