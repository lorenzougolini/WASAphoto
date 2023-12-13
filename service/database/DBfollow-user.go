package database

import (
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) FollowUser(user map[string]string, username string) error {

	sqlStmt := "SELECT * FROM follows WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, user["id"], username)
	if err != nil {
		return fmt.Errorf("Error adding the ban: %v", err)
	}
	if rowsAffected, err := res.RowsAffected(); rowsAffected != 0 || err != nil {
		return fmt.Errorf("The user '%s' is already followed", username)
	}

	sqlStmt = "INSERT INTO follows VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, user["id"], username)
	if err != nil {
		return fmt.Errorf("Error adding the follow: %v", err)
	}
	return nil
}

// unfollow user
func (db *appdbimpl) UnfollowUser(id string, username string) error {
	sqlStmt := "DELETE FROM follows WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, id, username)
	if err != nil {
		return fmt.Errorf("Error removing the follow: %v", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("The user '%s' is not followed at the moment", username)
	}
	return nil
}
