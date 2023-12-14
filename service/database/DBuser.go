package database

import (
	"database/sql"
	"fmt"
)

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetUser(id string, username string) error {
	stmt1, _ := db.c.Prepare("INSERT INTO users (userid, username) VALUES (?, ?);")
	stmt2, _ := db.c.Prepare("INSERT INTO profiles (userid) VALUES (?);")
	_, err1 := stmt1.Exec(id, username)
	_, err2 := stmt2.Exec(id)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("error in profie creation\n err1: %v\n err2: %v", err1, err2)
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
			return user, fmt.Errorf("userByUsername %s: no such user", username)
		}
		return user, fmt.Errorf("userByUsername %s: %v", username, err)
	}
	return user, nil
}

func (db *appdbimpl) CheckID(id string) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE userid = ?", id).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("search id %s: %v", id, err)
	}
	return count, nil
}

// add func getProfile
