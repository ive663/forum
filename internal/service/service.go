package service

import "github.com/ive663/forum/internal/repository"

type Service struct {
	Auth
	Post
	Comment
}

func NewServices(repositories *repository.Repository) *Service {
	return &Service{
		Auth:    newAuthService(repositories.Auth),
		Post:    newPostService(repositories.Post),
		Comment: newCommentService(repositories.Comment),
	}
}
