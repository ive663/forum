package module

import "time"

type Post struct {
	ID         int
	Title      string
	AuthorID   int
  Author     string
	Message    string
  Likes      int
  Dislikes   int
	Liked      bool
  CategoryID int
	Category   string
	Categories []Category
	Comments   []Comment
	Date       time.Time
  DateFormat string
}

func (p *Post) SetDateFormat() {
  p.DateFormat = p.Date.Format("02.01.2006 15:04")
}

type PostList []Post 

func (p PostList) PrepToView() PostList {
  for i := range p {
    p[i].SetDateFormat()
  }
  return p
}

