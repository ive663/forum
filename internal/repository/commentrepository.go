package repository

import (
	"database/sql"
	"log"

	"github.com/ive663/forum/internal/module"
)

type Comment interface {
	CreateComment(*module.Comment) error
	FindCommentsInPostID(postid int) ([]module.Comment, error)
	GetPostIdByCommentId(commentID int) (*module.Comment, error)
	GetCommentLikesByPostID(postID int) (map[int][]int, error)
	GetCommentDislikesByPostID(postID int) (map[int][]int, error)
	AddLikeByComment(commentID int, userdID int) error
	AddDislikeByComment(commentID int, userID int) error
	RemoveLikeByComment(commentID int, userID int) error
	RemoveDislikeByComment(commentID int, userID int) error
	CommentHasLike(commentID int, userID int) error
	CommentHasDislike(commentID int, userID int) error
}

type CommentRepository struct {
	db *sql.DB
}

func newCommentRepostiroy(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) GetPostIdByCommentId(commentID int) (*module.Comment, error) {
	c := &module.Comment{}
	err := r.db.QueryRow("SELECT post_id FROM comments WHERE id = ?", commentID).Scan(&c.PostID)
	if err == sql.ErrNoRows {
		log.Println("error:rep: no rows found in GetPostIdByCommentId")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	log.Println(c.PostID)
	return c, nil
}

// /=================///
func (r *CommentRepository) GetCommentLikesByPostID(postID int) (map[int][]int, error) {
	queryForCommentsId := "SELECT id FROM comments WHERE post_id = ?"
	queryForLikes := "SELECT likes FROM comments WHERE id = ?"
	commentLikes := make(map[int][]int)
	rowsComment, err := r.db.Query(queryForCommentsId, postID)
	if err != nil {
		return nil, err
	}
	defer rowsComment.Close()
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		var likes []int
		rowsLikes, err := r.db.Query(queryForLikes, id)
		if err != nil {
			return nil, err
		}
		defer rowsLikes.Close()
		for rowsLikes.Next() {
			var like int
			if err := rowsLikes.Scan(&like); err != nil {
				return nil, err
			}
			likes = append(likes, like)
		}
		commentLikes[id] = likes
	}
	return commentLikes, nil
}

func (r *CommentRepository) GetCommentDislikesByPostID(postID int) (map[int][]int, error) {
	queryForCommentsId := "SELECT id FROM comments WHERE post_id = ?"
	queryForDislikes := "SELECT dislikes FROM comments WHERE id = ?"
	commentDislikes := make(map[int][]int)
	rowsComment, err := r.db.Query(queryForCommentsId, postID)
	if err != nil {
		return nil, err
	}
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		var dislikes []int
		rowsDislikes, err := r.db.Query(queryForDislikes, id)
		if err != nil {
			return nil, err
		}
		for rowsDislikes.Next() {
			var dislike int
			if err := rowsDislikes.Scan(&dislike); err != nil {
				return nil, err
			}
			dislikes = append(dislikes, dislike)
		}
		commentDislikes[id] = dislikes
	}
	return commentDislikes, nil
}

func (r *CommentRepository) AddLikeByComment(commentID int, userID int) error {
	query := "INSERT INTO likes(user_id, comment_id) VALUES (?, ?)"
	if _, err := r.db.Exec(query, userID, commentID); err != nil {
		return err
	}
	query = "UPDATE comments SET likes = likes + 1 WHERE id = ?"
	if _, err := r.db.Exec(query, commentID); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *CommentRepository) AddDislikeByComment(commentID int, userID int) error {
	query := "INSERT INTO dislikes(user_id, comment_id) VALUES (?, ?)"
	if _, err := r.db.Exec(query, userID, commentID); err != nil {
		return err
	}
	query = "UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?"
	if _, err := r.db.Exec(query, commentID); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *CommentRepository) RemoveLikeByComment(commentID int, userID int) error {
	query := "DELETE FROM likes WHERE user_id = ? AND comment_id = ?"
	if _, err := r.db.Exec(query, userID, commentID); err != nil {
		return err
	}
	query = "UPDATE comments SET likes = likes - 1 WHERE id = ?"
	if _, err := r.db.Exec(query, commentID); err != nil {
		log.Print(err)
		return err
	}
	log.Println("Remove ", userID, "'s like from comment ", commentID)
	return nil
}

func (r *CommentRepository) RemoveDislikeByComment(commentID int, userID int) error {
	log.Println("dislike removing...")
	query := "DELETE FROM dislikes WHERE user_id = ? AND comment_id = ?"
	if _, err := r.db.Exec(query, userID, commentID); err != nil {
		log.Println("error db delete dislike: ", err)
		return err
	}
	log.Println("dislike removing update...")
	query = "UPDATE comments SET dislikes = dislikes - 1 WHERE id = ?"
	if _, err := r.db.Exec(query, commentID); err != nil {
		log.Println("error update db dislike: ", err)
		return err
	}
	log.Println("Remove ", userID, "'s dislike from comment ", commentID)
	return nil
}

func (r *CommentRepository) CommentHasLike(commentID int, userID int) error {
	var u int
	query := "SELECT user_id FROM likes WHERE comment_id = ? AND user_id = ?"
	err := r.db.QueryRow(query, commentID, userID).Scan(&u)
	if err != nil {
		log.Println("error:rep: no rows found in CommentHasLike")
		return err
	}
	return nil
}

func (r *CommentRepository) CommentHasDislike(commentID int, userID int) error {
	var u int
	query := "SELECT user_id FROM dislikes WHERE comment_id = ? AND user_id = ?"
	err := r.db.QueryRow(query, commentID, userID).Scan(&u)
	if err != nil {
		log.Println("error:rep: no rows found in CommentHasDislike")
		return err
	}
	log.Print("Comment ", commentID, " has dislike")
	return nil
}

func (r *CommentRepository) CreateComment(c *module.Comment) error {
	if _, err := r.db.Exec("INSERT INTO comments (author_id, author, post_id, message, date) VALUES(?, ?, ?, ?, ?)", c.AuthorID, c.Author, c.PostID, c.Message, c.Date); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *CommentRepository) FindCommentsInPostID(PostId int) ([]module.Comment, error) {
	u := module.User{}
	comments := u.Comments
	c := module.Comment{}
	rows, err := r.db.Query("SELECT id, author, message, date, likes, dislikes FROM comments WHERE post_id = ?", PostId)
	if err == sql.ErrNoRows {
		log.Println("error:rep: no rows found in FindCommentsInPostID")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Author, &c.Message, &c.Date, &c.Likes, &c.Dislikes); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
