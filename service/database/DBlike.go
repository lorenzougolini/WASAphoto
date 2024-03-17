package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

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

func (db *appdbimpl) GetLikeByLikeId(id string) (bool, Like, error) {

	var like Like
	err := db.c.QueryRow("SELECT * FROM likes WHERE likeid = ?", id).Scan(&like.LikeID, &like.PhotoID, &like.UserID, &like.DateAndTime)
	if errors.Is(err, sql.ErrNoRows) {
		return false, like, fmt.Errorf("error retreiving the like")
	}
	return true, like, nil
}

func (db *appdbimpl) GetLikeByUserId(userid string, photoid string) (bool, error) {

	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM likes WHERE photoid = ? AND userid = ?", photoid, userid).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return count > 0, fmt.Errorf("error retreiving the like")
	}
	return count > 0, nil
}

func (db *appdbimpl) GetLikesNamesByPhotoId(photoid string) ([]PhotoLike, error) {

	rows, err := db.c.Query("SELECT u.username, l.likeid FROM likes l JOIN users u ON l.userid = u.userid WHERE l.photoid = ? ORDER BY l.dateAndTime DESC", photoid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []PhotoLike
	for rows.Next() {
		var username string
		var likeid string

		if err := rows.Scan(&username, &likeid); err != nil {
			return nil, err
		}
		likes = append(likes, PhotoLike{LikeID: likeid, Username: username})
	}

	if rows.Err() != nil {
		return likes, rows.Err()
	}
	return likes, nil
}

func (db *appdbimpl) RemoveLike(likeid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM likes WHERE likeid=?")
	_, err := stmt.Exec(likeid)
	if err != nil {
		return fmt.Errorf("error removing the like from the photo")
	}
	return nil
}
