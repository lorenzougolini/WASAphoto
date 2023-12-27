package database

import (
	"fmt"
)

func (db *appdbimpl) GetStream(userid string) (Stream, error) {
	stream := Stream{}

	// get users followed
	followedRows, err := db.c.Query("SELECT followedid FROM follows WHERE userid = ?", userid)
	if err != nil {
		return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
	}
	defer followedRows.Close()

	for followedRows.Next() {
		var followedid string

		if err := followedRows.Scan(&followedid); err != nil {
			return stream, fmt.Errorf("error retrieving the profile. err: %v", err)

		}

		// get photos for each followed user
		photoRows, err := db.c.Query("SELECT photoid, picture, dateAndTime, description FROM photos WHERE userid = ?", followedid)
		if err != nil {
			return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photoid string
			photo := struct {
				file   string
				author string
				likes  struct {
					numberOfLikes int
					usernames     []string
				}
				comments struct {
					numberOfComments int
					comment          []struct {
						username    string
						commentText string
						dateAndTime string
					}
				}
				description string
				dateAndTime string
			}{}

			if err := photoRows.Scan(&photoid, &photo.file, &photo.dateAndTime, &photo.description); err != nil {
				return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
			}

			// retrieve likes and comments of each photo
			likeRows, err := db.c.Query("SELECT userid FROM likes WHERE photoid = ?", photoid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
			}
			defer likeRows.Close()

			for likeRows.Next() {
				var likeAuthor string

				if err := likeRows.Scan(&likeAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
				}

				if user, err := db.GetById(likeAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the profile. err: %v", err)

				} else {
					photo.likes.usernames = append(photo.likes.usernames, user.Username)
				}
			}
			photo.likes.numberOfLikes = len(photo.likes.usernames)

			commentRows, err := db.c.Query("SELECT userid, commentText, dateAndTime FROM comments WHERE photoid = ?", photoid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
			}
			defer commentRows.Close()

			for commentRows.Next() {
				var commentAuthor string
				comment := struct {
					username    string
					commentText string
					dateAndTime string
				}{}

				if err := commentRows.Scan(&commentAuthor, &comment.commentText, &comment.dateAndTime); err != nil {
					return stream, fmt.Errorf("error retrieving the profile. err: %v", err)
				}

				if user, err := db.GetById(commentAuthor); err != nil {
					return stream, fmt.Errorf("error retrieving the profile. err: %v", err)

				} else {
					comment.username = user.Username
				}
				photo.comments.comment = append(photo.comments.comment, comment)
			}
			photo.comments.numberOfComments = len(photo.comments.comment)

			stream.Photos = append(stream.Photos, photo)
		}
	}
	return stream, nil
}
