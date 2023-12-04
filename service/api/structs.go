package api

import (
	"github.com/google/uuid"
)

type User struct {
	userID   uuid.UUID
	username string
	profile  Profile
}

type Profile struct {
	photos    []Photo
	followers []User
	following []User
	banned    []User
}

type Photo struct {
	photoID     uuid.UUID
	userID      uuid.UUID
	picture     string
	dateAndTime string
	likes       []Like
	comments    []Comment
}

type Like struct {
	likeID      uuid.UUID
	username    string
	photoID     uuid.UUID
	dateAndTime string
}

type Comment struct {
	commentID   uuid.UUID
	username    string
	photoID     uuid.UUID
	commentText string
	dateAndTime string
}

// func find(list interface{}, elem interface{}, attr interface{}) (int, bool) {
// 	for ind, el := range list {
// 		if el.attr == elem {
// 			return ind, 1
// 		} else {
// 			return ind, 0
// 		}
// 	}
// }
