package database

import (
	"database/sql"
	"fmt"
)

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetUser(id string, username string) error {
	stmt, _ := db.c.Prepare("INSERT INTO users (userid, username) VALUES (?, ?);")
	_, err := stmt.Exec(id, username)
	if err != nil {
		return fmt.Errorf("error in profie creation. err: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetName(id string, newUsername string) error {
	stmt, err := db.c.Prepare("UPDATE users SET username = ? WHERE userid = ?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newUsername, id)
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
			return user, fmt.Errorf("user ByUsername %s: no such user", username)
		}
		return user, fmt.Errorf("user ByUsername %s: %w", username, err)
	}
	return user, nil
}

func (db *appdbimpl) GetById(userid string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT * FROM users WHERE userid=?", userid).Scan(&user.UserID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user ByUsername %s: no such user", userid)
		}
		return user, fmt.Errorf("user ByUsername %s: %w", userid, err)
	}
	return user, nil
}

func (db *appdbimpl) CheckID(id string) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE userid = ?", id).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("search id %s: %w", id, err)
	}
	return count, nil
}

func (db *appdbimpl) GetProfile(userid string) (Profile, error) {
	profile := Profile{}

	// get photos
	rows, err := db.c.Query("SELECT picture, dateAndTime, description FROM photos WHERE userid = ?", userid)
	if err != nil {
		return profile, fmt.Errorf("error retrieving the profile. err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		photo := struct {
			File        string
			Description string
			DateAndTime string
		}{}

		if err := rows.Scan(&photo.File, &photo.DateAndTime, &photo.Description); err != nil {
			return profile, fmt.Errorf("error retrieving the profile. err: %w", err)
		}

		profile.Posts.Photos = append(profile.Posts.Photos, photo)
	}
	profile.Posts.NumberOfPosts = len(profile.Posts.Photos)

	// get followers
	rows, err = db.c.Query("SELECT userid FROM follows WHERE followedid = ?", userid)
	if err != nil {
		return profile, fmt.Errorf("error retrieving the profile. err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var followerid string

		if err := rows.Scan(&followerid); err != nil {
			return profile, fmt.Errorf("error retrieving the profile. err: %w", err)
		}

		if user, err := db.GetById(followerid); err != nil {
			return profile, fmt.Errorf("error retrieving the profile. err: %w", err)

		} else {
			profile.Followers.Usernames = append(profile.Followers.Usernames, user.Username)
		}
	}
	profile.Followers.NumberOfFollowers = len(profile.Followers.Usernames)

	// get following
	rows, err = db.c.Query("SELECT followedid FROM follows WHERE userid = ?", userid)
	if err != nil {
		return profile, fmt.Errorf("error retrieving the profile. err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var followedid string

		if err := rows.Scan(&followedid); err != nil {
			return profile, fmt.Errorf("error retrieving the profile. err: %w", err)

		}

		if user, err := db.GetById(followedid); err != nil {
			return profile, fmt.Errorf("error retrieving the profile. err: %w", err)

		} else {
			profile.Following.Usernames = append(profile.Following.Usernames, user.Username)

		}
	}
	profile.Following.NumberOfFollowing = len(profile.Following.Usernames)

	return profile, nil
}
