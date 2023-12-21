package database

import (
	"encoding/json"
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) FollowUser(loggedUser string, username string) error {

	logged_user := User{}
	json.Unmarshal([]byte(loggedUser), &logged_user)

	// get the followed user object
	followedUser, err := db.GetByUsername(username)
	if err != nil {
		return err
	}

	// check followed user didn't ban the logged user
	if banned, err := db.IsBanned(followedUser.UserID, logged_user.UserID); banned || err != nil {
		return fmt.Errorf("impossible to follow the user")
	}

	// check user isn't already followed
	followed, err := db.IsFollowed(logged_user.UserID, followedUser.UserID)
	if followed {
		return fmt.Errorf("the loggedUser '%s' is already followed", username)

	} else if err != nil {
		return fmt.Errorf("error adding the ban: %v", err)

	} else {
		sqlStmt := "INSERT INTO follows VALUES (?, ?);"
		_, err = db.c.Exec(sqlStmt, logged_user.UserID, followedUser.UserID)
		if err != nil {
			return fmt.Errorf("error adding the follow: %v", err)
		}
		return nil
	}
}

// unfollow loggedUser
func (db *appdbimpl) UnfollowUser(loggedId string, followedUsername string) error {

	followed_user, err := db.GetByUsername(followedUsername)
	if err != nil {
		return fmt.Errorf("error removing the follow: %v", err)
	}

	followed, err := db.IsFollowed(loggedId, followed_user.UserID)
	if !followed {
		return fmt.Errorf("the user '%s' is not followed at the moment", followed_user.Username)

	} else if err != nil {
		return fmt.Errorf("error removing the follow: %v", err)

	} else {
		sqlStmt := "DELETE FROM follows WHERE userid=? AND follow=?"
		_, err := db.c.Exec(sqlStmt, loggedId, followed_user.UserID)
		return err
	}
}

func (db *appdbimpl) IsFollowed(id string, followedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE userid=? AND ban=?", id, followedId).Scan(&count)
	if err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
