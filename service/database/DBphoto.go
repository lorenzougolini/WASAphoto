package database

import (
	"fmt"
)

func (db *appdbimpl) GetPhotoById(id string) (bool, Photo, error) {

	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoid = ?").Scan(&photo)
	if err != nil {
		return false, photo, fmt.Errorf("error retreiving the photo")
	}
	return true, photo, nil
}

func (db *appdbimpl) PostPhoto(photo Photo) error {

	stmt, _ := db.c.Prepare("INSERT INTO photos VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(photo.PhotoID, photo.UserID, photo.Picture, photo.DateAndTime, photo.Description)
	if err != nil {
		return fmt.Errorf("error uploading the photo")
	}
	return nil
}

func (db *appdbimpl) RemovePhoto(photoid string, userid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM photos WHERE photoid=? AND userid=?")
	_, err := stmt.Exec(photoid, userid)
	if err != nil {
		return fmt.Errorf("error deleting the photo")
	}
	// remove likes and comments by photo
	return nil
}
