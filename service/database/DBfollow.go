package database

import (
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) FollowUser(user User, username string) error {

	// get the followed user
	followedUser, err := db.GetByUsername(username)
	if err != nil {
		return err
	}

	// check followed user didn't ban the following user
	count, err := db.IsBanned(user.UserID, followedUser.UserID)
	if err != nil {
		return err
	} else if count {
		return fmt.Errorf("impossible to follow the user")
	}

	sqlStmt := "SELECT * FROM follows WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, user.UserID, followedUser.UserID)
	if err != nil {
		return fmt.Errorf("error adding the ban: %v", err)
	}
	if rowsAffected, err := res.RowsAffected(); rowsAffected != 0 || err != nil {
		return fmt.Errorf("the user '%s' is already followed", username)
	}

	sqlStmt = "INSERT INTO follows VALUES (?, ?);"
	_, err = db.c.Exec(sqlStmt, user.UserID, followedUser.UserID)
	if err != nil {
		return fmt.Errorf("error adding the follow: %v", err)
	}
	return nil
}

// unfollow user
func (db *appdbimpl) UnfollowUser(id string, username string) error {

	// get the followed user
	unfollowedUser, err := db.GetByUsername(username)
	if err != nil {
		return err
	}

	sqlStmt := "DELETE FROM follows WHERE userid=? AND ban=?"
	res, err := db.c.Exec(sqlStmt, id, unfollowedUser.UserID)
	if err != nil {
		return fmt.Errorf("error removing the follow: %v", err)
	}
	if affectedRows, _ := res.RowsAffected(); affectedRows == 0 {
		return fmt.Errorf("the user '%s' is not followed at the moment", username)
	}
	return nil
}
