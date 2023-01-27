package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ive663/forum/internal/module"
	"github.com/ive663/forum/internal/repository"
)

var ErrInvalidComment = errors.New("Invalid typing comment")

type Comment interface {
	GetComments(postId int) (module.CommentList, error)
	CreateComment(comment *module.Comment) error
	///  added new interfaces for likes and dislikes ///
	GetCommentLikesByPostID(postID int) (map[int][]int, error)
	GetCommentDislikesByPostID(postID int) (map[int][]int, error)
	AddLikeByComment(commentID int, userID int) error
	AddDislikeByComment(commentID int, userID int) error
	GetPostIdByCommentId(commentID int) (*module.Comment, error)
}

type CommentService struct {
	repository repository.Comment
}

func newCommentService(repository repository.Comment) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

///===============================================///
///  added new interfaces for likes and dislikes ///
///=============================================///

// GetCommentLikesPostID
func (s *CommentService) GetCommentLikesByPostID(postID int) (map[int][]int, error) {
	comments, err := s.repository.GetCommentLikesByPostID(postID)
	if err != nil {
		log.Println("error:service:comment: GetCommentLikesByPostID")
		return nil, err
	}
	return comments, nil
}

// GetCommentDisLikesByPostID
func (s *CommentService) GetCommentDislikesByPostID(postID int) (map[int][]int, error) {
	comments, err := s.repository.GetCommentDislikesByPostID(postID)
	if err != nil {
		log.Println("error:service:comment: GetCommentDislikesByPostID")
		return nil, err
	}
	return comments, nil
}

// AddLikeByCommentID
func (s *CommentService) AddLikeByComment(commentID int, userID int) error {
	if err := s.repository.CommentHasLike(commentID, userID); err == nil {
		err := s.repository.RemoveLikeByComment(commentID, userID)
		return err

	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Println(err.Error())
		return err
	}
	if err := s.repository.CommentHasDislike(commentID, userID); err == nil {
		if err := s.repository.RemoveDislikeByComment(commentID, userID); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.AddLikeByComment(commentID, userID); err != nil {
		return err
	}
	return nil
}

// AddDisLikeByCommentID
func (s *CommentService) AddDislikeByComment(commentID int, userID int) error {
	if err := s.repository.CommentHasDislike(commentID, userID); err == nil {
		log.Println("Service comment has dislike", err)
		if err := s.repository.RemoveDislikeByComment(commentID, userID); err != nil {
			log.Println("Removing dislike error: ", err)
			return err
		}
		log.Print("Service comment removed dislike")
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.CommentHasLike(commentID, userID); err == nil {
		if err := s.repository.RemoveLikeByComment(commentID, userID); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.AddDislikeByComment(commentID, userID); err != nil {
		return err
	}
	return nil
}

///=======================================///

func (s *CommentService) GetComments(postId int) (module.CommentList, error) {
	comments, err := s.repository.FindCommentsInPostID(postId)
	if err != nil {
		log.Println("error:service:comment: GetComments")
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) CreateComment(comment *module.Comment) error {
	if err := ValidComment(comment); err != nil {
		log.Println("error:service:comment: CreateComment: ValidComment")
		return err
	}
	if err := s.repository.CreateComment(comment); err != nil {
		log.Println("error:service:comment:CreateComment: repo.CreateComment")
		return err
	}
	return nil
}

func (s *CommentService) GetPostIdByCommentId(commentID int) (*module.Comment, error) {
	c, err := s.repository.GetPostIdByCommentId(commentID)
	if err != nil {
		log.Println("error:service:comment: GetPostIdByCommentId")
		return nil, err
	}
	return c, nil
}

func ValidComment(comment *module.Comment) error {
	onlyspace := true
	for _, check := range comment.Message {
		if check < 32 || check > 127 {
			return ErrInvalidComment
		}
		if check != ' ' {
			onlyspace = false
			break
		}
	}
	if onlyspace {
		return ErrEmptyValue
	}
	return nil
}

func GetCommentId(comments []module.Comment) int {
	for _, comment := range comments {
		return comment.ID
	}
	return 0
}
