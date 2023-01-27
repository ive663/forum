package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ive663/forum/internal/module"
)

type Post interface {
	CreatePost(*module.Post) (int, error)
	CreateCategory(*module.Category) error
	GetPostByCategory(category string) ([]module.Post, error)
	GetOldPosts() ([]module.Post, error)
	GetNewPosts() ([]module.Post, error)
	GetPostIdByUserId(id int) (*module.Post, error)
	GetPostByPostId(id int) (*module.Post, error)
	GetAllCategoryByPostId(postid int) ([]module.Category, error)
	GetPostsByUserId(id int) ([]module.Post, error)

	///  added new interfaces for likes and dislikes ///
	GetLikesCountByPostID(postID int) (*module.Post, error)
	GetDisLikesCountByPostID(postID int) (*module.Post, error)
	AddLikeByPost(postID int, userID int) error
	AddDisLikeByPost(postID int, userID int) error
	RemoveLikeByPost(postID int, userID int) error
	RemoveDislikeByPost(postID int, userId int) error
	///=================///
	PostHasLike(postId int, userId int) error
	PostHasDisLike(postId int, userId int) error
	GetMyLikedPosts(userID int) ([]module.Post, error)
	///=================///
	GetAllPostsByUserId(id int) ([]module.Post, error)
	/// added new interfaces for sorting by likes///
	GetPostsByLikesHigh() ([]module.Post, error)
	GetPostsByLikesLow() ([]module.Post, error)
	///=================///
	GetPostsByDisLikesHigh() ([]module.Post, error)
	GetPostsByDisLikesLow() ([]module.Post, error)
}

var ErrRecordNotFound = errors.New("record not found")

type PostRepository struct {
	db *sql.DB
}

func newPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

/// added new interfaces for sorting by likes///

func (r *PostRepository) GetPostsByDisLikesLow() ([]module.Post, error) {
	var posts []module.Post
	rows, err := r.db.Query("SELECT id, title, author_id, message, likes, dislikes, date FROM posts ORDER BY dislikes ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Message, &post.AuthorID, &post.Likes, &post.Dislikes, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByDisLikesHigh() ([]module.Post, error) {
	var posts []module.Post
	rows, err := r.db.Query("SELECT id, title, author_id, message, likes, dislikes, date FROM posts ORDER BY dislikes DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Message, &post.AuthorID, &post.Likes, &post.Dislikes, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByLikesLow() ([]module.Post, error) {
	var posts []module.Post
	rows, err := r.db.Query("SELECT id, title, author_id, message, likes, dislikes, date FROM posts ORDER BY likes ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Message, &post.Likes, &post.Dislikes, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByLikesHigh() ([]module.Post, error) {
	var posts []module.Post
	rows, err := r.db.Query("SELECT id, title, author_id, message,  likes, dislikes, date FROM posts ORDER BY likes DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Message, &post.Likes, &post.Dislikes, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

///=================///

// Get all posts by user id
func (r *PostRepository) GetAllPostsByUserId(id int) ([]module.Post, error) {
	var posts []module.Post
	rows, err := r.db.Query("SELECT id, title, message, author_id, likes, dislikes, date FROM posts WHERE author_id = ?", id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := module.Post{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Message, &p.AuthorID, &p.Likes, &p.Dislikes, &p.Date); err != nil {
			log.Print(err)
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetMyLikedPosts(userID int) ([]module.Post, error) {
	var posts []module.Post
	queryLike := "SELECT post_id FROM likes WHERE user_id = ?"
	queryPosts := "SELECT * FROM posts WHERE id = ?"
	rowsLike, err := r.db.Query(queryLike, userID)
	if err != nil {
		return nil, err
	}
	defer rowsLike.Close()
	for rowsLike.Next() {
		var postid int
		if err := rowsLike.Scan(&postid); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			if postid == 0 {
				continue
			}
			return nil, err
		}
		rowsPosts, err := r.db.Query(queryPosts, postid)
		if err != nil {
			return nil, err
		}
		defer rowsPosts.Close()
		for rowsPosts.Next() {
			p := module.Post{}
			if err := rowsPosts.Scan(&p.ID, &p.Title, &p.AuthorID, &p.Author, &p.Message, &p.Likes, &p.Dislikes, &p.CategoryID, &p.Date); err != nil {
				return nil, err
			}
			posts = append(posts, p)
		}
	}
	return posts, nil
}

// add like to post by post id and return error
func (r *PostRepository) AddLikeByPost(postID int, userID int) error {
	query := "INSERT INTO likes(post_id, user_id) VALUES (?, ?)"
	if _, err := r.db.Exec(query, postID, userID); err != nil {
		return err
	}
	query = "UPDATE posts SET likes = likes + 1 WHERE id = ?;"
	if _, err := r.db.Exec(query, postID); err != nil {
		return err
	}
	return nil
}

// add dislike to post by post id and return  error
func (r *PostRepository) AddDisLikeByPost(postID int, userID int) error {
	query := "INSERT INTO dislikes(post_id, user_id) VALUES (?, ?)"
	if _, err := r.db.Exec(query, postID, userID); err != nil {
		return err
	}
	query = "UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?;"
	if _, err := r.db.Exec(query, postID); err != nil {
		return err
	}
	return nil
}

// remove dislike to post by post id and return error
func (r *PostRepository) RemoveDislikeByPost(postID int, userID int) error {
	query := "DELETE FROM dislikes WHERE post_id = ? AND user_id = ?"
	if _, err := r.db.Exec(query, postID, userID); err != nil {
		return err
	}
	query = "UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?;"
	if _, err := r.db.Exec(query, postID); err != nil {
		return err
	}
	return nil
}

// remove like to post by post id and return error
func (r *PostRepository) RemoveLikeByPost(postID int, userID int) error {
	query := "DELETE FROM likes WHERE post_id = ? AND user_id = ?"
	if _, err := r.db.Exec(query, postID, userID); err != nil {
		return err
	}
	query = "UPDATE posts SET likes = likes - 1 WHERE id = ?;"
	if _, err := r.db.Exec(query, postID); err != nil {
		return err
	}
	return nil
}

// get likes count by post id and return error
func (r *PostRepository) GetLikesCountByPostID(postID int) (*module.Post, error) {
	var post module.Post
	query := "SELECT likes FROM posts WHERE id = ?;"
	if err := r.db.QueryRow(query, postID).Scan(&post.Likes); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) PostHasLike(postId int, userId int) error {
	var u int
	query := "SELECT user_id FROM likes WHERE post_id = ? AND user_id = ?"
	err := r.db.QueryRow(query, postId, userId).Scan(&u)
	if err != nil {
		return fmt.Errorf("error in posthaslike, %w", err)
	}
	return nil
}

func (r *PostRepository) PostHasDisLike(postId int, userId int) error {
	var u int
	query := "SELECT user_id FROM dislikes WHERE post_id = ? AND user_id = ?"
	err := r.db.QueryRow(query, postId, userId).Scan(&u)
	if err != nil {
		return fmt.Errorf("error in posthasdislike, %w", err)
	}
	return nil
}

// get dislikes count by post id and return error
func (r *PostRepository) GetDisLikesCountByPostID(postID int) (*module.Post, error) {
	var post module.Post
	query := "SELECT dislikes FROM posts WHERE id = ?;"
	if err := r.db.QueryRow(query, postID).Scan(&post.Dislikes); err != nil {
		return nil, err
	}
	return &post, nil
}

///===================================================///

func (r *PostRepository) CreatePost(p *module.Post) (int, error) {
	query := "INSERT INTO posts(title, author_id, author, message, category_id, date) VALUES (?, ?, ?, ?, ?, ?) RETURNING id"
	var id int
	if err := r.db.QueryRow(query, p.Title, p.AuthorID, p.Author, p.Message, p.CategoryID, p.Date).Scan(&id); err != nil {
		log.Print(err)
		return 0, err
	}
	return id, nil
}

func (r *PostRepository) CreateCategory(c *module.Category) error {
	if _, err := r.db.Exec("INSERT INTO categories (tag, postid) VALUES(?, ?)", c.Tag, c.PostID); err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetNewPosts() ([]module.Post, error) {
	var posts []module.Post
	query := "SELECT id, title, author_id, author, message, likes, dislikes, date FROM posts ORDER by date DESC;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Author, &post.Message, &post.Likes, &post.Dislikes, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)

	}
	return posts, nil
}

func (r *PostRepository) GetPostByCategory(category string) ([]module.Post, error) {
	var posts []module.Post
	query := "SELECT * FROM posts WHERE id IN (SELECT postid FROM categories WHERE tag = ?);"
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	for rows.Next() {
		var post module.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Author, &post.Message, &post.Likes, &post.Dislikes, &post.CategoryID, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1 / /
func (r *PostRepository) GetOldPosts() ([]module.Post, error) {
	var posts []module.Post
	query := "SELECT id, title, author_id, author, message, likes, dislikes, date FROM posts ORDER by date;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting old posts: %w", err)
	}
	for rows.Next() {
		post := module.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Author, &post.Message, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)

	}
	return posts, nil
}

func (r *PostRepository) GetPostsByUserId(id int) ([]module.Post, error) {
	var posts []module.Post
	query := "SELECT * FROM posts WHERE author_id = ?"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post module.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Author, &post.Message, &post.Likes, &post.Dislikes, &post.CategoryID, &post.Date); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostIdByUserId(id int) (*module.Post, error) {
	var p module.Post
	err := r.db.QueryRow("SELECT id FROM posts WHERE author_id = ?", id).Scan(
		&p.AuthorID)
	if err == sql.ErrNoRows {
		return nil, errors.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepository) GetPostByPostId(postid int) (*module.Post, error) {
	p := &module.Post{}
	err := r.db.QueryRow("SELECT id, title, author_id, author, message, date FROM posts WHERE id = ?", postid).Scan(&p.ID, &p.Title, &p.AuthorID, &p.Author, &p.Message, &p.Date)
	if err == sql.ErrNoRows {
		return nil, errors.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepository) GetAllCategoryByPostId(postid int) ([]module.Category, error) {
	queryCategory := "SELECT tag FROM categories WHERE postid = ?"
	categoryRows, err := r.db.Query(queryCategory, postid)
	if err != nil {
		return nil, err
	}
	var category []module.Category
	for categoryRows.Next() {
		var oneCategory module.Category
		if err := categoryRows.Scan(&oneCategory.Tag); err != nil {
			return nil, err
		}
		category = append(category, oneCategory)
	}
	return category, nil
}
