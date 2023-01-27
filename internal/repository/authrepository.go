package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/ive663/forum/internal/module"
)

type Auth interface {
	CreateNewSession(*module.Session) error
	GetUserIdByUUID(UUID string) (int, error)
	Delete(uuid string) error
	CreateNewUser(*module.User) error
	FindByLogin(login string) (*module.User, error)
	GetUserByID(id int) (*module.User, error)
	DeleteExpiredSession() error
	UpdateSession(s *module.Session) error
	IsSessionExists(userID int) (bool, error)
}

type AuthRepository struct {
	db *sql.DB
}

func newAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) IsSessionExists(userID int) (bool, error) {
	user_id := "0"
	err := r.db.QueryRow("SELECT uuid FROM sessions WHERE user_id = $1", userID).Scan(&user_id)
	if err != nil {
		log.Println("repo:IsSessionExist: false")
		return false, nil
	}
	log.Println("repo:IsSessionExist: true")
	return user_id != "0", nil
}

func (r *AuthRepository) CreateNewSession(s *module.Session) error {
	query := "INSERT INTO sessions (user_id, uuid, created_at, expires_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, s.UserID, s.UUID, s.CreatedAt, s.ExpiresAt)
	log.Println("repo:auth: session ID creatted")
	if err != nil {
		log.Println("err:repo:auth: CreateNewSession")
		return err
	}
	return nil
}

func (r *AuthRepository) UpdateSession(s *module.Session) error {
	query := "UPDATE sessions SET uuid = ?, created_at = ?, expires_at = ? WHERE user_id = ?"
	_, err := r.db.Exec(query, s.UUID, s.CreatedAt, s.ExpiresAt, s.UserID)
	log.Println("authrepo: updateSession")
	if err != nil {
		log.Println("error:authRepo: UpdateSession")
	}
	return nil
}

func (r *AuthRepository) GetUserIdByUUID(UUID string) (int, error) {
	s := &module.Session{}
	if UUID == "" {
		return 0, errors.New("Nil in UUID")
	}
	err := r.db.QueryRow("SELECT user_id FROM sessions WHERE uuid = ?", UUID).Scan(
		&s.UserID)
	if err == sql.ErrNoRows {
		return 0, errors.New("Record not found")
	}
	if err != nil {
		return 0, err
	}
	return s.UserID, nil
}

func (r *AuthRepository) Delete(uuid string) error {
	if _, err := r.db.Exec("DELETE FROM sessions WHERE uuid = ?", uuid); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (r *AuthRepository) CreateNewUser(u *module.User) error {
	log.Println(u.Login, u.EncryptedPassword, u.Email)
	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
	if _, err := r.db.Exec(query, u.Login, u.EncryptedPassword, u.Email); err != nil {
		log.Printf("error:authRepo:CreatingNewUser %v\n", err)
		return err
	}
	return nil
}

func (r *AuthRepository) FindByLogin(login string) (*module.User, error) {
	if login == "" {
		return nil, errors.New("error:authRepo:findByLogin empty str")
	}
	u := &module.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password FROM users WHERE username = ?",
		login,
	).Scan(&u.ID, &u.Login, &u.EncryptedPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("error:authRepo:findByLogin: Record not found")
	}
	if err != nil {
		return nil, errors.New("error:authRepo:findByLogin: Record not found")
	}
	return u, nil
}

func (r *AuthRepository) GetUserByID(id int) (*module.User, error) {
	u := &module.User{}
	err := r.db.QueryRow("SELECT id, username FROM users WHERE id = ?", id).Scan(&u.ID, &u.Login)
	if err == sql.ErrNoRows {
		log.Println("error:authRepo:GetUserByID: Record not found")
		return nil, err
	}
	if err != nil {
		log.Println("error:authRepo:GetUserByID: DB error")
		return nil, err
	}
	return u, nil
}

func (r *AuthRepository) DeleteExpiredSession() error {
	if _, err := r.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now()); err != nil {
		log.Print("error:authRepo:DeleteExpiredSession")
		return err
	}
	return nil
}
