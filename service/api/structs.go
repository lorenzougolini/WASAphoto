package api

import (
	"WASAphoto.uniroma1.it/WASAphoto/service/database"
)

type User struct {
	UserID   string
	Username string
}

func (user *User) userToDBUser() database.User {
	return database.User{
		UserID:   user.UserID,
		Username: user.Username,
	}
}

func userFromDBUser(user database.User) User {
	return User{
		UserID:   user.UserID,
		Username: user.Username,
	}
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
