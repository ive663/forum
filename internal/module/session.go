package module

import "time"

type Session struct {
	ID       int
	UserID   int
	UUID     string
	CreatedAt time.Time
	ExpiresAt  time.Time
}
