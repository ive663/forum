package module

import "time"

type Comment struct {
	ID       int
	AuthorID int
  Author   string
  Likes    int
  Dislikes int
	PostID   int
	Message  string
	Date     time.Time
  DateFormat string
}

func (c *Comment) SetDateFormat() {
  c.DateFormat = c.Date.Format("02.01.2006 15:04")
}

type CommentList []Comment

func (c CommentList) PrepToView() CommentList {
  for i := range c {
    c[i].SetDateFormat()
  }
  return c
}
