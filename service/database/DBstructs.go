package database

type User struct {
	UserID   string
	Username string
}

type Photo struct {
	PhotoID     string
	UserID      string
	PicPath     string
	DateAndTime string
	Description string
}

type Like struct {
	LikeID      string
	PhotoID     string
	UserID      string
	DateAndTime string
}

type Comment struct {
	CommentID   string
	UserID      string
	PhotoID     string
	CommentText string
	DateAndTime string
}

type PhotoData struct {
	PhotoID     string
	Author      string
	Description string
	DateAndTime string
	Likes       []string
	Comments    []struct {
		Username    string
		CommentText string
		DateAndTime string
	}
}

type Profile struct {
	Posts     []PhotoData
	Followers struct {
		Usernames         []string
		NumberOfFollowers int
	}
	Following struct {
		Usernames         []string
		NumberOfFollowing int
	}
}

type Stream struct {
	Posts []PhotoData
}

var Logged = User{}
