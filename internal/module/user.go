package module

type User struct {
	ID                int
	Login             string
	Password          string
	EncryptedPassword string
	Email             string
	Posts             []Post
	Comments          []Comment
	Authorization     bool
}
