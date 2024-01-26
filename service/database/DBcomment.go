package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetCommentByCommentId(id string) (bool, Comment, error) {

	var comment Comment
	err := db.c.QueryRow("SELECT * FROM comments WHERE commentid = ?", id).Scan(&comment.CommentID, &comment.PhotoID, &comment.UserID, &comment.CommentText, &comment.DateAndTime)
	if errors.Is(err, sql.ErrNoRows) {
		return false, comment, fmt.Errorf("error retreiving the comment")
	}
	return true, comment, nil
}

func (db *appdbimpl) AddComment(comment string) error {

	added_comment := Comment{}
	_ = json.Unmarshal([]byte(comment), &added_comment)

	stmt, _ := db.c.Prepare("INSERT INTO comments VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(
		added_comment.CommentID,
		added_comment.PhotoID,
		added_comment.UserID,
		added_comment.CommentText,
		added_comment.DateAndTime)

	if err != nil {
		return fmt.Errorf("error liking the photo")
	}
	return nil
}

func (db *appdbimpl) RemoveComment(commentid string, photoid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM comments WHERE commentid=? AND photoid=?")
	_, err := stmt.Exec(commentid, photoid)
	if err != nil {
		return fmt.Errorf("error removing the comment from the photo")
	}
	return nil
}
