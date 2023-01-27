package module

type Category struct {
	ID     int
	PostID int
	Tag  string
	Posts  []Post
}
