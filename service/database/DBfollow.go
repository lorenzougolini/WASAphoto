package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) FollowUser(followerId string, followedId string) error {

	// check followed user didn't ban the logged user
	if banned, err := db.IsBannedBy(followerId, followedId); banned || err != nil {
		return fmt.Errorf("impossible to follow the user; %w", err)
	}

	// check user isn't already followed
	followed, err := db.IsFollowedBy(followedId, followerId)
	if followed {
		return fmt.Errorf("the user is already followed")
	}
	if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}

	stmt, _ := db.c.Prepare("INSERT INTO follows (followerid, followedid) VALUES (?, ?);")
	_, err = stmt.Exec(followerId, followedId)
	if err != nil {
		return fmt.Errorf("error adding the follow: %w", err)
	}
	return nil
}

// unfollow loggedUser
func (db *appdbimpl) UnfollowUser(followerId string, followedId string) error {

	followed, err := db.IsFollowedBy(followedId, followerId)
	if !followed {
		return fmt.Errorf("the user is not followed at the moment")
	}
	if err != nil {
		return fmt.Errorf("error removing the follow: %w", err)
	}

	sqlStmt := "DELETE FROM follows WHERE followerid=? AND followedid=?"
	_, err = db.c.Exec(sqlStmt, followerId, followedId)
	return err
}

func (db *appdbimpl) GetFollowers(userid string) ([]string, error) {

	rows, err := db.c.Query("SELECT u.username FROM users u JOIN follows f ON u.userid = f.followerid WHERE f.followedid = ?", userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		followers = append(followers, username)
	}
	if rows.Err() != nil {
		return followers, rows.Err()
	}
	return followers, nil
}

func (db *appdbimpl) GetFollowing(userid string) ([]string, error) {
	rows, err := db.c.Query("SELECT u.username FROM users u JOIN follows f ON u.userid = f.followedid WHERE f.followerid = ?", userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		following = append(following, username)
	}
	if rows.Err() != nil {
		return following, rows.Err()
	}
	return following, nil
}

func (db *appdbimpl) IsFollowedBy(followedId string, followerId string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE followerid=? AND followedid=?", followerId, followedId).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return count > 0, nil
}
