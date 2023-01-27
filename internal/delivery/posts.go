package delivery

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ive663/forum/internal/module"
	"github.com/ive663/forum/internal/service"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		h.Errors(w, http.StatusForbidden, "You can't post comment.")
		return
	}
	user_authorization := true
	if user_id == 0 {
		user_authorization = false
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/post.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		post, err := h.services.GetPostByPostId(postid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.Errors(w, http.StatusNotFound, err.Error())
				return
			}
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		postlikes, err := h.services.GetLikesCountByPostID(postid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		postdislikes, err := h.services.GetDisLikesCountByPostID(postid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		comment, err := h.services.GetComments(post.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.Errors(w, http.StatusNotFound, err.Error())
				return
			}
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentlikes, err := h.services.GetCommentLikesByPostID(postid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		commentdislikes, err := h.services.GetCommentDislikesByPostID(postid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		pageContent := module.PostPage{
			Post:             post,
			PostLikes:        postlikes.Likes,
			PostDislikes:     postdislikes.Dislikes,
			CommentsLikes:    commentlikes,
			CommentsDislikes: commentdislikes,
			Comments:         comment.PrepToView(),
			Authorization:    user_authorization,
		}

		if err := t.Execute(w, pageContent); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "POST":
		if user_id != 0 {
			err := r.ParseForm()
			if err != nil {
				h.Errors(w, http.StatusBadRequest, "Error parsing")
				return
			}
			comment := r.Form["comment"]
			if !ok {
				log.Print("Error in input comment")
				h.Errors(w, http.StatusBadRequest, "Bad typing message")
				return
			}
			author, err := h.services.GetUserByUserID(user_id)
			if err != nil {
				h.Errors(w, http.StatusInternalServerError, err.Error())
				return
			}
			newComment := &module.Comment{
				AuthorID: user_id,
				Author:   author.Login,
				PostID:   postid,
				Message:  comment[0],
				Date:     time.Now(),
			}
			if err := h.services.Comment.CreateComment(newComment); err != nil {
				if errors.Is(err, service.ErrInvalidComment) || errors.Is(err, service.ErrEmptyValue) {
					h.Errors(w, http.StatusBadRequest, err.Error())
					return
				}
			}
			http.Redirect(w, r, "post?id="+strconv.Itoa(postid), http.StatusSeeOther)
		}
		tmpl, err := template.ParseFiles("./templates/post.html")
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, "Error in parsing")
			log.Print(err)
			return
		}
		if err = tmpl.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error in execute")
		}
	}
}

func (h *Handler) createpost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createpost" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	user_id, ok := r.Context().Value(keyUserID).(int)
	if !ok {
		log.Print("Not ok")
		return
	}
	if user_id == 0 {
		h.Errors(w, http.StatusForbidden, "")
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/createpost.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error parsing file")
			return
		}
		if err = t.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error executing")
			return
		}
	case "POST":
		user, err := h.services.GetUserByUserID(user_id)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err = r.ParseForm(); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusBadRequest, "Error parsing")
			return
		}
		title, ok := r.Form["title"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Bad typing title")
			return
		}
		message, ok := r.Form["message"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Bad typing message")
			return
		}
		titlecategory, ok := r.Form["category"]
		tags := strings.Fields(strings.Join(titlecategory, " "))
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Bad typing title category")
		}
		newPost := &module.Post{
			Title:    title[0],
			Message:  message[0],
			AuthorID: user.ID,
			Author:   user.Login,
			Date:     time.Now(),
		}
		err = h.services.CreatePost(newPost, tags)
		if err != nil {
			if errors.Is(err, service.ErrEmptyValue) || errors.Is(err, service.ErrInvalidTypingPost) {
				h.Errors(w, http.StatusBadRequest, err.Error())
				return
			}
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		t, err := template.ParseFiles("templates/createpost.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error Parsing")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		if err = t.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, "Error executing")
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
