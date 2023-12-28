package database

import (
	"fmt"
)

func (db *appdbimpl) GetStream(userid string) (Stream, error) {
	stream := Stream{}

	// get users followed
	followedRows, err := db.c.Query("SELECT followedid FROM follows WHERE userid = ?", userid)
	if err != nil {
		return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
	}
	defer followedRows.Close()

	for followedRows.Next() {
		var followedid string

		if err := followedRows.Scan(&followedid); err != nil {
			return stream, fmt.Errorf("error retrieving the stream. err: %v", err)

		}

		// get photos for each followed user
		photoRows, err := db.c.Query("SELECT photoid, picture, dateAndTime, description FROM photos WHERE userid = ?", followedid)
		if err != nil {
			return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photoid string
			photo := struct {
				File   string
				Author string
				Likes  struct {
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
			}{}

			if err := photoRows.Scan(&photoid, &photo.File, &photo.DateAndTime, &photo.Description); err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
			}

			// retrieve likes and comments of each photo
			likeRows, err := db.c.Query("SELECT userid FROM likes WHERE photoid = ?", photoid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
			}
			defer likeRows.Close()

			for likeRows.Next() {
				var likeAuthor string

				if err := likeRows.Scan(&likeAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
				}

				if user, err := db.GetById(likeAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the stream. err: %v", err)

				} else {
					photo.Likes.Usernames = append(photo.Likes.Usernames, user.Username)
				}
			}
			photo.Likes.NumberOfLikes = len(photo.Likes.Usernames)

			commentRows, err := db.c.Query("SELECT userid, commentText, dateAndTime FROM comments WHERE photoid = ?", photoid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
			}
			defer commentRows.Close()

			for commentRows.Next() {
				var commentAuthor string
				comment := struct {
					Username    string
					CommentText string
					DateAndTime string
				}{}

				if err := commentRows.Scan(&commentAuthor, &comment.CommentText, &comment.DateAndTime); err != nil {
					return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
				}

				if user, err := db.GetById(commentAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the stream. err: %v", err)

				} else {
					comment.Username = user.Username
				}
				photo.Comments.Comment = append(photo.Comments.Comment, comment)
			}
			photo.Comments.NumberOfComments = len(photo.Comments.Comment)

			photoAuthor, err := db.GetById(followedid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %v", err)
			}
			photo.Author = photoAuthor.Username
			stream.Photos = append(stream.Photos, photo)
		}
	}
	return stream, nil
}
