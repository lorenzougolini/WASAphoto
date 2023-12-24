package database

import (
	"encoding/json"
	"fmt"
)

func (db *appdbimpl) GetPhotoById(photoid string) (bool, Photo, error) {

	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoid = ?", photoid).Scan(&photo.PhotoID, &photo.UserID, &photo.Picture, &photo.DateAndTime, &photo.Description)
	if err != nil {
		return false, photo, fmt.Errorf("error retreiving the photo")
	}
	return true, photo, nil
}

func (db *appdbimpl) PostPhoto(photo string) error {

	posted_photo := Photo{}
	json.Unmarshal([]byte(photo), &posted_photo)
	// fmt.Printf("%v", posted_photo)

	stmt, _ := db.c.Prepare("INSERT INTO photos VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(posted_photo.PhotoID, posted_photo.UserID, posted_photo.Picture, posted_photo.DateAndTime, posted_photo.Description)
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
