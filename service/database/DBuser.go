package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) SetUser(userid string, username string) error {
	_, err := db.c.Exec("INSERT INTO users (userid, username) VALUES (?, ?);", userid, username)
	if err != nil {
		return fmt.Errorf("error in profie creation. err: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetName(reqUserId string, newUsername string) error {

	stmt, _ := db.c.Prepare("UPDATE users SET username = ? WHERE userid = ?;")
	_, err := stmt.Exec(newUsername, reqUserId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.UserID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, fmt.Errorf("error retrieving the user: %w", err)
	}
	return user, nil
}

func (db *appdbimpl) GetById(userid string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT * FROM users WHERE userid=?", userid).Scan(&user.UserID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no such user")
		}
		return user, fmt.Errorf("error retrieving the user: %w", err)
	}
	return user, nil
}

func (db *appdbimpl) CheckIDExistence(userID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE userid = ?)", userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}
	return exists, nil
}

func (db *appdbimpl) GetProfile(userid string) (Profile, error) {

	profile := Profile{}

	// get user
	user, err := db.GetById(userid)
	if err != nil {
		return profile, fmt.Errorf(errRetrievingProfile, err)
	}
	profile.Username = user.Username

	// get photos
	photoRows, err := db.c.Query("SELECT photoid FROM photos WHERE authorid = ? ORDER BY dateAndTime DESC;", userid)
	if err != nil {
		return profile, fmt.Errorf(errRetrievingProfile, err)
	}
	defer photoRows.Close()

	for photoRows.Next() {

		var photoid string

		if err := photoRows.Scan(&photoid); err != nil {
			return profile, fmt.Errorf(errRetrievingProfile, err)
		}

		photoData, err := db.GetPhotoData(photoid)
		if err != nil {
			return profile, fmt.Errorf(errRetrievingProfile, err)
		}
		profile.Posts = append(profile.Posts, photoData)
	}

	// get followers
	profile.Followers.Usernames, err = db.GetFollowers(userid)
	if err != nil {
		return profile, fmt.Errorf(errRetrievingProfile, err)
	}
	profile.Followers.NumberOfFollowers = len(profile.Followers.Usernames)

	// get following
	profile.Following.Usernames, err = db.GetFollowing(userid)
	if err != nil {
		return profile, fmt.Errorf(errRetrievingProfile, err)
	}
	profile.Following.NumberOfFollowing = len(profile.Following.Usernames)

	// get ban status
	profile.Banned = false

	if photoRows.Err() != nil {
		return profile, photoRows.Err()
	}
	return profile, nil
}
