package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetPhotoById(photoid string) (bool, Photo, error) {

	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoid = ?", photoid).Scan(&photo.PhotoID, &photo.UserID, &photo.PicPath, &photo.DateAndTime, &photo.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return false, photo, fmt.Errorf("error retreiving the photo")
	}
	return true, photo, nil
}

func (db *appdbimpl) PostPhoto(photo string) error {

	posted_photo := Photo{}
	_ = json.Unmarshal([]byte(photo), &posted_photo)
	// fmt.Printf("%v", posted_photo)

	stmt, _ := db.c.Prepare("INSERT INTO photos VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(posted_photo.PhotoID, posted_photo.UserID, posted_photo.PicPath, posted_photo.DateAndTime, posted_photo.Description)
	if err != nil {
		return fmt.Errorf("error uploading the photo")
	}
	return nil
}

func (db *appdbimpl) RemovePhoto(userid string, photoid string) error {

	// no need to delete all likes and comments associated, the db has delete on cascade

	stmt, _ := db.c.Prepare("DELETE FROM photos WHERE photoid=? AND userid=?")
	_, err := stmt.Exec(photoid, userid)
	if err != nil {
		return fmt.Errorf("error deleting the photo")
	}
	return nil
}

func (db *appdbimpl) GetPhotoData(photoid string) (PhotoData, error) {

	var photodata PhotoData
	err := db.c.QueryRow("SELECT userid, description, dateAndTime FROM photos WHERE photoid = ?", photoid).Scan(&photodata.Author, &photodata.Description, &photodata.DateAndTime)
	if errors.Is(err, sql.ErrNoRows) {
		return photodata, fmt.Errorf("error retreiving the photo data")
	}

	// get author
	if author, err := db.GetById(photodata.Author); err != nil {
		return photodata, fmt.Errorf("error retreiving the photo data")
	} else {
		photodata.Author = author.Username
	}

	// get likes
	if likes, err := db.c.Query("SELECT userid FROM likes WHERE photoid = ?", photoid); err != nil {
		return photodata, fmt.Errorf("error retreiving the photo data")
	} else {
		defer likes.Close()
		for likes.Next() {
			var idWhoLikes string
			if err := likes.Scan(&idWhoLikes); err != nil {
				return photodata, fmt.Errorf("error retreiving the photo data")
			}
			if userWhoLikes, err := db.GetById(idWhoLikes); err != nil {
				return photodata, fmt.Errorf("error retreiving the photo data")
			} else {
				photodata.Likes = append(photodata.Likes, userWhoLikes.Username)
			}
		}
	}

	// get comments
	if comments, err := db.c.Query("SELECT userid, commentText, dateAndTime FROM comments WHERE photoid = ?", photoid); err != nil {
		return photodata, fmt.Errorf("error retreiving the photo data")
	} else {
		defer comments.Close()
		for comments.Next() {
			var comment struct {
				Username    string
				CommentText string
				DateAndTime string
			}

			if err := comments.Scan(&comment.Username, &comment.CommentText, &comment.DateAndTime); err != nil {
				return photodata, fmt.Errorf("error retreiving the photo data")
			} else {

				if userWhoComments, err := db.GetById(comment.Username); err != nil {
					return photodata, fmt.Errorf("error retreiving the photo data")
				} else {
					comment.Username = userWhoComments.Username
				}
			}
			photodata.Comments = append(photodata.Comments, comment)
		}
	}
	photodata.PhotoID = photoid

	return photodata, nil

}
