package api

import (
	"time"

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
	dateAndTime time.Time
}

type Comment struct {
	commentID   uuid.UUID
	username    string
	photoID     uuid.UUID
	commentText string
	dateAndTime time.Time
}
