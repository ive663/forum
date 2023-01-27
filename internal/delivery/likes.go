package delivery

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't like post.")
		return
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}

	switch r.Method {
	case "GET":
		if err := h.services.Post.AddLikeByPost(postid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "post?id="+strconv.Itoa(postid), 303)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
func (h *Handler) likePostIndex(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't like post.")
		return
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}

	switch r.Method {
	case "GET":
		if err := h.services.Post.AddLikeByPost(postid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	commentid, err := strconv.Atoi(r.URL.Query().Get("commentid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't like commnet")
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}
	switch r.Method {
	case "GET":
		if err := h.services.Comment.AddLikeByComment(commentid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		postid, err := h.services.Comment.GetPostIdByCommentId(commentid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "post?id="+strconv.Itoa(postid.PostID), 303)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}

func (h *Handler) dislikePost(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't like post.")
		return
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}

	switch r.Method {
	case "GET":
		if err := h.services.Post.AddDislikeByPost(postid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "post?id="+strconv.Itoa(postid), 303)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}

func (h *Handler) dislikePostIndex(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't like post.")
		return
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}

	switch r.Method {
	case "GET":
		if err := h.services.Post.AddDislikeByPost(postid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}

func (h *Handler) dislikeComment(w http.ResponseWriter, r *http.Request) {
	commentid, err := strconv.Atoi(r.URL.Query().Get("commentid"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't dislike commnet")
	}
	if user_id == 0 {
		http.Redirect(w, r, "/signin", 303)
	}
	switch r.Method {
	case "GET":
		if err := h.services.Comment.AddDislikeByComment(commentid, user_id); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		postid, err := h.services.Comment.GetPostIdByCommentId(commentid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Println(postid)
		http.Redirect(w, r, "post?id="+strconv.Itoa(postid.PostID), 303)
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
