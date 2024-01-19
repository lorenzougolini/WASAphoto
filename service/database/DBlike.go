package database

import (
	"encoding/json"
	"fmt"
)

func (db *appdbimpl) GetLikeById(id string) (bool, Like, error) {

	var like Like
	err := db.c.QueryRow("SELECT * FROM likes WHERE likeid = ?", id).Scan(&like.LikeID, &like.PhotoID, &like.UserID, &like.DateAndTime)
	if err != nil {
		return false, like, fmt.Errorf("error retreiving the like")
	}
	return true, like, nil
}

func (db *appdbimpl) AddLike(like string) error {

	added_like := Like{}
	_ = json.Unmarshal([]byte(like), &added_like)

	stmt, _ := db.c.Prepare("INSERT INTO likes VALUES (?, ?, ?, ?)")
	_, err := stmt.Exec(
		added_like.LikeID,
		added_like.PhotoID,
		added_like.UserID,
		added_like.DateAndTime)
	if err != nil {
		return fmt.Errorf("error liking the photo")
	}
	return nil
}

func (db *appdbimpl) RemoveLike(likeid string, photoid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM likes WHERE likeid=? AND photoid=?")
	_, err := stmt.Exec(likeid, photoid)
	if err != nil {
		return fmt.Errorf("error removing the like from the photo")
	}
	return nil
}

// func (db *appdbimpl) GetPhotoLikes(photoid string) ([]string, error) {

// 	var like_list []string
// 	rows, err := db.c.Query("SELECT * FROM likes WHERE photoid = ?", photoid)
// 	if err != nil {
// 		return []string{}, fmt.Errorf("%v", err)
// 	}
// 	for rows.Next() {
// 		var likeid string
// 		err := rows.Scan(&likeid)
// 		if err != nil {
// 			return []string{}, fmt.Errorf("%v", err)
// 		}
// 		like_list = append(like_list, likeid)
// 	}
// 	return like_list, nil
// }
