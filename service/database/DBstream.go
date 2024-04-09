package database

import (
	"fmt"
)

func (db *appdbimpl) GetStream(userid string) (Stream, error) {
	stream := Stream{}

	// get users followed
	followedRows, err := db.c.Query("SELECT followedid FROM follows WHERE followerid = ?", userid)
	if err != nil {
		return stream, fmt.Errorf(errRetrievingStream, err)
	}
	defer followedRows.Close()

	for followedRows.Next() {
		var followedid string

		if err := followedRows.Scan(&followedid); err != nil {
			return stream, fmt.Errorf(errRetrievingStream, err)

		}

		// get photo data for each followed user
		photoRows, err := db.c.Query("SELECT photoid FROM photos WHERE authorid = ? AND authorid NOT IN ( SELECT bannerid FROM bans WHERE bannedid = ? ) ORDER BY dateAndTime DESC;", followedid, userid)
		if err != nil {
			return stream, fmt.Errorf(errRetrievingStream, err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photoid string

			if err := photoRows.Scan(&photoid); err != nil {
				return stream, fmt.Errorf(errRetrievingStream, err)
			}

			photoData, err := db.GetPhotoData(photoid)
			if err != nil {
				return stream, fmt.Errorf(errRetrievingStream, err)
			}
			stream.Posts = append(stream.Posts, photoData)
		}
		if photoRows.Err() != nil {
			return stream, photoRows.Err()
		}
	}

	if followedRows.Err() != nil {
		return stream, followedRows.Err()
	}
	return stream, nil
}
