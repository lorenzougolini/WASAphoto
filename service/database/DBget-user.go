package database

import (
	"database/sql"
	"fmt"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUser(username string) (User, error) {
	var user User
	row := db.c.QueryRow("SELECT * FROM users WHERE username=?", username)
	if err := row.Scan(&user.UserID, &user.Username); err != nil {
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
		return 0, fmt.Errorf("Search id %s: %v", id, err)
	}
	return count, nil
}

// add func getProfile
