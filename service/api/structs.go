package api

import "strings"

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

// var Logged = User{}
var errUncorrectLogin string = "User is not correctly authenticated"

func formatId(id string) string {
	parts := strings.Split(id, "-")
	lastPart := parts[len(parts)-1]
	return lastPart[len(lastPart)-12:]
}
