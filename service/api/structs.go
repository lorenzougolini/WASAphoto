package api

import (
	"net/http"
)

type User struct {
	UserID   string
	Username string
}

type Profile struct {
	Posts struct {
		Photos []struct {
			File        string
			Description string
			DateAndTime string
		}
		NumberOfPosts int
	}
	Followers struct {
		Usernames         []string
		NumberOfFollowers int
	}
	Following struct {
		Usernames         []string
		NumberOfFollowing int
	}
}

type Photo struct {
	PhotoID     string
	UserID      string
	Picture     string
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

type Stream struct {
	Photos []struct {
		File   string
		Author string
		Likes  []struct {
			NumberOfLikes int
			Usernames     []string
		}
		Comments struct {
			NumberOfComments int
			Comment          []struct {
				Username    string
				CommentText string
				DateAndTime string
			}
		}
		Description string
		DateAndTime string
	}
}

var Logged = User{}
var uncorrectLogin string = "User is not correctly authenticated"

// var UsernameToId = make(map[string]string)
// var Users = make(map[string]User)
// var Photos = make(map[string]Photo)
// var Likes = make(map[string]Like)
// var Comments = make(map[string]Comment)

func checkLogin(r *http.Request) bool {
	// return id == Logged.UserID && username == Logged.Username
	authToken := r.Header.Get("Authorization")
	if authToken == "" || authToken != Logged.UserID {
		return false
	}
	return true
}
