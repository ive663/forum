package service

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/ive663/forum/internal/repository"

	"github.com/ive663/forum/internal/module"
)

var (
	ErrInvalidQueryRequest   = errors.New("invalid query request")
	ErrEmptyValue            = errors.New("Empty value")
	ErrInvalidTypingPost     = errors.New("Invalid typing post")
	ErrInvalidTypingCategory = errors.New("Invalid typing category")
)

type Post interface {
	CreatePost(post *module.Post, category []string) error
	CreateCategory(category *module.Category) error
	GetNewPosts() (module.PostList, error)
	GetAllPostBy(userid int, query map[string][]string) (module.PostList, error)
	GetPostIdByUserId(id int) (*module.Post, error)
	GetPostByPostId(id int) (*module.Post, error)

	///  added new interfaces for likes and dislikes ///
	GetLikesCountByPostID(postID int) (*module.Post, error)
	GetDisLikesCountByPostID(postID int) (*module.Post, error)
	AddLikeByPost(postID int, userID int) error
	AddDislikeByPost(postID int, userID int) error
	///=================///
}

type PostService struct {
	repository repository.Post
}

func newPostService(repository repository.Post) *PostService {
	return &PostService{
		repository: repository,
	}
}

// /  added new interfaces for likes and dislikes ///
// GetLikesCountByPostID returns the number of likes for a post and error
func (s *PostService) GetLikesCountByPostID(postID int) (*module.Post, error) {
	post, err := s.repository.GetLikesCountByPostID(postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetDisLikesCountByPostID returns the number of dislikes for a post and error
func (s *PostService) GetDisLikesCountByPostID(postID int) (*module.Post, error) {
	post, err := s.repository.GetDisLikesCountByPostID(postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// AddLikeByPostID adds a like to a post and returns an error
func (s *PostService) AddLikeByPost(postID int, userID int) error {
	if err := s.repository.PostHasLike(postID, userID); err == nil {
		if err := s.repository.RemoveLikeByPost(postID, userID); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.PostHasDisLike(postID, userID); err == nil {
		if err := s.repository.RemoveDislikeByPost(postID, userID); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.AddLikeByPost(postID, userID); err != nil {
		return err
	}
	return nil
}

// AddDisLikeByPostID adds a dislike to a post and returns an error
func (s *PostService) AddDislikeByPost(postID int, userID int) error {
	if err := s.repository.PostHasDisLike(postID, userID); err == nil {
		if err := s.repository.RemoveDislikeByPost(postID, userID); err != nil {
			return err
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.PostHasLike(postID, userID); err == nil {
		if err := s.repository.RemoveLikeByPost(postID, userID); err != nil {
			return err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := s.repository.AddDisLikeByPost(postID, userID); err != nil {
		return err
	}
	return nil
}

///===============================================///

func (s *PostService) CreatePost(post *module.Post, categories []string) error {
	err := validPost(post)
	if err != nil {
		return err
	}

	id, err := s.repository.CreatePost(post)
	if err != nil {
		log.Println("error:service:post:CreatePost:", err)
		return err
	}
	for _, category := range categories {
		c := &module.Category{
			PostID: id,
			Tag:    category,
		}

		if err := s.repository.CreateCategory(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) CreateCategory(category *module.Category) error {
	err := validCategory(category)
	if err != nil {
		return err
	}
	err = s.repository.CreateCategory(category)
	if err != nil {
		log.Println("error:service:post:CreateCategory:", err)
		return err
	}
	return nil
}

func (s *PostService) GetNewPosts() (module.PostList, error) {
	posts, err := s.repository.GetNewPosts()
	if err != nil {
		log.Println("error:service:post:GetNewPosts:", err)
		return nil, err
	}
	for i := range posts {
		category, err := s.repository.GetAllCategoryByPostId(posts[i].ID)
		likes, err := s.repository.GetLikesCountByPostID(posts[i].ID)
		dislikes, err := s.repository.GetDisLikesCountByPostID(posts[i].ID)
		if err != nil {
			log.Println("error:service:post:GetNewPosts:", err)
			return nil, err
		}
		posts[i].Categories = category
		posts[i].Likes = likes.Likes
		posts[i].Dislikes = dislikes.Dislikes
	}
	return posts, nil
}

func (s PostService) GetPostIdByUserId(id int) (*module.Post, error) {
	post, err := s.repository.GetPostIdByUserId(id)
	if err != nil {
		log.Println("error:service:post:GetPostIdInUserId:", err)
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByPostId(id int) (*module.Post, error) {
	p, err := s.repository.GetPostByPostId(id)
	if err != nil {
		log.Println("error:service:post:GetPostInPostId:", err)
		return nil, err
	}
	return p, nil
}

func validPost(post *module.Post) error {
	whiteSpaceTitle := true
	for _, title := range post.Title {
		if title < 32 || title > 127 {
			return ErrInvalidTypingPost
		}
		if title != ' ' {
			whiteSpaceTitle = false
			break
		}
	}
	if whiteSpaceTitle {
		return ErrEmptyValue
	}
	whiteSpaceMessage := true
	for _, message := range post.Message {
		if message < 32 || message > 127 {
			return ErrInvalidTypingPost
		}
		if message != ' ' {
			whiteSpaceMessage = false
			break
		}
	}
	if whiteSpaceMessage {
		return ErrEmptyValue
	}
	return nil
}

func validCategory(category *module.Category) error {
	whiteSpaceTag := true
	for _, tag := range category.Tag {
		if tag < 32 || tag > 127 {
			return ErrInvalidTypingPost
		}
		if tag != ' ' {
			whiteSpaceTag = false
			break
		}
	}
	if whiteSpaceTag || category.Tag == "" {
		return ErrEmptyValue
	}
	return nil
}

func (s *PostService) GetAllPostBy(userid int, query map[string][]string) (module.PostList, error) {
	var (
		posts []module.Post
		err   error
	)
	if userid == 0 {
		for key, val := range query {
			switch key {
			case "category":
				posts, err = s.repository.GetPostByCategory(strings.Join(val, ""))
				if err != nil {
					log.Println("error:post:GetPostByCategory: ", err)
					return nil, err
				}
				if err != nil {
					log.Println("error:post:GetAllPostBy:time: ", err)
					return nil, err
				}
			}
		}
	}
	for key, val := range query {
		switch key {
		case "category":
			posts, err = s.repository.GetPostByCategory(strings.Join(val, ""))
			if err != nil {
				if errors.Is(err, repository.ErrRecordNotFound) {
					return nil, ErrInvalidQueryRequest
				}
				log.Println("error:post:GetPostByCategory: ", err)
				return nil, err
			}
		case "mypost":
			switch strings.Join(val, "") {
			case "mypost":
				posts, err = s.repository.GetPostsByUserId(userid)
				if err != nil {
					log.Println("error:post:GetAllPostBy:mypost: ", err)
					return nil, err
				}
			default:
				return nil, ErrInvalidQueryRequest
			}

		case "mylikedposts":
			switch strings.Join(val, "") {
			case "mylikedposts":
				posts, err = s.repository.GetMyLikedPosts(userid)
				if err != nil {
					log.Println("error:post:GetAllPostBy:likedpost: ", err)
					return nil, err
				}
			default:
				return nil, ErrInvalidQueryRequest
			}
		default:
			log.Println("error:post:GetAllPostBy:default: ", err)
			return nil, ErrInvalidQueryRequest
		}
	}
	for i := range posts {
		log.Println(posts[i].ID)
		category, err := s.repository.GetAllCategoryByPostId(posts[i].ID)
		if err != nil {
			log.Println("error:post:GetAllPostBy:category: ", err)
			return nil, err
		}
		posts[i].Categories = category
	}
	return posts, nil
}
