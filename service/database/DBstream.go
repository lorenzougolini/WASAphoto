package database

import (
	"fmt"
)

func (db *appdbimpl) GetStream(userid string) (Stream, error) {
	stream := Stream{}

	// get users followed
	followedRows, err := db.c.Query("SELECT followedid FROM follows WHERE userid = ?", userid)
	if err != nil {
		return stream, fmt.Errorf("error retrieving the stream. err: %w", err)
	}
	defer followedRows.Close()

	for followedRows.Next() {
		var followedid string

		if err := followedRows.Scan(&followedid); err != nil {
			return stream, fmt.Errorf("error retrieving the stream. err: %w", err)

		}

		// get photo data for each followed user
		photoRows, err := db.c.Query("SELECT photoid FROM photos WHERE userid = ? ORDER BY dateAndTime DESC", followedid)
		if err != nil {
			return stream, fmt.Errorf("error retrieving the stream. err: %w", err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photoid string

			if err := photoRows.Scan(&photoid); err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %w", err)
			}

			photoData, err := db.GetPhotoData(photoid)
			if err != nil {
				return stream, fmt.Errorf("error retrieving the stream. err: %w", err)
			}
			stream.Posts = append(stream.Posts, photoData)
		}
	}
	return stream, nil
}
