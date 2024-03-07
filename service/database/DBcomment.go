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

func (db *appdbimpl) AddComment(commentObj string) error {

	add_comment := Comment{}
	_ = json.Unmarshal([]byte(commentObj), &add_comment)

	stmt, _ := db.c.Prepare("INSERT INTO comments VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(
		add_comment.CommentID,
		add_comment.PhotoID,
		add_comment.UserID,
		add_comment.CommentText,
		add_comment.DateAndTime)

	if err != nil {
		return fmt.Errorf("error liking the photo")
	}
	return nil
}

func (db *appdbimpl) RemoveComment(commentid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM comments WHERE commentid=?")
	_, err := stmt.Exec(commentid)
	if err != nil {
		return fmt.Errorf("error removing the comment from the photo")
	}
	return nil
}
