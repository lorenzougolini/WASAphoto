package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddComment(commentObj string) error {

	add_comment := Comment{}
	_ = json.Unmarshal([]byte(commentObj), &add_comment)

	stmt, _ := db.c.Prepare("INSERT INTO comments VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(
		add_comment.CommentID,
		add_comment.PhotoID,
		add_comment.User,
		add_comment.CommentText,
		add_comment.DateAndTime)

	if err != nil {
		return fmt.Errorf("error liking the photo")
	}
	return nil
}

func (db *appdbimpl) GetCommentByCommentId(id string) (bool, Comment, error) {

	var comment Comment
	err := db.c.QueryRow("SELECT * FROM comments WHERE commentid = ?", id).Scan(&comment.CommentID, &comment.PhotoID, &comment.User, &comment.CommentText, &comment.DateAndTime)
	if errors.Is(err, sql.ErrNoRows) {
		return false, comment, fmt.Errorf("error retreiving the comment")
	}
	return true, comment, nil
}

func (db *appdbimpl) GetCommentsByPhotoId(photoid string) ([]PhotoComment, error) {

	rows, err := db.c.Query("SELECT u.username, c.commentid, c.commentText, c.dateAndTime FROM comments c JOIN users u ON c.userid = u.userid WHERE c.photoid = ?", photoid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []PhotoComment
	for rows.Next() {
		var comment PhotoComment
		if err := rows.Scan(&comment.Username, &comment.CommentID, &comment.CommentText, &comment.DateAndTime); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return comments, rows.Err()
	}
	return comments, nil
}

func (db *appdbimpl) RemoveComment(commentid string) error {

	stmt, _ := db.c.Prepare("DELETE FROM comments WHERE commentid=?")
	_, err := stmt.Exec(commentid)
	if err != nil {
		return fmt.Errorf("error removing the comment from the photo")
	}
	return nil
}
