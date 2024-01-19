package database

import (
	"encoding/json"
	"fmt"
)

func (db *appdbimpl) FollowUser(loggedUser string, followedUser string) error {

	logged_user := User{}
	followed_user := User{}
	_ = json.Unmarshal([]byte(loggedUser), &logged_user)
	_ = json.Unmarshal([]byte(followedUser), &followed_user)

	// check followed user didn't ban the logged user
	if banned, err := db.IsBanned(followed_user.UserID, logged_user.UserID); banned || err != nil {
		return fmt.Errorf("impossible to follow the user; %w", err)
	}

	// check user isn't already followed
	followed, err := db.IsFollowed(logged_user.UserID, followed_user.UserID)
	if followed {
		return fmt.Errorf("the user '%s' is already followed", followed_user.Username)

	} else if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}

	stmt, _ := db.c.Prepare("INSERT INTO follows (userid, followedid) VALUES (?, ?);")
	_, err = stmt.Exec(logged_user.UserID, followed_user.UserID)
	// fmt.Println(res)
	if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}
	return nil
}

// unfollow loggedUser
func (db *appdbimpl) UnfollowUser(loggedId string, unfollowedUsername string) error {

	unfollowed_user, err := db.GetByUsername(unfollowedUsername)
	if err != nil {
		return fmt.Errorf("error removing the follow: %w", err)
	}

	followed, err := db.IsFollowed(loggedId, unfollowed_user.UserID)
	if !followed {
		return fmt.Errorf("the user '%s' is not followed at the moment", unfollowed_user.Username)

	} else if err != nil {
		return fmt.Errorf("error removing the follow: %w", err)

	}
	sqlStmt := "DELETE FROM follows WHERE userid=? AND followedid=?"
	_, err = db.c.Exec(sqlStmt, loggedId, unfollowed_user.UserID)
	return err
}

func (db *appdbimpl) IsFollowed(id string, followedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE userid=? AND followedid=?", id, followedId).Scan(&count)
	if err != nil {
		return count > 0, err
	}
	return count > 0, nil
}
