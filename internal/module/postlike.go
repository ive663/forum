package module

import "time"

type PostLike struct {
	PostID int
	UserID int
	Date   time.Time
}
