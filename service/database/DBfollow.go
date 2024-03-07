package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) FollowUser(reqUserId string, followUserId string) error {

	// check followed user didn't ban the logged user
	if banned, err := db.IsBanned(followUserId, reqUserId); banned || err != nil {
		return fmt.Errorf("impossible to follow the user; %w", err)
	}

	// check user isn't already followed
	followed, err := db.IsFollowed(reqUserId, followUserId)
	if followed {
		return fmt.Errorf("the user is already followed")
	}
	if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}

	stmt, _ := db.c.Prepare("INSERT INTO follows (userid, followedid) VALUES (?, ?);")
	_, err = stmt.Exec(reqUserId, followUserId)
	if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}
	return nil
}

// unfollow loggedUser
func (db *appdbimpl) UnfollowUser(reqUserId string, unfollowUserId string) error {

	followed, err := db.IsFollowed(reqUserId, unfollowUserId)
	if !followed {
		return fmt.Errorf("the user is not followed at the moment")
	}
	if err != nil {
		return fmt.Errorf("error removing the follow: %w", err)
	}

	sqlStmt := "DELETE FROM follows WHERE userid=? AND followedid=?"
	_, err = db.c.Exec(sqlStmt, reqUserId, unfollowUserId)
	return err
}

func (db *appdbimpl) IsFollowed(reqUserId string, followedId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE userid=? AND followedid=?", reqUserId, followedId).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return count > 0, nil
}
