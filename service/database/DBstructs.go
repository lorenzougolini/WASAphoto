package database

type User struct {
	UserID   string
	Username string
}

type Photo struct {
	PhotoID     string
	AuthorID    string
	PicFile     []byte
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
	User        string
	PhotoID     string
	CommentText string
	DateAndTime string
}

type PhotoLike struct {
	LikeID   string
	Username string
}

type PhotoComment struct {
	CommentID   string
	Username    string
	CommentText string
	DateAndTime string
}

type PhotoData struct {
	PhotoID     string
	File        []byte
	Author      string
	Description string
	DateAndTime string
	Likes       []PhotoLike
	Comments    []PhotoComment
}

type Profile struct {
	Username  string
	Posts     []PhotoData
	Followers struct {
		Usernames         []string
		NumberOfFollowers int
	}
	Following struct {
		Usernames         []string
		NumberOfFollowing int
	}
	Banned bool
}

type Stream struct {
	Posts []PhotoData
}

var errRetrievingProfile = "error retrieving the profile: %w"

// var errRetrievingPhotoData = "error retrieving the photo data: %w"
var errRetrievingStream = "error retrieving the stream: %w"
