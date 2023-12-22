package database

import (
	"encoding/json"
	"fmt"
)

func (db *appdbimpl) FollowUser(loggedUser string, followedUser string) error {

	logged_user := User{}
	followed_user := User{}
	json.Unmarshal([]byte(loggedUser), &logged_user)
	json.Unmarshal([]byte(followedUser), &followed_user)

	// check followed user didn't ban the logged user
	if banned, err := db.IsBanned(followed_user.UserID, logged_user.UserID); banned || err != nil {
		return fmt.Errorf("impossible to follow the user; %v", err)
	}

	// check user isn't already followed
	followed, err := db.IsFollowed(logged_user.UserID, followed_user.UserID)
	if followed {
		return fmt.Errorf("the user '%s' is already followed", followed_user.Username)

	} else if err != nil {
		return fmt.Errorf("error adding the follow: %v", err)
	}

	stmt, _ := db.c.Prepare("INSERT INTO follows (userid, followedid) VALUES (?, ?);")
	res, err := stmt.Exec(logged_user.UserID, followed_user.UserID)
	fmt.Println(res)
	if err != nil {
		return fmt.Errorf("error adding the follow: %v", err)
	}
	return nil
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
		sqlStmt := "DELETE FROM follows WHERE userid=? AND followedid=?"
		_, err := db.c.Exec(sqlStmt, loggedId, followed_user.UserID)
		return err
	}
}

func (db *appdbimpl) IsFollowed(id string, followedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE userid=? AND followedid=?", id, followedId).Scan(&count)
	if err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
