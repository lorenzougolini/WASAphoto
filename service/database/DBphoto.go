package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetPhotoById(photoid string) (bool, Photo, error) {

	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoid = ?", photoid).Scan(&photo.PhotoID, &photo.AuthorID, &photo.PicFile, &photo.DateAndTime, &photo.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return false, photo, fmt.Errorf("error retreiving the photo")
	}
	return true, photo, nil
}

func (db *appdbimpl) PostPhoto(photo string) error {

	posted_photo := Photo{}
	_ = json.Unmarshal([]byte(photo), &posted_photo)

	stmt, _ := db.c.Prepare("INSERT INTO photos VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(posted_photo.PhotoID, posted_photo.AuthorID, posted_photo.PicFile, posted_photo.DateAndTime, posted_photo.Description)
	if err != nil {
		return fmt.Errorf("error uploading the photo")
	}
	return nil
}

func (db *appdbimpl) RemovePhoto(photoid string) error {

	// no need to delete all likes and comments associated, the db has delete on cascade

	stmt, _ := db.c.Prepare("DELETE FROM photos WHERE photoid=?")
	_, err := stmt.Exec(photoid)
	if err != nil {
		return fmt.Errorf("error deleting the photo")
	}
	return nil
}

func (db *appdbimpl) GetPhotoData(photoid string) (PhotoData, error) {

	var photodata PhotoData
	photodata.PhotoID = photoid

	err := db.c.QueryRow("SELECT picFile, authorid, description, dateAndTime FROM photos WHERE photoid = ?", photoid).Scan(&photodata.File, &photodata.Author, &photodata.Description, &photodata.DateAndTime)
	if errors.Is(err, sql.ErrNoRows) {
		return photodata, err
	}

	// get author
	if author, err := db.GetById(photodata.Author); err != nil {
		return photodata, err
	} else {
		photodata.Author = author.Username
	}

	// get likes
	photodata.Likes, err = db.GetLikesNamesByPhotoId(photoid)
	if err != nil {
		return photodata, err
	}

	// get comments
	photodata.Comments, err = db.GetCommentsByPhotoId(photoid)
	if err != nil {
		return photodata, err
	}

	return photodata, nil
}
