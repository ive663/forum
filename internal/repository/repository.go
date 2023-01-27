package repository

import "database/sql"

type Repository struct {
	Post
	Comment
	Auth
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Post:    newPostRepository(db),
		Comment: newCommentRepostiroy(db),
		Auth:    newAuthRepository(db),
	}
}
