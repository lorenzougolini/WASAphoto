package database

type User struct {
	UserID   string
	Username string
}

type Photo struct {
	PhotoID     string
	AuthorID    string
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
	User        string
	PhotoID     string
	CommentText string
	DateAndTime string
}

type PhotoComment struct {
	Username    string
	CommentText string
	DateAndTime string
}

type PhotoData struct {
	PhotoID     string
	Author      string
	Description string
	DateAndTime string
	Likes       []string
	Comments    []PhotoComment
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
	Bans struct {
		Usernames    []string
		NumberOfBans int
	}
}

type Stream struct {
	Posts []PhotoData
}

var errRetrievingProfile = "error retrieving the profile: %w"
var errRetrievingPhotoData = "error retrieving the photo data: %w"
