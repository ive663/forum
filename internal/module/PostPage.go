package module

// // Added likes and dislikes
type PostPage struct {
	Post             *Post
	Comments         []Comment
	PostLikes        int
	PostDislikes     int
	CommentsLikes    map[int][]int
	CommentsDislikes map[int][]int
	Authorization    bool
}

