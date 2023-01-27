package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/ive663/forum/internal/module"
	"github.com/ive663/forum/internal/repository"

	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("Incorrect login or password")
	ErrInvalidUserName = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type Auth interface {
	CreateNewUser(user *module.User) (newUser *module.User, err error)
	GenerateSessionToken(login, password string) (string, error)
	ParseSessionToken(token string) (*module.User, error)
	DeleteSessionToken(token string) error
	GetUserIdByUUID(token string) (int, error)
	GetUserByUserID(id int) (*module.User, error)
	DeleteExpiredSessions() error
}

type AuthService struct {
	repository repository.Auth
}

func newAuthService(repository repository.Auth) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func BeforeCreate(u *module.User) error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			log.Println("Error:service:auth:BeforeCreate: encryptString")
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

func Sanitize(u *module.User) {
	u.Password = ""
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error:service:auth:encryptString")
		return "", err
	}
	return string(b), nil
}

func ComparePassword(u *module.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) != nil
}

func validUser(u *module.User) error {
	for _, char := range u.Login {
		if char < 32 || char > 127 {
			log.Println("Error:service:auth:validUser: invalid username")
			return ErrInvalidUserName
		}
	}
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, u.Email)
	// validEmail, err := mail.ParseAddress(u.Email)
	if err != nil {
		log.Println("Error:service:auth:validUser: ParseAddress")
		return err
	}
	if !validEmail {
		log.Println("Error:service:auth:validUser: invalid email")
		return ErrInvalidEmail
	}
	if len(u.Login) < 4 || len(u.Login) > 36 {
		log.Println("Error:service:auth:validUser: invalid login")
		return ErrInvalidUserName
	}
	return nil
}

func (s *AuthService) GenerateSessionToken(username, password string) (string, error) {
	log.Println("service:auth:GenerateSession: password: ", password)
	user, err := s.repository.FindByLogin(username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error:service:auth:GenerateSessionToken: FindByLogin")
			return "", ErrUserNotFound
		}
		log.Println("Error:service:auth:GenerateSessionToken: FindByLogin: ", err)
		return "", err
	}
	if ComparePassword(user, password) {
		log.Println(user, password)
		log.Println("Error:service:auth:GenerateSessionToken: ComparePassword: ", err)
		return "", ErrUserNotFound
	}
	token := uuid.NewV4()
	fmt.Println("TOKEN: " + token.String())
	session := &module.Session{
		UserID:    user.ID,
		UUID:      token.String(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}
	isSessionExists, err := s.repository.IsSessionExists(user.ID)
	if err != nil {
		log.Println("Error:service:auth: IsSessionExists: ", err)
		return "", err
	}

	if isSessionExists {
		log.Println("service:GenSessionToken: updating session token")
		err = s.repository.UpdateSession(session)
		if err != nil {
			log.Println("Error:service:auth: updateSession: ", err)
			return "", err
		}
	} else {
		err = s.repository.CreateNewSession(session)
		if err != nil {
			log.Println("Error:Service:Auth: CreateNewSession: ", err)
			return "", err
		}
	}
	return token.String(), nil
}

func (s *AuthService) CreateNewUser(user *module.User) (*module.User, error) {
	err := validUser(user)
	if err != nil {
		log.Println("Error:service:auth:CreateNewUser validUser: ", err)
		return nil, err
	}
	user.EncryptedPassword, err = encryptString(user.Password)
	log.Println("service:auth:CreateNewUser: encrypted password: ", user.EncryptedPassword)
	if err != nil {
		log.Println("error:service:auth:CreateNewUser: encrypted password: ", err)
		return nil, err
	}
	err = s.repository.CreateNewUser(user)
	log.Println("service:auth:CreateNewUser: create new user: ", user)
	if err != nil {
		log.Println("error:service:auth:CreateNewUser: create new user: ", err)
		return nil, err
	}
	newUser, err := s.repository.FindByLogin(user.Login)
	log.Println("service:auth:CreateNewUser: find by login: ", newUser)
	if err != nil {
		log.Println("Error:service:auth:CreateNewUser: FindByLogin: ", err)
		return nil, err
	}
	return newUser, nil
}

func (s *AuthService) ParseSessionToken(token string) (*module.User, error) {
	uid, err := s.repository.GetUserIdByUUID(token)
	if err != nil {
		log.Println("Error:service:auth:ParseSessionToken: GetUUIDFindID")
		return nil, err
	}
	user, err := s.repository.GetUserByID(uid)
	if err != nil {
		log.Println("Error:service:auth:ParseSessionToken: FindIdUser")
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetUserIdByUUID(token string) (int, error) {
	if token == "" {
		return 0, ErrEmptyValue
	}
	uid, err := s.repository.GetUserIdByUUID(token)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (s *AuthService) GetUserByUserID(id int) (*module.User, error) {
	user, err := s.repository.GetUserByID(id)
	if err != nil {
		log.Println("Error:service:auth:GetUserByUserID: FindIdUser")
		return nil, err
	}
	if id == 0 {
		log.Println("Error:service:auth:GetUserByUserID: id is 0")
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) DeleteSessionToken(token string) error {
	err := s.repository.Delete(token)
	if err != nil {
		log.Println("Error:service:auth:DeleteSessionToken: Delete")
		return err
	}
	return nil
}

func (s *AuthService) DeleteExpiredSessions() error {
	err := s.repository.DeleteExpiredSession()
	if err != nil {
		log.Println("Error:service:auth:DeleteExpiredSessions: DeleteExpiredSession")
		return err
	}
	return nil
}
